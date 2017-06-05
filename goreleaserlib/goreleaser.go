package goreleaserlib

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v1"

	"fmt"

	"github.com/rai-project/goreleaser/config"
	"github.com/rai-project/goreleaser/context"
	"github.com/rai-project/goreleaser/pipeline"
	"github.com/rai-project/goreleaser/pipeline/archive"
	"github.com/rai-project/goreleaser/pipeline/brew"
	"github.com/rai-project/goreleaser/pipeline/build"
	"github.com/rai-project/goreleaser/pipeline/checksums"
	"github.com/rai-project/goreleaser/pipeline/defaults"
	"github.com/rai-project/goreleaser/pipeline/env"
	"github.com/rai-project/goreleaser/pipeline/fpm"
	"github.com/rai-project/goreleaser/pipeline/git"
	"github.com/rai-project/goreleaser/pipeline/release"
)

var pipes = []pipeline.Pipe{
	defaults.Pipe{},  // load default configs
	git.Pipe{},       // get and validate git repo state
	env.Pipe{},       // load and validate environment variables
	build.Pipe{},     // build
	archive.Pipe{},   // archive (tar.gz, zip, etc)
	fpm.Pipe{},       // archive via fpm (deb, rpm, etc)
	checksums.Pipe{}, // checksums of the files
	release.Pipe{},   // release to github
	brew.Pipe{},      // push to brew tap
}

func init() {
	log.SetFlags(0)
}

// Flags interface represents an extractor of cli flags
type Flags interface {
	IsSet(s string) bool
	String(s string) string
	Bool(s string) bool
}

// Release runs the release process with the given flags
func Release(flags Flags) error {
	var file = flags.String("config")
	var notes = flags.String("release-notes")
	cfg, err := config.Load(file)
	if err != nil {
		// Allow file not found errors if config file was not
		// explicitly specified
		_, statErr := os.Stat(file)
		if !os.IsNotExist(statErr) || flags.IsSet("config") {
			return err
		}
		log.Printf("WARNING: Could not load %v\n", file)
	}
	var ctx = context.New(cfg)
	ctx.Validate = !flags.Bool("skip-validate")
	ctx.Publish = !flags.Bool("skip-publish")
	if notes != "" {
		bts, err := ioutil.ReadFile(notes)
		if err != nil {
			return err
		}
		log.Println("Loaded custom release notes from", notes)
		ctx.ReleaseNotes = string(bts)
	}
	ctx.Snapshot = flags.Bool("snapshot")
	if ctx.Snapshot {
		log.Println("Publishing disabled in snapshot mode")
		ctx.Publish = false
	}
	for _, pipe := range pipes {
		log.Println(pipe.Description())
		log.SetPrefix(" -> ")
		if err := pipe.Run(ctx); err != nil {
			return err
		}
		log.SetPrefix("")
	}
	log.Println("Done!")
	return nil
}

// InitProject creates an example goreleaser.yml in the current directory
func InitProject(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s already exists", filename)
	}

	var ctx = context.New(config.Project{})
	var pipe = defaults.Pipe{}
	if err := pipe.Run(ctx); err != nil {
		return err
	}
	out, err := yaml.Marshal(ctx.Config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, out, 0644)
}
