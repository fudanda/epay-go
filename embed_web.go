package epaygo

import "embed"

// WebDist holds the built frontend assets produced by `npm run build` in web/.
//
// all: is required because Vite emits helper chunks whose filenames may start
// with an underscore, and directory embeds skip those by default.
//
//go:embed all:web/dist
var WebDist embed.FS
