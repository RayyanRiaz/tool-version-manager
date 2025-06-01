package models

import (
	"log/slog"
	"strings"

	"github.com/hashicorp/go-version"
)

type ToolComparerWithVersionParsing struct{}

func (t *ToolComparerWithVersionParsing) CompareVersions(tool Tool, v1 ToolVersion, v2 ToolVersion) (int, error) {
	// remove all non-numeric characters from the beginning of the version strings
	stripFunc := func(r rune) bool {
		return r < '0' || r > '9'
	}
	v1Stripped := strings.TrimLeftFunc(string(v1), stripFunc)
	v2Stripped := strings.TrimLeftFunc(string(v2), stripFunc)

	// if either version is empty after stripping, fall back to string comparison
	if v1Stripped == "" || v2Stripped == "" {
		slog.Warn("One of the versions is empty after stripping non-numeric characters (%v, %v), falling back to string comparison", string(v1), string(v2))
		return strings.Compare(string(v1), string(v2)), nil
	}

	// parse the stripped versions using hashicorp/go-version. If parsing fails, fall back to string comparison
	v1Parsed, e1 := version.NewVersion(v1Stripped)
	v2Parsed, e2 := version.NewVersion(v2Stripped)

	if e1 != nil || e2 != nil {
		slog.Warn("Failed to parse versions(%v, %v), falling back to string comparison", string(v1), string(v2))
		return strings.Compare(string(v1), string(v2)), nil
	}

	return v1Parsed.Compare(v2Parsed), nil
}

var _ ToolComparer = &ToolComparerWithVersionParsing{}
