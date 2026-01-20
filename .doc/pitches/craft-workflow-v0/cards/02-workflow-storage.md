# Card: Workflow Storage

## Summary

Implement file-based storage for workflow state using a single Markdown file with YAML front matter.

## Scope

### In Scope
- Read/write `.craft/workflow.md`
- YAML front matter parsing (state, schema_version, checksum)
- Markdown body for intent and notes
- Checksum generation and validation (SHA-256, first 8 hex chars)
- Atomic writes (write to temp, rename)
- Error handling for missing/corrupted files

### Out of Scope
- CLI commands (separate cards)
- State machine logic (Card 01)

## Tasks

- [ ] Define `Workflow` struct (state, version, checksum, intent, notes)
- [ ] Implement `LoadWorkflow(path string) (*Workflow, error)`
- [ ] Implement `SaveWorkflow(path string, w *Workflow) error`
- [ ] Implement `ComputeChecksum(content []byte) string`
- [ ] Implement `ValidateChecksum(w *Workflow) error`
- [ ] Implement `EnsureCraftDir() error`
- [ ] Write tests for load/save round-trip
- [ ] Write tests for checksum validation
- [ ] Write tests for tampered file detection

## File Format

```markdown
---
state: thinking
schema_version: 1
checksum: a1b2c3d4
---

# Intent
<user-provided intent>

## Notes
- <note 1>
- <note 2>
```

## Acceptance Criteria

- Workflow loads and saves correctly
- Checksum is computed over the content (excluding checksum line itself)
- Tampered files are detected and rejected
- Missing `.craft/` directory is created automatically
- Atomic writes prevent partial state
