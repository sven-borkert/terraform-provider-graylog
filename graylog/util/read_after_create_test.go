package util

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func setReadAfterCreateRetry(t *testing.T, attempts int, delay time.Duration) {
	t.Helper()

	prevAttempts := readAfterCreateAttempts
	prevDelay := readAfterCreateDelay
	readAfterCreateAttempts = attempts
	readAfterCreateDelay = delay
	t.Cleanup(func() {
		readAfterCreateAttempts = prevAttempts
		readAfterCreateDelay = prevDelay
	})
}

func newTestResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
}

func TestReadAfterCreate(t *testing.T) {
	t.Run("success on first read", func(t *testing.T) {
		setReadAfterCreateRetry(t, 3, time.Millisecond)

		d := newTestResourceData(t)
		calls := 0
		err := ReadAfterCreate(d, nil, "resource-id", func(d *schema.ResourceData, _ interface{}) error {
			calls++
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 1, calls)
		require.Equal(t, "resource-id", d.Id())
	})

	t.Run("retries when read clears id", func(t *testing.T) {
		setReadAfterCreateRetry(t, 4, time.Millisecond)

		d := newTestResourceData(t)
		calls := 0
		err := ReadAfterCreate(d, nil, "resource-id", func(d *schema.ResourceData, _ interface{}) error {
			calls++
			if calls < 3 {
				d.SetId("")
			}
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 3, calls)
		require.Equal(t, "resource-id", d.Id())
	})

	t.Run("returns read error", func(t *testing.T) {
		setReadAfterCreateRetry(t, 3, time.Millisecond)

		d := newTestResourceData(t)
		expectedErr := errors.New("boom")
		err := ReadAfterCreate(d, nil, "resource-id", func(d *schema.ResourceData, _ interface{}) error {
			return expectedErr
		})
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("returns error when still missing", func(t *testing.T) {
		setReadAfterCreateRetry(t, 3, time.Millisecond)

		d := newTestResourceData(t)
		calls := 0
		err := ReadAfterCreate(d, nil, "resource-id", func(d *schema.ResourceData, _ interface{}) error {
			calls++
			d.SetId("")
			return nil
		})
		require.EqualError(t, err, "resource resource-id not found after create")
		require.Equal(t, 3, calls)
	})

	t.Run("validates input", func(t *testing.T) {
		d := newTestResourceData(t)

		err := ReadAfterCreate(d, nil, "", func(d *schema.ResourceData, _ interface{}) error {
			return nil
		})
		require.EqualError(t, err, "id is required")

		err = ReadAfterCreate(d, nil, "resource-id", nil)
		require.EqualError(t, err, "read function is required")
	})
}
