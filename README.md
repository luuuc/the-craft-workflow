# The Craft Workflow

Deliberate judgment before execution.

## The Problem

Decisions drift. Work advances before thinking is done. Nothing enforces the first decision.

## The Solution

A CLI that says "no" when you try to skip thinking.

```
thinking → shaping → building → shipped
```

You can't ship without building. You can't build without shaping (or skipping it). You can't shape without thinking first.

## Commands

```
craft start "<intent>"   Begin with explicit intent
craft think [--review]   Review where you are (optionally invoke reviewer)
craft accept [note]      Confirm alignment, advance to shaping
craft reject [note]      Record concern, stay in thinking
craft shape              Show shaping status
craft shape --generate   Generate pitch and cards via AI
craft approve            Approve structure, advance to building
craft revise "note"      Record concern during shaping
craft ship               Finalize the work
craft status             Show current state and valid actions
craft reset              Abandon current workflow
craft init [flags]       Copy AI integration templates
```

Use `craft accept --skip-shaping` to go directly to building for simple tasks.

## Installation

```bash
go install github.com/luuuc/craft@latest
```

Or download the binary from [Releases](https://github.com/luuuc/craft/releases).

## Example

```
$ craft start "Add rate limiting to API"
Workflow started. State: thinking

$ craft accept "Token bucket algorithm"
Intent frozen. State: shaping

$ craft shape
Shaping: Add rate limiting to API
Structure: (none)
Next: craft shape --generate OR create .craft/pitch.md

$ craft shape --generate
Generating via AI...
Created:
  .craft/pitch.md
  .craft/cards/01-rate-limiter.md
  .craft/cards/02-middleware.md
Next: craft approve

$ craft approve
Structure approved. State: building

$ craft ship
Workflow complete. State: shipped
```

The CLI is boring. That's the point.

## State Storage

All state lives in a single file:

```
.craft/workflow.md
```

Markdown with YAML front matter. Human-readable. Machine-parseable. Includes timestamps and history for accountability. A checksum detects tampering.

## What This Tool Does Not Do

- No task management
- No daemon or background process
- No configuration files

## AI Integration

craft works with AI coding assistants. Use `craft init` to copy integration templates:

```bash
craft init --claude    # Claude Code (CLAUDE.md + slash command)
craft init --cursor    # Cursor AI (.cursorrules)
craft init --all       # All templates
```

Or see [templates/INTEGRATION.md](templates/INTEGRATION.md) for manual setup.

## Optional Review

Invoke external reviewers during thinking:

```
craft think --review          # Auto-detect best reviewer
craft think --review=ai       # Use AI reviewer
craft think --review=council  # Use council-cli
```

Environment variables for AI review:
- `CRAFT_AI_API_KEY` — Required
- `CRAFT_AI_MODEL` — Optional (default: gpt-4o-mini)
- `CRAFT_AI_BASE_URL` — Optional (default: OpenAI)

For local models (Ollama):
```bash
export CRAFT_AI_BASE_URL=http://localhost:11434/v1
export CRAFT_AI_MODEL=llama3
export CRAFT_AI_API_KEY=unused  # Required but not validated by Ollama
```

Without configuration, falls back to self-review prompts.

## Development Workflow

This project is built using craft.

Every feature follows the workflow:
1. `craft start "<intent>"` - Begin with explicit intent
2. `craft think` - Deliberate before building
3. `craft accept` - Freeze intent, advance to shaping
4. `craft shape --generate` or create pitch manually
5. `craft approve` - Approve structure, start implementation
6. `craft ship` - Finalize when complete

We eat our own dog food.

## Development

```bash
go test ./...
go build -o craft .
```

## License

MIT
