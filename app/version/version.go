package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

func Info() string {
	return fmt.Sprintf("Version: %s, Git Commit: %s, Build Time: %s, Go Version: %s",
		Version, GitCommit, BuildTime, GoVersion)
}

func Short() string {
	commit := GitCommit
	if len(commit) > 8 {
		commit = commit[:8]
	}
	return fmt.Sprintf("v%s (%s)", Version, commit)
}
