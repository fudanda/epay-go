package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadAppliesDotEnvValues(t *testing.T) {
	t.Cleanup(viper.Reset)
	unsetEnv(t, "DB_NAME")
	unsetEnv(t, "DEFAULT_ADMIN_PASSWORD")

	dir := t.TempDir()
	writeTestConfig(t, dir)
	writeTestDotEnv(t, dir, "DB_NAME=from_dot_env\nDEFAULT_ADMIN_PASSWORD=FromDotEnv123!\n")

	if err := Load(filepath.Join(dir, "config.yaml")); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := Get().Database.DBName; got != "from_dot_env" {
		t.Fatalf("database dbname = %q, want %q", got, "from_dot_env")
	}

	if got := os.Getenv("DEFAULT_ADMIN_PASSWORD"); got != "FromDotEnv123!" {
		t.Fatalf("DEFAULT_ADMIN_PASSWORD = %q, want %q", got, "FromDotEnv123!")
	}
}

func TestLoadKeepsExistingEnvironmentPriority(t *testing.T) {
	t.Cleanup(viper.Reset)
	t.Setenv("DB_NAME", "from_process_env")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "FromProcessEnv123!")

	dir := t.TempDir()
	writeTestConfig(t, dir)
	writeTestDotEnv(t, dir, "DB_NAME=from_dot_env\nDEFAULT_ADMIN_PASSWORD=FromDotEnv123!\n")

	if err := Load(filepath.Join(dir, "config.yaml")); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := Get().Database.DBName; got != "from_process_env" {
		t.Fatalf("database dbname = %q, want %q", got, "from_process_env")
	}

	if got := os.Getenv("DEFAULT_ADMIN_PASSWORD"); got != "FromProcessEnv123!" {
		t.Fatalf("DEFAULT_ADMIN_PASSWORD = %q, want %q", got, "FromProcessEnv123!")
	}
}

func writeTestConfig(t *testing.T, dir string) {
	t.Helper()

	content := `server:
  port: 8090
  mode: debug

database:
  host: 127.0.0.1
  port: 25432
  user: postgres.lzserver
  password: password
  dbname: from_config
  sslmode: disable

redis:
  host: 127.0.0.1
  port: 26379
  password: redis_password
  db: 1

jwt:
  secret: secret
  expire_hour: 168
`

	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(config.yaml) error = %v", err)
	}
}

func writeTestDotEnv(t *testing.T, dir string, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(.env) error = %v", err)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	oldValue, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv(%s) error = %v", key, err)
	}
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, oldValue)
			return
		}
		_ = os.Unsetenv(key)
	})
}
