# Pitch: Craft Workflow v0

## Problem

When building software—with or without AI tools—decisions drift.

Work advances before thinking is done. Intent gets lost between steps. Humans re-decide the same things because nothing enforced the first decision.

The problem isn't lack of intelligence or tooling. The problem is **lack of enforcement around judgment**.

## Appetite

**Small batch: 2 weeks**

This is a deliberately constrained v0. We're building the minimum surface that enforces judgment discipline. No integrations. No AI. Just a state machine that says "no" when you try to skip thinking.

## Solution

A CLI with five commands and one file.

### The Commands

```
craft start "<intent>"   # Begin with explicit intent
craft think              # Review where you are
craft accept [note]      # Confirm alignment, advance to building
craft reject [note]      # Record concern, stay in thinking
craft ship               # Finalize the work
craft status             # Show state and valid actions
craft reset              # Abandon current workflow
```

### The State Machine

```
thinking → building → shipping
```

That's it. Three states. You can't skip `thinking`. You can't ship without building. The CLI enforces every transition.

### The File

```
.craft/workflow.md
```

Markdown with YAML front matter. Human-readable. Machine-parseable. One source of truth. A checksum detects tampering.

## Rabbit Holes

- **External integrations** — v0 works standalone. No Council, no AI reviewers. Add later.
- **Multiple workflows** — One workflow at a time. Multi-workflow support is a different problem.
- **Config files** — No configuration for v0. Sensible defaults only.
- **Undo/history** — `reset` abandons. No partial undo. Keep it simple.

## No-Gos

- No task management
- No AI suggestions
- No daemon or background process
- No network calls
- No velocity tracking

## Fat Marker Sketch

```
$ craft start "Add rate limiting to API"
Workflow started. State: thinking

$ craft accept
Intent frozen. State: building

$ craft ship
Workflow complete. State: shipped

$ craft start "something else"
Error: Workflow exists. Run 'craft reset' to abandon.
```

The CLI is boring. That's the point.
