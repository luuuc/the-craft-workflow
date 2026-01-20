# Card: State Machine and Transitions

## Summary

Implement the core state machine that enforces workflow transitions.

## Scope

### In Scope
- Define three states: `thinking`, `building`, `shipping`
- Define valid transitions:
  - `thinking` → `building` (via accept)
  - `building` → `shipping` (via ship)
  - `any` → `none` (via reset)
- Reject invalid transitions with clear error messages
- Expose state validation as a testable module

### Out of Scope
- File I/O (separate card)
- CLI parsing (separate card)
- Checksum logic (part of storage card)

## Tasks

- [ ] Define `State` type with three values
- [ ] Define `Transition` type
- [ ] Implement `ValidateTransition(from, to State) error`
- [ ] Implement `NextValidActions(current State) []string`
- [ ] Write tests for all valid transitions
- [ ] Write tests for all invalid transitions

## Acceptance Criteria

- `thinking` can only advance to `building`
- `building` can only advance to `shipping`
- `shipping` cannot advance (terminal state)
- Invalid transitions return descriptive errors
- 100% test coverage on the state machine
