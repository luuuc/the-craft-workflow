# The Craft Workflow

**Subtitle:** Deliberate judgment before execution

---

## 1. Overview

The Craft Workflow is a local, human‑in‑the‑loop CLI that enforces deliberate judgment before execution when working with AI tools.

It acts as a **control plane** over existing AI CLIs (Council, Shape, editors, etc.) by:
- capturing intent once
- enforcing explicit human approvals
- preventing semantic drift
- reducing repetitive prompting

The Craft Workflow does not decide, automate judgment, or optimize creativity. It exists to make *good judgment easier to sustain*.

---

## 2. Problem Statement

When working with AI tools:
- context is repeatedly restated
- decisions are implicit and easily drift
- AI tools prematurely suggest next steps
- humans re‑prompt instead of reacting

Even disciplined workflows become laborious due to repetition and cognitive overhead.

The problem is not lack of intelligence, but lack of **workflow rigor** that preserves human judgment while using AI.

---

## 3. Goals

### Primary Goals
- Make judgment explicit and durable
- Reduce repetitive prompting
- Prevent AI‑driven workflow drift
- Keep the human as the sole authority

### Secondary Goals
- Work with any AI CLI or editor
- Require near‑zero learning curve
- Be boring, predictable, and inspectable

---

## 4. Non‑Goals

The Craft Workflow will **not**:
- make decisions on behalf of the user
- manage tasks or backlogs
- track productivity or velocity
- act as an AI agent or daemon
- require a specific editor or LLM

---

## 5. Target Users

Primary users are:
- senior engineers
- founders
- technical leads
- individual contributors who value judgment over speed

Users are comfortable with the CLI and prefer explicit control over automation.

---

## 6. Core Workflow

### High‑Level Flow

```
start → think → shape → build → ship
```

Each phase requires explicit human approval to advance.

---

## 7. Workflow States

### 7.1 Deliberate

**Purpose:** Reach a position the user can defend.

- User provides intent and chosen path
- Council critiques aggressively
- User reacts with: agree / partial / disagree
- Loop continues until alignment is reached

Exit condition: explicit human decision.

---

### 7.2 Decided

**Purpose:** Freeze intent.

Artifacts frozen:
- intent
- rationale
- rejected alternatives

This decision becomes the canonical reference for all downstream work.

---

### 7.3 Shaping

**Purpose:** Turn decision into structure.

Outputs:
- pitch
- cards
- scope boundaries

Council reviews structure for fidelity to the decision.

Exit condition: approved structure.

---

### 7.4 Crafting

**Purpose:** Implement approved structure.

- Editor or Craft tool is invoked
- No new intent is introduced
- Council reviews implementation against decision and structure

Exit condition: approved implementation.

---

### 7.5 Shipping

**Purpose:** Package work for others.

Outputs:
- commit suggestions
- commit messages
- PR description
- tags / labels

User is final approver.

---

## 8. Command Surface (v0)

The CLI exposes a minimal, natural set of commands:

```bash
craft start "<intent>"
craft think
craft agree | craft partial "note" | craft disagree "note"
craft decide

craft shape
craft approve | craft revise "note"

craft build
craft review
craft approve

craft ship
```

Commands map to human intent, not internal state.

---

## 9. Enforcement Model

### Human Authority

Only explicit human commands can:
- advance workflow state
- approve artifacts
- freeze decisions

### AI Constraints

AI tools:
- never advance state
- never approve or decide
- never introduce new intent outside deliberate phase

The CLI enforces all transitions.

---

## 10. State Management

Workflow state is stored locally on disk:

```
.craft/
  state.json
  intent.md
  decision.md
  structure.md
  implementation.md
```

State is authoritative and never inferred from AI output.

---

## 11. Integrations

### Council CLI
- Used for critique and review
- Invoked with full, rehydrated context
- Evaluates fidelity, not next steps

=> available locally at /Users/luc/Developer/luuuc/council-cli, and online at https://github.com/luuuc/council-cli

### Shape CLI
- Used to generate structure
- Receives frozen decision only
- Cannot introduce new intent

=> available locally at /Users/luc/Developer/luuuc/shape-cli, and online at https://github.com/luuuc/shape-cli

### Editors / Craft Tools
- Used for implementation
- Context injected, but workflow remains external

=> Claude Code, Opencode, Windsurf, Cursor, any AI tool capable of using this project.


The Craft Workflow does not require Council, or Shape CLIs to function.
Council and Shape are optional accelerators that reduce effort but do not change the workflow’s rigor.

---

## 12. UX Principles

- No configuration required for v0
- No flags unless strictly necessary
- Clear refusal messages on invalid transitions
- Predictable, boring output

---

## 13. Technical Stack (v0)

- Language: Go
- Distribution: single static binary
- Storage: filesystem
- No daemon, no server

---

## 14. Risks & Mitigations

### Risk: AI drift or over‑helping
**Mitigation:** strict state enforcement outside the AI

### Risk: Workflow feels restrictive
**Mitigation:** minimal commands, natural verbs

### Risk: Feature creep
**Mitigation:** explicit non‑goals and refusal to automate judgment

---

## 15. Success Criteria

The Craft Workflow is successful when:
- users stop re‑prompting AI tools
- decisions are easy to revisit and justify
- AI tools feel quieter, not louder
- judgment remains explicit and respected

---

## 16. Status

v0 focuses exclusively on correctness, clarity, and restraint.

Expansion is intentionally resisted.
