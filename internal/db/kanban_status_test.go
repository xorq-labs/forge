package db

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKanbanEnumTagsMatchKanbanStatuses ties the enum tags that document the
// kanban status field in API schemas to KanbanStatuses. The tags are written by
// hand, so without this guard a new status could reach shell completion and the
// generated schema separately.
func TestKanbanEnumTagsMatchKanbanStatuses(t *testing.T) {
	expected := make([]string, 0, len(KanbanStatuses()))
	for _, status := range KanbanStatuses() {
		expected = append(expected, string(status))
	}

	for _, tc := range []struct {
		value any
		field string
	}{
		{MergeRequest{}, "KanbanStatus"},
		{Issue{}, "WorkflowStatus"},
	} {
		field, ok := reflect.TypeOf(tc.value).FieldByName(tc.field)
		require.True(t, ok, "field %s must exist", tc.field)

		tag, ok := field.Tag.Lookup("enum")
		require.True(t, ok, "field %s must carry an enum tag", tc.field)

		assert.ElementsMatch(t, expected, strings.Split(tag, ","), "enum tag on %s", tc.field)
	}
}
