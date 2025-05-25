// Package version provides information about the build version, git commit, and build date set by ldflags at build time.
package version

// All of these vars are replaced at link time using ldflags.
var (
	// Version is a valid semantic version
	Version = "unknown"
	// GitCommit contains the git SHA of the commit
	GitCommit = "unknown"
	// BuildDate is the date of the build
	BuildDate = "unknown"
)
