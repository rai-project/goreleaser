package archiveformat

import (
	"testing"

	"github.com/rai-project/goreleaser/config"
	"github.com/rai-project/goreleaser/context"
	"github.com/stretchr/testify/assert"
)

func TestFormatFor(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{
		Config: config.Project{
			Archive: config.Archive{
				Format: "tar.gz",
				FormatOverrides: []config.FormatOverride{
					{
						Goos:   "windows",
						Format: "zip",
					},
				},
			},
		},
	}
	assert.Equal("zip", For(ctx, "windowsamd64"))
	assert.Equal("tar.gz", For(ctx, "linux386"))
}
