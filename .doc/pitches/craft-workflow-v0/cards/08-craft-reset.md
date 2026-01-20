# Card: craft reset

## Summary

Implement the `craft reset` command that abandons the current workflow.

## Scope

### In Scope
- Delete `.craft/workflow.md`
- Work from any state
- Require confirmation (or `--force` flag)

### Out of Scope
- Archiving old workflows
- Undo/recovery

## Tasks

- [ ] Implement `reset` subcommand
- [ ] Check for existing workflow
- [ ] Prompt for confirmation (unless `--force`)
- [ ] Delete workflow file
- [ ] Display confirmation message
- [ ] Handle missing workflow gracefully
- [ ] Write tests for reset
- [ ] Write tests for force flag
- [ ] Write tests for no workflow

## CLI Behavior

```
$ craft reset
Abandon workflow "Add rate limiting to API"? [y/N] y
Workflow abandoned.

$ craft reset --force
Workflow abandoned.

$ craft reset
No workflow to reset.
```

## Acceptance Criteria

- Deletes workflow file
- Prompts for confirmation by default
- `--force` skips confirmation
- Works from any state
- Clear message when no workflow exists
