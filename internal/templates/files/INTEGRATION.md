# Integrating craft with AI Tools

This document explains how to integrate craft workflow with any AI coding assistant.

## Principles

1. **AI reads state from `.craft/workflow.md`**
   The workflow file is markdown with YAML front matter. Parse it or run `craft status`.

2. **AI behavior changes based on state**
   Each state has specific allowed behaviors. See State Behaviors below.

3. **AI never advances state**
   Only the human runs craft commands. AI suggests, human decides.

4. **AI suggests appropriate craft commands**
   Guide the user through the workflow with timely suggestions.

## State Machine

```
thinking → building → shipped
    ↓         ↓         ↓
  reset     reset     reset
```

- `thinking`: Initial state. Deliberation phase.
- `building`: Intent frozen. Implementation phase.
- `shipped`: Terminal state. Work complete.
- Any state can `reset` to abandon the workflow.

## State Behaviors

### thinking

**Purpose**: Deliberate before implementing.

**AI should**:
- Ask clarifying questions
- Challenge assumptions
- Discuss alternatives
- Help refine the intent
- Point out edge cases

**AI should NOT**:
- Write implementation code
- Suggest specific implementations (discuss approaches, don't code them)
- Advance the workflow

**Suggest when ready**:
- `craft accept [note]` - Freeze intent and start building
- `craft reject [note]` - Record a concern, continue thinking

### building

**Purpose**: Implement the frozen intent.

**AI should**:
- Implement exactly as decided in thinking phase
- Stay within the agreed scope
- Reference the notes from deliberation

**AI should NOT**:
- Introduce new features or scope
- Suggest changes to the intent
- Start new unrelated work

**Suggest when complete**:
- `craft ship` - Finalize the work

### shipped

**Purpose**: Mark work as complete.

**AI should**:
- Confirm the workflow is complete
- Refuse modifications without a new workflow

**Suggest for new work**:
- `craft reset` followed by `craft start "<new intent>"`

## Workflow File Format

Location: `.craft/workflow.md`

```yaml
---
state: thinking
schema_version: 1
checksum: abc12345
---

# Intent
Add rate limiting to API

## Notes
- Considering token bucket algorithm
- Per-user limits
```

## Adapting to Your Tool

1. **Find your tool's system prompt location**
   - Claude Code: `CLAUDE.md` in project root
   - Cursor: `.cursorrules` or `.cursor/rules`
   - Others: Check tool documentation

2. **Add state-checking instructions**
   Tell the AI to run `craft status` before any code changes.

3. **Include behavior rules per state**
   Copy the State Behaviors section above.

4. **Test with a sample workflow**
   ```bash
   craft start "Test integration"
   # Verify AI asks questions instead of coding
   craft accept
   # Verify AI now implements
   craft ship
   ```

## Example Integration Prompt

```
Before any code changes, run `craft status` to check workflow state.

If state is "thinking":
- Help deliberate. Ask questions. No implementation code.
- Suggest `craft accept` when ready.

If state is "building":
- Implement the frozen intent.
- Suggest `craft ship` when done.

If state is "shipped":
- Work is complete.
- Suggest `craft reset` for new work.
```

This minimal prompt captures the essential behavior. Adapt as needed for your tool.
