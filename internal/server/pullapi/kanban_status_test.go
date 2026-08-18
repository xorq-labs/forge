package pullapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kenn.io/forge/internal/db"
)

// TestValidKanbanStatesMatchesKanbanStatuses ties the statuses the kanban
// mutation accepts to db.KanbanStatuses, which shell completion offers. The
// list filter passes its value straight to the query layer, so this guards the
// mutation boundary: without it a new status could become settable without
// being offered, or offered without being settable.
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
