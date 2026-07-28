package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func testRoot(t *testing.T) *cobra.Command {
	t.Helper()
	return NewRootCmd(nil, "")
}

// walk visits every command in the tree, depth first.
func walk(c *cobra.Command, fn func(*cobra.Command)) {
	fn(c)
	for _, sub := range c.Commands() {
		walk(sub, fn)
	}
}

func fullName(c *cobra.Command) string {
	var parts []string
	for cur := c; cur != nil; cur = cur.Parent() {
		parts = append([]string{cur.Name()}, parts...)
	}
	return strings.Join(parts, " ")
}

// A local flag reusing a global shorthand makes cobra panic at invocation time,
// not at registration — so nothing catches it until someone runs that exact
// command. Adding the global -p/--project-id collided with subscribers list's
// local -p/--platform and hard-panicked `rc subscribers list`.
func TestNoShorthandCollisionsWithPersistentFlags(t *testing.T) {
	root := testRoot(t)
	global := map[string]string{}
	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Shorthand != "" {
			global[f.Shorthand] = f.Name
		}
	})
	if len(global) == 0 {
		t.Fatal("expected the root command to define shorthand persistent flags")
	}

	walk(root, func(c *cobra.Command) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Shorthand == "" {
				return
			}
			if name, taken := global[f.Shorthand]; taken && name != f.Name {
				t.Errorf("%q: local -%s/--%s collides with global -%s/--%s; cobra will panic when this command runs",
					fullName(c), f.Shorthand, f.Name, f.Shorthand, name)
			}
		})
	})
}

// Executing --help on every command exercises cobra's flag merging, which is
// where a collision actually blows up.
func TestEveryCommandHelpDoesNotPanic(t *testing.T) {
	root := testRoot(t)
	var names []string
	walk(root, func(c *cobra.Command) {
		if c.Name() != "completion" && c.Name() != "help" {
			names = append(names, fullName(c))
		}
	})
	if len(names) < 20 {
		t.Fatalf("expected a large command tree, found %d", len(names))
	}

	for _, name := range names {
		args := append(strings.Split(name, " ")[1:], "--help")
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic running 'rc %s --help': %v", strings.Join(args, " "), r)
				}
			}()
			c, _, err := root.Find(args)
			if err != nil {
				t.Fatalf("could not find command: %v", err)
			}
			// Merging persistent flags is what panics on a shorthand collision.
			c.InitDefaultHelpFlag()
			c.LocalFlags()
			c.InheritedFlags()
		})
	}
}
