# Card: craft ship

## Summary

Implement the `craft ship` command that finalizes the workflow.

## Scope

### In Scope
- Advance from `building` to `shipping`
- Mark workflow as complete
- Display final summary

### Out of Scope
- Commit message generation (future)
- PR description generation (future)
- Any external integrations

## Tasks

- [ ] Implement `ship` subcommand
- [ ] Validate current state is `building`
- [ ] Advance state to `shipping`
- [ ] Update checksum
- [ ] Display completion message with summary
- [ ] Write tests for state transition
- [ ] Write tests for wrong state refusal

## CLI Behavior

```
$ craft ship
Workflow complete. State: shipped

Intent: Add rate limiting to API

$ craft ship
Error: Invalid transition. Current state: thinking
Must accept before shipping.

$ craft ship
Error: Invalid transition. Current state: shipped
Workflow already complete.
```

## Acceptance Criteria

- Advances from `building` to `shipping`
- Refuses when not in `building` state
- Displays summary on completion
- Clear error messages for invalid transitions
