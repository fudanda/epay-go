[CmdletBinding()]
param(
    [string]$OutputDir = ".\releases",
    [ValidateSet("linux")]
    [string]$GoOs = "linux",
    [ValidateSet("amd64", "arm64")]
    [string]$GoArch = "amd64",
    [switch]$Force
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Require-Command {
    param([Parameter(Mandatory = $true)][string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Missing required command: $Name"
    }
}

function New-CleanDirectory {
    param([Parameter(Mandatory = $true)][string]$Path)

    if (Test-Path -LiteralPath $Path) {
        Remove-Item -LiteralPath $Path -Recurse -Force
    }
    New-Item -ItemType Directory -Path $Path -Force | Out-Null
}

function Copy-IfExists {
    param(
        [Parameter(Mandatory = $true)][string]$Source,
        [Parameter(Mandatory = $true)][string]$Destination
    )

    if (Test-Path -LiteralPath $Source) {
        Copy-Item -LiteralPath $Source -Destination $Destination -Recurse -Force
    }
}

function Compress-Directory {
    param(
        [Parameter(Mandatory = $true)][string]$SourceDir,
        [Parameter(Mandatory = $true)][string]$ZipPath
    )

    Add-Type -AssemblyName System.IO.Compression
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    if (Test-Path -LiteralPath $ZipPath) {
        Remove-Item -LiteralPath $ZipPath -Force
    }

    $archive = [System.IO.Compression.ZipFile]::Open($ZipPath, [System.IO.Compression.ZipArchiveMode]::Create)
    try {
        $sourceRoot = (Resolve-Path -LiteralPath $SourceDir).Path
        $rootLength = $sourceRoot.Length + 1
        Get-ChildItem -LiteralPath $sourceRoot -Recurse -File | ForEach-Object {
            $relativePath = $_.FullName.Substring($rootLength).Replace('\', '/')
            [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
                $archive,
                $_.FullName,
                $relativePath,
                [System.IO.Compression.CompressionLevel]::Optimal
            ) | Out-Null
        }
    }
    finally {
        $archive.Dispose()
    }
}

$repoRoot = Split-Path -Parent $PSCommandPath
$outputRoot = [System.IO.Path]::GetFullPath((Join-Path $repoRoot $OutputDir))
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$bundleName = "epay-go-$GoOs-$GoArch-$timestamp"
$stagingDir = Join-Path $outputRoot $bundleName
$zipPath = Join-Path $outputRoot "$bundleName.zip"
$binaryPath = Join-Path $stagingDir "epay-server"
$configPath = Join-Path $repoRoot "config.yaml"

Require-Command "go"
Require-Command "npm"

if (-not (Test-Path -LiteralPath $configPath)) {
    throw "Missing required runtime config: $configPath"
}

if ((Test-Path -LiteralPath $zipPath) -and -not $Force) {
    throw "Package already exists: $zipPath`nRe-run with -Force to overwrite."
}

New-Item -ItemType Directory -Path $outputRoot -Force | Out-Null
New-CleanDirectory -Path $stagingDir

$webDir = Join-Path $repoRoot "web"
Push-Location $webDir
try {
    if (-not (Test-Path -LiteralPath (Join-Path $webDir "node_modules"))) {
        npm ci
    }
    npm run build
}
finally {
    Pop-Location
}

$previousEnv = @{
    CGO_ENABLED = $env:CGO_ENABLED
    GOOS        = $env:GOOS
    GOARCH      = $env:GOARCH
}

Push-Location $repoRoot
try {
    $env:CGO_ENABLED = "0"
    $env:GOOS = $GoOs
    $env:GOARCH = $GoArch

    go build -trimpath -ldflags="-s -w" -o $binaryPath .\cmd\server\main.go
}
finally {
    $env:CGO_ENABLED = $previousEnv.CGO_ENABLED
    $env:GOOS = $previousEnv.GOOS
    $env:GOARCH = $previousEnv.GOARCH
    Pop-Location
}

Copy-Item -LiteralPath $configPath -Destination $stagingDir -Force
Copy-IfExists -Source (Join-Path $repoRoot ".env") -Destination $stagingDir
Copy-IfExists -Source (Join-Path $repoRoot ".env.example") -Destination $stagingDir
Copy-IfExists -Source (Join-Path $repoRoot "config.example.yaml") -Destination $stagingDir
Copy-IfExists -Source (Join-Path $repoRoot "README.md") -Destination $stagingDir
Copy-IfExists -Source (Join-Path $repoRoot "DEPLOYMENT.md") -Destination $stagingDir
Copy-IfExists -Source (Join-Path $repoRoot "deploy") -Destination (Join-Path $stagingDir "deploy")

$manifest = @(
    "epay-server"
    "config.yaml"
    ".env"
    ".env.example"
    "config.example.yaml"
    "README.md"
    "DEPLOYMENT.md"
    "deploy/"
) -join [Environment]::NewLine
Set-Content -LiteralPath (Join-Path $stagingDir "package-manifest.txt") -Value $manifest -NoNewline

Compress-Directory -SourceDir $stagingDir -ZipPath $zipPath
Remove-Item -LiteralPath $stagingDir -Recurse -Force

Write-Host "Created package: $zipPath"
