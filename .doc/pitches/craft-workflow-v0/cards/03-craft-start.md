# Card: craft start

## Summary

Implement the `craft start` command that begins a new workflow with explicit intent.

## Scope

### In Scope
- Parse `craft start "<intent>"` command
- Refuse if workflow already exists
- Create `.craft/workflow.md` with initial state
- Set state to `thinking`
- Store the provided intent

### Out of Scope
- Overwriting existing workflows (use reset)
- Multi-workflow support

## Tasks

- [ ] Implement `start` subcommand
- [ ] Validate intent is provided and non-empty
- [ ] Check for existing workflow, refuse if present
- [ ] Create new workflow file with `thinking` state
- [ ] Print confirmation message
- [ ] Write tests for successful start
- [ ] Write tests for missing intent
- [ ] Write tests for existing workflow refusal

## CLI Behavior

```
$ craft start "Add rate limiting to API"
Workflow started. State: thinking

$ craft start "Something else"
Error: Workflow already exists. Run 'craft reset' to abandon.

$ craft start
Error: Intent required. Usage: craft start "<intent>"
```

## Acceptance Criteria

- Creates `.craft/workflow.md` when no workflow exists
- Stores intent in the file body
- Sets state to `thinking`
- Refuses when workflow exists with clear error
- Requires non-empty intent
