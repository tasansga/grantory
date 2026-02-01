package cli

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Version and Commit are injected at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "unknown"
)

func versionString() string {
	version := strings.TrimSpace(Version)
	commit := strings.TrimSpace(Commit)

	if version == "" || version == "dev" {
		if commit == "" || commit == "unknown" {
			if rev, modified := buildInfoRevision(); rev != "" {
				if modified {
					commit = rev + "-dirty"
				} else {
					commit = rev
				}
			}
		}
		if commit != "" && commit != "unknown" {
			return fmt.Sprintf("dev (%s)", commit)
		}
		return "dev"
	}

	if commit != "" && commit != "unknown" {
		return fmt.Sprintf("%s (%s)", version, commit)
	}
	return version
}

func buildInfoRevision() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	var revision string
	var modified bool
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	return revision, modified
}
