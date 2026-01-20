# Craft Workflow v0

**Appetite:** Small batch (2 weeks)

## Overview

A CLI that enforces judgment discipline with five commands and one file.

## Documents

- [pitch.md](pitch.md) - Full pitch with problem, solution, rabbit holes, no-gos

## Cards

Work is organized into cards to be completed in sequence:

| # | Card | Summary |
|---|------|---------|
| 00 | [CLI Skeleton](cards/00-cli-skeleton.md) | Entry point, routing, help |
| 01 | [State Machine](cards/01-state-machine.md) | Core state transitions |
| 02 | [Workflow Storage](cards/02-workflow-storage.md) | File I/O, checksums |
| 03 | [craft start](cards/03-craft-start.md) | Begin workflow |
| 04 | [craft think](cards/04-craft-think.md) | Review state |
| 05 | [craft accept/reject](cards/05-craft-accept-reject.md) | Decision commands |
| 06 | [craft ship](cards/06-craft-ship.md) | Finalize workflow |
| 07 | [craft status](cards/07-craft-status.md) | Display state |
| 08 | [craft reset](cards/08-craft-reset.md) | Abandon workflow |

## Suggested Order

1. **00 + 01 + 02** can be built in parallel (no dependencies)
2. **03** depends on 01, 02
3. **04, 05, 06, 07, 08** depend on 01, 02, 03

## Source PRD

See [the_craft_workflow_prd_v_0.md](../../the_craft_workflow_prd_v_0.md)
