# Craft Workflow Rules for Cursor

This project uses craft for deliberate judgment before execution.

## Before any code changes

Check `.craft/workflow.md` for current state. Run `craft status` if unsure.

## State behaviors

### thinking
- Help deliberate. Ask questions. Challenge assumptions.
- DO NOT write implementation code.
- Suggest: `craft accept` when thinking is complete
- Suggest: `craft reject "reason"` if concerns arise

### building
- Implement the frozen intent exactly as decided.
- Stay within the agreed scope.
- DO NOT introduce new features or changes beyond the intent.
- Suggest: `craft ship` when implementation is complete

### shipped
- Work is complete.
- Suggest: `craft reset` to start new work

## Commands

```
craft start "<intent>"   Begin with explicit intent
craft accept [note]      Freeze intent, advance to building
craft reject [note]      Record concern, stay in thinking
craft ship               Finalize the work
craft status             Show current state
craft reset              Abandon workflow
```

## Key principle

Think first. Build second. The workflow enforces this order.
