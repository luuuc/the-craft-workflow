# Card: craft accept / craft reject

## Summary

Implement the `craft accept` and `craft reject` commands for recording decisions during the thinking phase.

## Scope

### In Scope
- `craft accept [note]` - advance from thinking to building
- `craft reject [note]` - record concern, stay in thinking
- Optional note appended to workflow file
- State transition for accept only

### Out of Scope
- Accept/reject in other states (invalid)

## Tasks

- [ ] Implement `accept` subcommand
- [ ] Implement `reject` subcommand
- [ ] Parse optional note argument
- [ ] Validate current state is `thinking`
- [ ] For accept: advance state to `building`
- [ ] For reject: keep state as `thinking`
- [ ] Append note to workflow file if provided
- [ ] Recompute and update checksum
- [ ] Write tests for accept state transition
- [ ] Write tests for reject staying in thinking
- [ ] Write tests for notes being recorded
- [ ] Write tests for wrong state refusal

## CLI Behavior

```
$ craft accept
Intent frozen. State: building

$ craft accept "Decided to use token bucket algorithm"
Intent frozen. State: building

$ craft reject "Need to consider rate limit headers"
Concern recorded. State: thinking

$ craft accept
Error: Invalid transition. Current state: building
```

## Acceptance Criteria

- `accept` advances from `thinking` to `building`
- `reject` stays in `thinking`
- Notes are appended to the workflow file
- Invalid state transitions are refused
- Checksum is updated after modifications
