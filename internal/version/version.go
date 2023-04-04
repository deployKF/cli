package version

import (
	"flag"
	"runtime"
)

var (
	// version is the current version of deployKF
	// NOTE: value is overwritten automatically during build
	version = "v0.0.0"

	// gitCommit is the git sha1
	// NOTE: value is overwritten automatically during build
	gitCommit = ""

	// gitTreeState is the state of the git tree
	// NOTE: value is overwritten automatically during build
	gitTreeState = ""
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`

	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`

	// GitTreeState is the state of the git tree.
	GitTreeState string `json:"git_tree_state,omitempty"`

	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// GetVersion returns the semver string of the version
func GetVersion() string {
	return version
}

// Get returns build info
func Get() BuildInfo {
	v := BuildInfo{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	// strip out GoVersion during a test run for consistent test output
	if flag.Lookup("test.v") != nil {
		v.GoVersion = ""
	}
	return v
}
