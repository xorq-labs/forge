package pullapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kenn.io/forge/internal/db"
)

// TestValidKanbanStatesMatchesKanbanStatuses ties the kanban states this API
// accepts to db.KanbanStatuses, which shell completion offers. Without this
// guard a new status could reach one set and not the other, so the CLI would
// complete a value the API rejects or reject a value it completes.
func TestValidKanbanStatesMatchesKanbanStatuses(t *testing.T) {
	expected := make([]string, 0, len(db.KanbanStatuses()))
	for _, status := range db.KanbanStatuses() {
		expected = append(expected, string(status))
	}

	accepted := make([]string, 0, len(validKanbanStates))
	for state := range validKanbanStates {
		accepted = append(accepted, state)
	}

	assert.ElementsMatch(t, expected, accepted)
}
