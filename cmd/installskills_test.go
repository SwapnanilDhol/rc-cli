package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// The real embed lives in the root package, so exercise the copy logic against
// an in-memory FS with the same shape. Passing it in per-test means no shared
// state and no ordering between tests.
func testSkills() (fs.FS, string) {
	return fstest.MapFS{
		".claude/skills/rc-cli-usage/SKILL.md":      {Data: []byte("---\nname: rc-cli-usage\n---\nbody")},
		".claude/skills/rc-asc-bridge/SKILL.md":     {Data: []byte("---\nname: rc-asc-bridge\n---\nbody")},
		".claude/skills/rc-asc-bridge/ref/extra.md": {Data: []byte("nested")},
		".claude/skills/README.md":                  {Data: []byte("index, not a skill")},
	}, ".claude/skills"
}

func TestBundledSkillsListsOnlyDirectories(t *testing.T) {
	names, err := bundledSkills(testSkills())
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("got %v, want the 2 skill directories (README.md must not be listed)", names)
	}
}

func TestCopySkillIncludesNestedFiles(t *testing.T) {
	fsys, root := testSkills()
	dest := filepath.Join(t.TempDir(), "rc-asc-bridge")
	if err := copySkill(fsys, root, "rc-asc-bridge", dest); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"SKILL.md", filepath.Join("ref", "extra.md")} {
		if _, err := os.Stat(filepath.Join(dest, rel)); err != nil {
			t.Errorf("expected %s to be copied: %v", rel, err)
		}
	}
	// The skill root itself must not be nested inside the destination again.
	if _, err := os.Stat(filepath.Join(dest, "rc-asc-bridge")); err == nil {
		t.Error("skill was copied one directory too deep")
	}
}

func TestCopySkillWritesReadableContent(t *testing.T) {
	fsys, root := testSkills()
	dest := filepath.Join(t.TempDir(), "s")
	if err := copySkill(fsys, root, "rc-cli-usage", dest); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("copied skill is empty")
	}
}

func TestBundledSkillsErrorsWhenNothingEmbedded(t *testing.T) {
	if _, err := bundledSkills(nil, ""); err == nil {
		t.Fatal("expected an error when no skills are embedded")
	}
}
