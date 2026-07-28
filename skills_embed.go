package main

import "embed"

// Skills ship inside the binary so `rc install-skills` works from a Homebrew
// install with no repo checkout. The embed directive must live in the root
// package because go:embed cannot reference paths above its own directory.
//
//go:embed all:.claude/skills
var skillsFS embed.FS

const skillsRoot = ".claude/skills"
