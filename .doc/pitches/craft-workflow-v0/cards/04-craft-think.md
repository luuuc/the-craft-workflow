# Card: craft think

## Summary

Implement the `craft think` command that displays the current workflow state for deliberation.

## Scope

### In Scope
- Display current intent and notes
- Display current state
- Show what actions are available
- Work only in `thinking` state (informational in other states)

### Out of Scope
- External reviewer invocation (future version)
- Opening files in editor (future version)

## Tasks

- [ ] Implement `think` subcommand
- [ ] Load current workflow
- [ ] Display intent prominently
- [ ] Display any recorded notes
- [ ] Display current state
- [ ] Display valid next actions
- [ ] Handle missing workflow gracefully
- [ ] Write tests for output format

## CLI Behavior

```
$ craft think
# Intent
Add rate limiting to API

## Notes
(none)

State: thinking
Actions: accept, reject, reset

$ craft think
Error: No workflow found. Run 'craft start' to begin.
```

## Acceptance Criteria

- Displays intent and notes clearly
- Shows current state
- Lists valid actions
- Does not modify state (read-only)
- Clear error when no workflow exists
