package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.kenn.io/forge/internal/cli/serve"
)

func newCompletionTestRoot(t *testing.T, stdout io.Writer) *cobra.Command {
	t.Helper()
	if stdout == nil {
		stdout = io.Discard
	}
	return newRootCommand(cliOptions{
		Stdin:     strings.NewReader(""),
		Stdout:    stdout,
		Stderr:    io.Discard,
		RunServer: func(serve.Options) error { return nil },
	})
}

func TestAgentFlagCompletion(t *testing.T) {
	for _, name := range []string{"run", "install", "uninstall"} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			cmd := mustFindCommand(newCompletionTestRoot(t, nil), "agent-hook", name)

			completion, ok := cmd.GetFlagCompletionFunc("agent")
			require.True(t, ok)

			got, directive := completion(cmd, nil, "")

			assert.Contains(got, cobra.Completion("claude"))
			assert.Contains(got, cobra.Completion("codex"))
			assert.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
		})
	}
}

func TestOutputFlagCompletion(t *testing.T) {
	root := newCompletionTestRoot(t, nil)

	completion, ok := root.GetFlagCompletionFunc("output")
	require.True(t, ok)

	got, directive := completion(root, nil, "")

	assert.ElementsMatch(t, []cobra.Completion{"json", "yaml", "jsonl"}, got)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestKanbanFlagCompletion(t *testing.T) {
	cmd := mustFindCommand(newCompletionTestRoot(t, nil), "pulls")

	completion, ok := cmd.GetFlagCompletionFunc("kanban")
	require.True(t, ok)

	got, directive := completion(cmd, nil, "")

	assert.ElementsMatch(t, []cobra.Completion{"new", "reviewing", "waiting", "awaiting_merge"}, got)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestProviderArgCompletion(t *testing.T) {
	for _, name := range []string{"pulls", "issues"} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			cmd := mustFindCommand(newCompletionTestRoot(t, nil), name, "get")
			require.NotNil(t, cmd.ValidArgsFunction)

			got, directive := cmd.ValidArgsFunction(cmd, nil, "")

			assert.ElementsMatch([]cobra.Completion{"github", "gitlab", "forgejo", "gitea"}, got)
			assert.Equal(cobra.ShellCompDirectiveNoFileComp, directive)

			owner, ownerDirective := cmd.ValidArgsFunction(cmd, []string{"github"}, "")
			assert.Empty(owner)
			assert.Equal(cobra.ShellCompDirectiveNoFileComp, ownerDirective)
		})
	}
}

func TestAPIMethodArgCompletion(t *testing.T) {
	assert := assert.New(t)
	cmd := mustFindCommand(newCompletionTestRoot(t, nil), "api")
	require.NotNil(t, cmd.ValidArgsFunction)

	got, directive := cmd.ValidArgsFunction(cmd, nil, "")

	assert.ElementsMatch(
		[]cobra.Completion{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		got,
	)
	assert.Equal(cobra.ShellCompDirectiveNoFileComp, directive)

	path, pathDirective := cmd.ValidArgsFunction(cmd, []string{"GET"}, "")
	assert.Empty(path)
	assert.Equal(cobra.ShellCompDirectiveNoFileComp, pathDirective)
}

func TestConfigKeyArgCompletion(t *testing.T) {
	assert := assert.New(t)
	cmd := mustFindCommand(newCompletionTestRoot(t, nil), "config", "read")
	require.NotNil(t, cmd.ValidArgsFunction)

	got, directive := cmd.ValidArgsFunction(cmd, nil, "")

	assert.ElementsMatch([]cobra.Completion{"port"}, got)
	assert.Equal(cobra.ShellCompDirectiveNoFileComp, directive)
}

func TestConfigReadKeysAllResolve(t *testing.T) {
	require := require.New(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	require.NoError(os.WriteFile(cfgPath, []byte("port = 9123\n"), 0o644))

	for _, key := range ConfigReadKeys {
		var stdout bytes.Buffer
		require.NoError(readConfigValue(cfgPath, key, &stdout), "completed key %q must resolve", key)
		require.NotEmpty(stdout.String())
	}
}

func TestShellCompletionOmitsHiddenCommands(t *testing.T) {
	assert := assert.New(t)
	var stdout bytes.Buffer
	root := newCompletionTestRoot(t, &stdout)
	root.SetArgs([]string{"__complete", ""})

	require.NoError(t, root.Execute())

	assert.Contains(stdout.String(), "pulls")
	assert.Contains(stdout.String(), "completion")
	assert.NotContains(stdout.String(), "pty-owner")
}

func TestCompletionCommandGeneratesShellScripts(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish"} {
		t.Run(shell, func(t *testing.T) {
			var stdout bytes.Buffer
			root := newCompletionTestRoot(t, &stdout)
			root.SetArgs([]string{"completion", shell})

			require.NoError(t, root.Execute())

			assert.Contains(t, stdout.String(), "kenn-forge")
			assert.Greater(t, stdout.Len(), 1000)
		})
	}
}
