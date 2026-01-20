# Craft Workflow Integration

This project uses **craft** for deliberate judgment before execution.

## Before ANY code changes

1. Run `craft status` to check current workflow state
2. Respect the state:
   - **thinking**: Help deliberate. Do NOT write implementation code.
   - **building**: Implement the frozen intent. Stay within decided scope.
   - **shipped**: Work is complete. Suggest `craft reset` for new work.

## State meanings

| State | Purpose | Your role |
|-------|---------|-----------|
| `thinking` | Decision in progress | Ask questions. Challenge assumptions. Discuss alternatives. No implementation code. |
| `building` | Intent is frozen | Implement exactly as decided. Do not introduce new intent. |
| `shipped` | Work is complete | No modifications without a new workflow. |

## Commands to suggest

| When | Suggest |
|------|---------|
| Thinking is complete | `craft accept [note]` |
| Concerns arise during thinking | `craft reject "reason"` |
| Implementation is complete | `craft ship` |
| Need to check state | `craft status` |
| Start fresh | `craft reset` |

## Workflow file

State lives in `.craft/workflow.md`. Read it to understand:
- Current state
- Original intent
- Notes from deliberation

## Example session

```
User: Let's add rate limiting
AI: I see no craft workflow is active. Let's start one:
    craft start "Add rate limiting to API"

User: [runs craft start]
AI: Good, we're in thinking state. Before coding:
    - What rate limit algorithm? Token bucket? Sliding window?
    - Per-user or global limits?
    - What happens when limits are hit?

User: Token bucket, per-user, return 429
AI: That's a clear decision. Ready to proceed?
    Run: craft accept "Token bucket, per-user, 429 on limit"

User: [runs craft accept]
AI: Intent frozen. Now I'll implement token bucket rate limiting...
```

The workflow ensures we think before we build.
