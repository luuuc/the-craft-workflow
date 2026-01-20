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
craft think              Review where you are
craft accept [note]      Confirm alignment, advance to building
craft reject [note]      Record concern, stay in thinking
craft ship               Finalize the work
craft status             Show current state and valid actions
craft reset              Abandon current workflow
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

Markdown with YAML front matter. Human-readable. Machine-parseable. A checksum detects tampering.

## What This Tool Does Not Do

- No task management
- No AI suggestions
- No daemon or background process
- No network calls
- No configuration files

## Development

```bash
go test ./...
go build -o craft .
```

## License

MIT
