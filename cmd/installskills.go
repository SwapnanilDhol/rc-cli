package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func initInstallSkills(root *cobra.Command, skillsFS fs.FS, skillsRoot string) {
	cmd := &cobra.Command{
		Use:   "install-skills",
		Short: "Install the rc agent skills so coding agents can use them anywhere",
		Long: `Copy this CLI's agent skills into your global Claude skills directory
(~/.claude/skills by default), making them available in every project rather
than only inside a checkout of this repo.

Skills are embedded in the binary, so this works from a Homebrew install.

  rc install-skills            # install, skipping any that already exist
  rc install-skills --force    # overwrite existing copies
  rc install-skills --list     # show what would be installed
  rc install-skills --dir ./x  # install somewhere else`,
		RunE: func(c *cobra.Command, args []string) error {
			return runInstallSkills(c, skillsFS, skillsRoot)
		},
	}
	cmd.Flags().Bool("force", false, "Overwrite skills that are already installed")
	cmd.Flags().Bool("list", false, "List the bundled skills without installing")
	cmd.Flags().String("dir", "", "Target directory (default ~/.claude/skills)")
	root.AddCommand(cmd)
}

// bundledSkills returns the skill directory names embedded in the binary.
func bundledSkills(skillsFS fs.FS, skillsRoot string) ([]string, error) {
	if skillsFS == nil {
		return nil, fmt.Errorf("no skills embedded in this binary")
	}
	entries, err := fs.ReadDir(skillsFS, skillsRoot)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func runInstallSkills(cmd *cobra.Command, skillsFS fs.FS, skillsRoot string) error {
	names, err := bundledSkills(skillsFS, skillsRoot)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no skills embedded in this binary")
	}

	list, _ := cmd.Flags().GetBool("list")
	if list {
		if jsonRequested(cmd) {
			return emitJSON(names)
		}
		fmt.Println(appsStyle.Render("\nBundled skills:\n"))
		for _, n := range names {
			fmt.Printf("  %s\n", cyanStyle.Render(n))
		}
		return nil
	}

	target, _ := cmd.Flags().GetString("dir")
	if target == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory; pass --dir: %w", err)
		}
		target = filepath.Join(home, ".claude", "skills")
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	force, _ := cmd.Flags().GetBool("force")
	var installed, skipped []string
	for _, name := range names {
		dest := filepath.Join(target, name)
		if _, err := os.Stat(dest); err == nil && !force {
			skipped = append(skipped, name)
			continue
		}
		if err := copySkill(skillsFS, skillsRoot, name, dest); err != nil {
			return fmt.Errorf("installing %s: %w", name, err)
		}
		installed = append(installed, name)
	}

	if jsonRequested(cmd) {
		return emitJSON(map[string]interface{}{
			"target": target, "installed": installed, "skipped": skipped,
		})
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	for _, n := range installed {
		fmt.Printf("%s %s\n", green.Render("✓"), n)
	}
	for _, n := range skipped {
		fmt.Printf("%s %s %s\n", gray.Render("-"), n, gray.Render("(already installed; --force to overwrite)"))
	}
	fmt.Println(gray.Render(fmt.Sprintf("\n%d installed, %d skipped → %s", len(installed), len(skipped), target)))
	return nil
}

// copySkill writes one embedded skill directory to dest.
func copySkill(skillsFS fs.FS, skillsRoot, name, dest string) error {
	srcRoot := skillsRoot + "/" + name
	return fs.WalkDir(skillsFS, srcRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(p, srcRoot), "/")
		out := filepath.Join(dest, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(out, 0755)
		}
		data, err := fs.ReadFile(skillsFS, p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return err
		}
		return os.WriteFile(out, data, 0644)
	})
}
