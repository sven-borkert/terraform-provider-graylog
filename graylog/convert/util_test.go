package convert

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapToListSkipsInvalidElements(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"valid-a": map[string]interface{}{"value": "a"},
		"invalid": "not-a-map",
		"valid-b": map[string]interface{}{"value": "b"},
	}

	list := MapToList(data, "id")
	require.Len(t, list, 2)

	ids := make([]string, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		require.True(t, ok)

		id, ok := m["id"].(string)
		require.True(t, ok)
		ids = append(ids, id)
	}

	sort.Strings(ids)
	require.Equal(t, []string{"valid-a", "valid-b"}, ids)
}

func TestListToMapSkipsInvalidElements(t *testing.T) {
	t.Parallel()

	data := []interface{}{
		map[string]interface{}{"id": "one", "value": "a"},
		"not-a-map",
		map[string]interface{}{"id": 123, "value": "b"},
		map[string]interface{}{"id": "two", "value": "c"},
	}

	m := ListToMap(data, "id")
	require.Len(t, m, 2)

	require.Equal(t, map[string]interface{}{"value": "a"}, m["one"])
	require.Equal(t, map[string]interface{}{"value": "c"}, m["two"])
}

func TestInterfaceListToStringListSkipsInvalidElements(t *testing.T) {
	t.Parallel()

	list := InterfaceListToStringList([]interface{}{"foo", 42, "bar", true})
	require.Equal(t, []string{"foo", "bar"}, list)
}
