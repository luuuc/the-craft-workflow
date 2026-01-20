# Card: craft status

## Summary

Implement the `craft status` command that displays current workflow state and valid actions.

## Scope

### In Scope
- Display current state
- Display intent
- Display notes
- List valid next actions
- Warn if checksum mismatch detected

### Out of Scope
- Modifying state
- Refusing to display on checksum mismatch (warn only)

## Tasks

- [ ] Implement `status` subcommand
- [ ] Load workflow file
- [ ] Validate checksum (warn on mismatch, don't fail)
- [ ] Display state, intent, notes
- [ ] Display valid actions from current state
- [ ] Handle missing workflow gracefully
- [ ] Write tests for output format
- [ ] Write tests for checksum warning

## CLI Behavior

```
$ craft status
State: thinking
Intent: Add rate limiting to API

Notes:
- Need to consider rate limit headers

Actions: accept, reject, reset

$ craft status
Warning: Workflow file modified externally. State may be inconsistent.

State: thinking
Intent: Add rate limiting to API
...

$ craft status
No workflow found. Run 'craft start' to begin.
```

## Acceptance Criteria

- Displays state, intent, notes, valid actions
- Warns on checksum mismatch but still displays
- Clear message when no workflow exists
- Read-only operation
