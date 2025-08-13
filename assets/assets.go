package assets

import "embed"

var (
	//go:embed templates
	Templates embed.FS

	//go:embed dist
	Dist embed.FS
)
