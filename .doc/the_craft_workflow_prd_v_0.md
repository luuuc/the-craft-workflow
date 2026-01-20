# The Craft Workflow – PRD (v0)

**Subtitle:** Deliberate judgment before execution

---

## 1. Purpose (v0)

The Craft Workflow (v0) is a **local CLI that enforces judgment discipline** when working on a piece of work.

It does this by:
- making intent explicit
- refusing invalid workflow transitions
- requiring human confirmation to move forward

v0 is intentionally small. It is designed to be useful **without any AI tools installed**.

---

## 2. Problem Statement

When building with or without AI:
- decisions are implicit and easily drift
- work advances prematurely
- context is lost between steps
- humans end up re-deciding the same things

The problem is not lack of intelligence or tooling.
The problem is **lack of enforcement around judgment**.

---

## 3. Goals (v0)

### Primary Goals
- Make intent explicit and durable
- Prevent accidental or premature progression
- Keep the human as the sole authority
- Be predictable, boring, and inspectable

### Secondary Goals
- Require almost no learning curve
- Work fully offline
- Be usable within one week of installation

---

## 4. Non-Goals (v0)

The Craft Workflow (v0) will **not**:
- make decisions
- suggest next steps
- manage tasks or backlogs
- optimize prompts
- require AI tools
- expose internal workflow theory to users

---

## 5. Target User

v0 is for individual engineers or founders who:
- think carefully before building
- want guardrails, not automation
- prefer clarity over speed
- are comfortable with the CLI

---

## 6. Workflow Model (v0)

v0 implements a **minimal state machine** with three internal states.

These states are **not user-facing concepts**.

```
thinking → building → shipping
```

The workflow always begins in `thinking`.

---

## 7. Commands (v0)

v0 exposes **five commands** only.

```bash
craft start "<intent>"
craft think
craft accept ["note"]
craft reject ["note"]
craft ship
```

Optional but recommended:

```bash
craft status
```

All commands are verbs. No flags are required for normal use.

---

## 8. Command Semantics

### craft start

- Creates a new workflow
- Stores the initial intent
- Sets internal state to `thinking`

---

### craft think

- Represents the deliberation loop
- May invoke an external reviewer (if present)
- Otherwise opens the workflow file for manual reasoning

This command **never advances state**.

---

### craft accept

- Signals that the user is aligned enough to proceed
- Freezes current intent and rationale
- Advances state from `thinking` → `building`

An optional note may be provided.

---

### craft reject

- Records a rejection or concern
- Keeps the workflow in `thinking`

Rejection never advances state.

---

### craft ship

- Packages the work
- Finalizes the workflow
- Advances state from `building` → `shipping`

Shipping always requires explicit human action.

---

### craft status

- Displays the current state
- Lists allowed next actions
- Shows the recorded intent and notes

This command is essential for debugging and trust.

---

## 9. State Enforcement

- Workflow state is authoritative
- State transitions are validated explicitly
- Invalid transitions are refused with clear errors

The tool never infers state from AI output or file contents alone.

---

## 10. State Storage (v0)

v0 uses **a single authoritative file**:

```
.craft/workflow.md
```

### Format

- Markdown with YAML front matter
- Human-readable
- Machine-parseable

Example:

```markdown
---
state: thinking
version: 1
checksum: abc123
---

# Intent
Add anomaly detection to Kamal metrics.

## Notes
- Initial concern about false positives
```

The checksum is used to detect manual tampering.

---

## 11. Failure Modes (v0)

### Interrupted Command
- State remains unchanged
- Partial output is discarded

### Manual File Edit
- Checksum mismatch detected
- Workflow refuses to advance

### External Tool Failure
- Error is reported
- State is unchanged

### Disk / IO Error
- Command fails fast
- No partial state is written

v0 prioritizes correctness over recovery.

---

## 12. Integrations (v0)

v0 works **without any integrations**.

If present:
- External reviewers may be invoked during `craft think`

Integrations are accelerators, not requirements.

---

## 13. Technical Stack (v0)

- Language: Go
- Distribution: single static binary
- Runtime: local CLI only
- No daemon, no server

---

## 14. Testing Strategy (v0)

- Unit tests for every valid state transition
- Unit tests for every invalid transition
- Tests for interrupted commands
- Tests for manual file modification

The state machine is the primary test surface.

---

## 15. Success Criteria (v0)

v0 is successful when:
- users cannot accidentally skip thinking
- workflow state is always clear
- the tool feels boring and safe
- the user trusts it to say "no"

---

## 16. Status

This PRD defines **v0 only**.

Future versions are intentionally unspecified.
