# rc CLI agent skills

Agent skills for the `rc` CLI. Claude Code discovers these automatically when
working in this repo.

| Skill | Covers |
|---|---|
| [rc-cli-usage](rc-cli-usage/SKILL.md) | Start here. Two APIs, two credentials, `--json`, project selection, which command to reach for |
| [rc-offerings-metadata](rc-offerings-metadata/SKILL.md) | Offering metadata, paywall config, attaching packages |
| [rc-projects-workflow](rc-projects-workflow/SKILL.md) | Projects, entitlements, products, apps, admin |
| [rc-charts-analytics](rc-charts-analytics/SKILL.md) | Charts, A/B price experiments, subscriber lists |

## Keeping them accurate

These document a real command surface, so they go stale silently. When changing
flags or commands, verify against the built binary rather than memory:

```bash
go build -o rc . && ./rc internal offerings update --help
```

The previous versions of these skills documented six flags and commands that did
not exist — every one of which cost an agent a failed command and a recovery loop.
