package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"go.kenn.io/forge/internal/db"
	"go.kenn.io/forge/internal/platform"
	"go.kenn.io/kit/agenthook"
)

// registerCompletions wires shell completion for enum-valued flags and
// positional arguments onto the assembled command tree. Every completer
// returns ShellCompDirectiveNoFileComp; flags that take paths, such as
// --config and --binary, keep Cobra's default file completion.
func registerCompletions(root *cobra.Command) {
	registerOutputCompletion(root)
	registerKanbanCompletion(mustFindCommand(root, "pulls"))
	registerProviderArgCompletion(mustFindCommand(root, "pulls", "get"))
	registerProviderArgCompletion(mustFindCommand(root, "issues", "get"))
	registerAPIMethodCompletion(mustFindCommand(root, "api"))
	registerConfigKeyCompletion(mustFindCommand(root, "config", "read"))
	for _, name := range []string{"run", "install", "uninstall"} {
		registerAgentCompletion(mustFindCommand(root, "agent-hook", name))
	}
}

// mustFindCommand resolves a subcommand by its name path.
// Panics if the command doesn't exist on the tree (programming error).
func mustFindCommand(root *cobra.Command, names ...string) *cobra.Command {
	current := root
	for _, name := range names {
		var next *cobra.Command
		for _, candidate := range current.Commands() {
			if candidate.Name() == name {
				next = candidate
				break
			}
		}
		if next == nil {
			panic(fmt.Sprintf("finding command %s under %s for completion", name, current.Name()))
		}
		current = next
	}
	return current
}

// registerAgentCompletion registers shell completion for the --agent flag.
// Panics if the flag doesn't exist on the command (programming error).
func registerAgentCompletion(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc("agent", func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		profiles := agenthook.Profiles()
		agents := make([]cobra.Completion, 0, len(profiles))
		for _, profile := range profiles {
			agents = append(agents, string(profile.Agent))
		}
		return agents, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(fmt.Sprintf("registering agent completion for %s: %v", cmd.Name(), err))
	}
}

// registerOutputCompletion registers shell completion for the --output flag.
// Panics if the flag doesn't exist on the command (programming error).
func registerOutputCompletion(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc("output", func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []cobra.Completion{"json", "yaml", "jsonl"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(fmt.Sprintf("registering output completion for %s: %v", cmd.Name(), err))
	}
}

// registerKanbanCompletion registers shell completion for the --kanban flag.
// Panics if the flag doesn't exist on the command (programming error).
func registerKanbanCompletion(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc("kanban", func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []cobra.Completion{
			string(db.KanbanStatusNew),
			string(db.KanbanStatusReviewing),
			string(db.KanbanStatusWaiting),
			string(db.KanbanStatusAwaitingMerge),
		}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(fmt.Sprintf("registering kanban completion for %s: %v", cmd.Name(), err))
	}
}

// registerProviderArgCompletion completes the PROVIDER positional on
// "get PROVIDER OWNER NAME NUMBER" commands from the canonical provider list.
func registerProviderArgCompletion(cmd *cobra.Command) {
	cmd.ValidArgsFunction = func(_ *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		kinds := platform.Kinds()
		providers := make([]cobra.Completion, 0, len(kinds))
		for _, kind := range kinds {
			providers = append(providers, string(kind))
		}
		return providers, cobra.ShellCompDirectiveNoFileComp
	}
}

// registerAPIMethodCompletion completes the METHOD positional on
// "api METHOD PATH [body...]".
func registerAPIMethodCompletion(cmd *cobra.Command) {
	cmd.ValidArgsFunction = func(_ *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []cobra.Completion{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		}, cobra.ShellCompDirectiveNoFileComp
	}
}

// registerConfigKeyCompletion completes the KEY positional on "config read".
func registerConfigKeyCompletion(cmd *cobra.Command) {
	cmd.ValidArgsFunction = func(_ *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return append([]cobra.Completion(nil), ConfigReadKeys...), cobra.ShellCompDirectiveNoFileComp
	}
}
