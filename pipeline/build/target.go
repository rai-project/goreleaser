package build

import (
	"fmt"
	"log"
	"runtime"

	"github.com/rai-project/goreleaser/context"
)

var runtimeTarget = buildTarget{runtime.GOOS, runtime.GOARCH, ""}

// a build target
type buildTarget struct {
	goos, goarch, goarm string
}

func (t buildTarget) String() string {
	return fmt.Sprintf("%v%v%v", t.goos, t.goarch, t.goarm)
}

func (t buildTarget) PrettyString() string {
	return fmt.Sprintf("%v/%v%v", t.goos, t.goarch, t.goarm)
}

func buildTargets(ctx *context.Context) (targets []buildTarget) {
	for _, target := range allBuildTargets(ctx) {
		if !valid(target) {
			log.Println("Skipped invalid build target:", target.PrettyString())
			continue
		}
		if ignored(ctx, target) {
			log.Println("Skipped ignored build target:", target.PrettyString())
			continue
		}
		targets = append(targets, target)
	}
	return
}

func allBuildTargets(ctx *context.Context) (targets []buildTarget) {
	for _, goos := range ctx.Config.Build.Goos {
		for _, goarch := range ctx.Config.Build.Goarch {
			if goarch == "arm" {
				for _, goarm := range ctx.Config.Build.Goarm {
					targets = append(targets, buildTarget{goos, goarch, goarm})
				}
				continue
			}
			targets = append(targets, buildTarget{goos, goarch, ""})
		}
	}
	return
}

func ignored(ctx *context.Context, target buildTarget) bool {
	for _, ig := range ctx.Config.Build.Ignore {
		var ignored = buildTarget{ig.Goos, ig.Goarch, ig.Goarm}
		if ignored == target {
			return true
		}
	}
	return false
}

func valid(target buildTarget) bool {
	var s = target.goos + target.goarch
	for _, a := range valids {
		if a == s {
			return true
		}
	}
	return false
}
