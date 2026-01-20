# /craft - Workflow Status

Display the current craft workflow state and suggest appropriate actions.

## Instructions

1. Run `craft status` to get current workflow state
2. Parse the output and display it in a structured format
3. Provide state-specific guidance

## Output Format

```markdown
## Craft Workflow

**State:** [state]
**Intent:** [intent]

**Notes:**
[list of notes or "(none)"]

**Valid Actions:**
[list of valid craft commands]

**Guidance:**
[state-specific guidance]
```

## State-Specific Guidance

### If state is `thinking`

```
You're in the thinking phase. I'll help you deliberate but won't write
implementation code until you run `craft accept`.

Questions to consider:
- Is the intent clear and specific?
- Have alternatives been explored?
- Are there edge cases to handle?
```

### If state is `building`

```
Intent is frozen. I'll implement exactly as decided.
Stay focused on the agreed scope. If new concerns arise,
consider running `craft reset` and starting fresh.
```

### If state is `shipped`

```
This workflow is complete. To start new work:
  craft reset
  craft start "<new intent>"
```

### If no workflow exists

```
No active workflow. To begin deliberate work:
  craft start "<your intent>"

The craft workflow ensures you think before you build.
```

## Example Output

```markdown
## Craft Workflow

**State:** thinking
**Intent:** Add rate limiting to API

**Notes:**
- Considering token bucket algorithm
- Per-user limits preferred

**Valid Actions:**
- `craft accept [note]` - Freeze intent, start building
- `craft reject [note]` - Record concern, continue thinking
- `craft reset` - Abandon workflow

**Guidance:**
You're in the thinking phase. I'll help you deliberate but won't write
implementation code until you run `craft accept`.
```
