# The Craft Workflow

Deliberate judgment before execution.

## The Problem

Decisions drift. Work advances before thinking is done. Nothing enforces the first decision.

## The Solution

A CLI that says "no" when you try to skip thinking.

```
thinking → building → shipped
```

You can't ship without building. You can't build without thinking first.

## Commands

```
craft start "<intent>"   Begin with explicit intent
craft think [--review]   Review where you are (optionally invoke reviewer)
craft accept [note]      Confirm alignment, advance to building
craft reject [note]      Record concern, stay in thinking
craft ship               Finalize the work
craft status             Show current state and valid actions
craft reset              Abandon current workflow
craft init [flags]       Copy AI integration templates
```

## Installation

```bash
go install github.com/luuuc/craft@latest
```

Or download the binary from [Releases](https://github.com/luuuc/craft/releases).

## Example

```
$ craft start "Add rate limiting to API"
Workflow started. State: thinking

$ craft ship
Error: Invalid transition. Current state: thinking
Must accept before shipping.

$ craft accept "Decided on token bucket algorithm"
Intent frozen. State: building

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
3. `craft accept` - Freeze intent, start implementation
4. `craft ship` - Finalize when complete

We eat our own dog food.

## Development

```bash
go test ./...
go build -o craft .
```

## License

MIT
