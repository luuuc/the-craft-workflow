package main

import (
	"fmt"
	"os"

	"craft/cmd"
)

const version = "0.5.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}

	switch args[0] {
	case "--help", "-h", "help":
		printHelp()
		return 0
	case "--version", "-v":
		fmt.Println("craft version " + version)
		return 0
	case "start":
		return cmd.Start(args[1:])
	case "think":
		return cmd.Think(args[1:])
	case "accept":
		return cmd.Accept(args[1:])
	case "reject":
		return cmd.Reject(args[1:])
	case "shape":
		return cmd.Shape(args[1:])
	case "approve":
		return cmd.Approve(args[1:])
	case "revise":
		return cmd.Revise(args[1:])
	case "ship":
		return cmd.Ship(args[1:])
	case "status":
		return cmd.Status(args[1:])
	case "reset":
		return cmd.Reset(args[1:])
	case "init":
		return cmd.Init(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", args[0])
		fmt.Fprintln(os.Stderr, "Run 'craft --help' for usage.")
		return 1
	}
}

func printHelp() {
	fmt.Print(`craft - deliberate judgment before execution

Usage:
  craft <command> [arguments]

Commands:
  start "<intent>"   Begin a new workflow with the given intent
  think              Review the current workflow state
  accept [note]      Confirm alignment and advance to shaping
  reject [note]      Record a concern, stay in thinking
  shape              Show shaping status
  shape --generate   Generate pitch and cards using AI
  approve            Approve structure and advance to building
  revise "note"      Record a concern during shaping
  ship               Finalize the workflow
  status             Show current state and valid actions
  reset              Abandon current workflow
  init [flags]       Copy AI integration templates

Accept flags:
  --skip-shaping     Skip shaping phase, advance directly to building

Init flags:
  --claude           Copy Claude Code templates (CLAUDE.md, .claude/commands/)
  --cursor           Copy Cursor rules (.cursorrules)
  --all              Copy all templates

Workflow states: thinking → shaping → building → shipped
State is stored in .craft/workflow.md
`)
}
