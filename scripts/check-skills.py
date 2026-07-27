#!/usr/bin/env python3
"""Verify every command and flag referenced in .claude/skills/ actually exists.

The skills describe a live command surface, so they rot silently: a renamed flag
still reads fine but costs an agent a failed command and a recovery loop. This
walks every `rc …` line in the skills and checks it against the built binary.

Usage: ./scripts/check-skills.py [path-to-rc-binary]
"""

import glob
import os
import re
import shlex
import shutil
import subprocess
import sys

RC = sys.argv[1] if len(sys.argv) > 1 else "./rc"
SKILLS = ".claude/skills/*/SKILL.md"

# Third-party CLIs the skills reference. Checked the same way when installed, and
# skipped otherwise so CI does not depend on them being present.
EXTERNAL = {"asc": shutil.which("asc")}

# Placeholders that appear in documented examples but are not real subcommands.
PLACEHOLDER_PREFIXES = ("<", "{", "'", '"', "$", "|")


# Subcommands are always lowercase words. Anything else (an HTTP verb, a URL path,
# a placeholder) is a positional argument and ends the subcommand path.
SUBCOMMAND_RE = re.compile(r"^[a-z][a-z0-9-]*$")

# Commands whose arguments are free-form and must not be read as subcommands.
POSITIONAL_CMDS = {("api",), ("internal", "api")}


def tokenize(line):
    """Split a documented shell line, tolerating unbalanced quotes in prose."""
    line = line.strip().rstrip("\\").strip()
    for splitter in (
        lambda s: shlex.split(s, comments=True),
        lambda s: shlex.split(s),
        lambda s: s.split(),
    ):
        try:
            return splitter(line)
        except ValueError:
            continue
    return line.split()


def parse(line):
    """Yield (subcommand path, flags) for each `rc …` invocation on a line.

    Skill files put several commands on one line in reference tables, so a line
    can contain more than one invocation.
    """
    toks = tokenize(line)
    binaries = {"rc", *EXTERNAL}
    if not binaries & set(toks):
        return

    # Split into one segment per binary occurrence.
    segments, cur, cur_bin = [], None, None
    for t in toks:
        if t in binaries:
            if cur is not None:
                segments.append((cur_bin, cur))
            cur, cur_bin = [], t
        elif cur is not None:
            cur.append(t)
    if cur:
        segments.append((cur_bin, cur))

    for binary, seg in segments:
        path, flags = [], []
        seen_flag = False
        for t in seg:
            if t.startswith("-") and len(t) > 1 and not t[1].isdigit():
                seen_flag = True
                flags.append(t.lstrip("-"))
                continue
            if seen_flag:
                continue  # a flag's value
            if not SUBCOMMAND_RE.match(t) or t.startswith(PLACEHOLDER_PREFIXES):
                break
            path.append(t)
            if tuple(path) in POSITIONAL_CMDS:
                seen_flag = True  # stop reading positionals as subcommands
        if path:
            yield binary, tuple(path), flags


def main():
    if not os.path.exists(RC):
        sys.exit(f"binary not found: {RC} (run: go build -o rc .)")

    files = sorted(glob.glob(SKILLS))
    if not files:
        sys.exit(f"no skills found at {SKILLS}")

    cache = {}
    missing_cmds, missing_flags = [], []
    n_cmds = n_flags = 0

    def help_for(binary, path):
        """Return (exists, help text) for a subcommand path.

        Cobra prints the *parent's* help and exits 0 for an unknown subcommand, so
        neither the exit code nor an "unknown command" string is reliable. The
        Usage line is: it echoes the full path only when the path really resolves.
        """
        key = (binary, path)
        if key not in cache:
            exe = RC if binary == "rc" else EXTERNAL[binary]
            r = subprocess.run([exe, *path, "--help"], capture_output=True, text=True)
            text = r.stdout + r.stderr
            usage = ""
            lines = text.splitlines()
            for i, ln in enumerate(lines):
                if ln.strip() == "Usage:" and i + 1 < len(lines):
                    usage = lines[i + 1].strip()
                    break
            # asc prints "USAGE" and does not echo the full path, so fall back to
            # requiring a clean exit plus no unknown-command complaint.
            if binary == "rc":
                exists = r.returncode == 0 and usage.startswith("rc " + " ".join(path))
            else:
                exists = r.returncode == 0 and "unknown" not in text.lower()[:400]
            cache[key] = (exists, text)
        return cache[key]

    for f in files:
        skill = f.split(os.sep)[2]
        in_code = False
        for line in open(f):
            # Only fenced code blocks are commands; prose mentions rc too.
            if line.lstrip().startswith("```"):
                in_code = not in_code
                continue
            if not in_code:
                continue
            for binary, path, flags in parse(line):
                if binary != "rc" and not EXTERNAL.get(binary):
                    continue  # not installed here; skip rather than fail
                exists, text = help_for(binary, path)
                n_cmds += 1
                if not exists:
                    missing_cmds.append((skill, f"{binary} " + " ".join(path)))
                    continue

                for fl in flags:
                    if fl in ("help", "h"):
                        continue  # universal, and not always self-documented
                    n_flags += 1
                    # Long flags appear as --name, short ones as -x, in the help.
                    needle = f"--{fl}" if len(fl) > 1 else f"-{fl}"
                    if needle not in text:
                        missing_flags.append((skill, f"{binary} " + " ".join(path), fl))

    print(f"checked {n_cmds} command usages and {n_flags} flag usages "
          f"across {len(files)} skills")

    if not missing_cmds and not missing_flags:
        print("OK — skills match the CLI")
        return 0

    for skill, cmd in sorted(set(missing_cmds)):
        print(f"MISSING COMMAND  {skill}: {cmd}", file=sys.stderr)
    for skill, cmd, fl in sorted(set(missing_flags)):
        print(f"MISSING FLAG     {skill}: '{cmd}' has no {fl}", file=sys.stderr)
    print("\nUpdate the skill, or the CLI, so they agree.", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
