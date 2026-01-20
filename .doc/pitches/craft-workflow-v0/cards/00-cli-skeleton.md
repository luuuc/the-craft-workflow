# Card: CLI Skeleton and Help

## Summary

Implement the CLI entry point with subcommand routing and help text.

## Scope

### In Scope
- Main entry point
- Subcommand routing (start, think, accept, reject, ship, status, reset)
- `--help` flag and help text
- `--version` flag
- Unknown command handling

### Out of Scope
- Individual command implementations (separate cards)
- Configuration files

## Tasks

- [ ] Set up Go module (`craft`)
- [ ] Implement main.go with subcommand dispatch
- [ ] Implement help text
- [ ] Implement version flag
- [ ] Handle unknown commands gracefully
- [ ] Write tests for command routing
- [ ] Write tests for help output

## Help Text

```
craft - deliberate judgment before execution

Usage:
  craft <command> [arguments]

Commands:
  start "<intent>"   Begin a new workflow with the given intent
  think              Review the current workflow state
  accept [note]      Confirm alignment and advance to building
  reject [note]      Record a concern, stay in thinking
  ship               Finalize the workflow
  status             Show current state and valid actions
  reset              Abandon current workflow

State is stored in .craft/workflow.md
```

## Acceptance Criteria

- `craft --help` displays help text
- `craft --version` displays version
- Unknown commands show error and help hint
- Subcommands are routed correctly
- Clean exit codes (0 success, 1 error)
