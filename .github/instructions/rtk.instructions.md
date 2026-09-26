---
description: 'Guidance for using RTK (Rust Token Killer) to optimize Copilot token usage'
applyTo: '*'
---

# RTK (Rust Token Killer) Instructions

RTK optimizes GitHub Copilot token usage by filtering command output. Installed and auto-configured in the dev container.

## Quick Start

```bash
rtk --version       # Verify installation
rtk gain            # See cumulative token savings (grows after Copilot commands)
```

## Usage

- **Transparent**: Hook auto-rewrites commands (e.g., `npm run test` → `rtk npm run test`). You never type `rtk` manually.
- **Safe**: If RTK is unavailable or disabled, commands run unchanged.
- **Bypass one command**: `RTK_DISABLED=1 npm run test` (runs unfiltered).
- **Retrieve trimmed output**: Terminal shows `[full output: rtk recall HASH_ID]`. Run it exactly to get full output.

## Diagnostics

```bash
rtk gain                    # Token savings dashboard
rtk config                  # View configuration
rtk init --copilot --show   # Check hook status
RTK_DISABLED=1 <cmd>        # Bypass for single command
rtk recall <HASH>           # Retrieve trimmed output
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `rtk: command not found` | Container not rebuilt. Rebuild devcontainer. |
| `rtk gain` shows no data | Run Copilot commands first, then try again. |
| Commands not rewritten | Restart VS Code Copilot Chat. |
| Need full trimmed output | Run `rtk recall <HASH>` (shown in terminal). |

## References

- [RTK Official Docs](https://www.rtk-ai.app/)
- [Configuration Guide](https://www.rtk-ai.app/docs/getting-started/configuration/)
- [What RTK Optimizes](https://www.rtk-ai.app/docs/resources/what-rtk-covers/)

**Note**: RTK configured at project scope (`.github/hooks/rtk-rewrite.json`).

Instructions used: [copilot-instructions.md, performance-optimization.instructions.md]
