package repo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWindowStore(t *testing.T) {
	ws := NewWindowStore()

	assert.IsType(t, &WindowStore{}, ws)
}

func TestWindowStore_PutGet(t *testing.T) {
	ws := NewWindowStore()
	w := make(Window)

	ws.Put(int64(10), &w)

	assert.Equal(t, &w, ws.Get(int64(10)))
}

func TestWindowStore_Get_Defaults(t *testing.T) {
	ws := NewWindowStore()

	assert.IsType(t, &Window{}, ws.Get(int64(10)))
}

func TestWindowStore_PopOlderThan(t *testing.T) {
	ws := NewWindowStore()
	oldW := make(Window)
	newW := make(Window)
	oldW["hello"] = int64(0)
	newW["hello"] = int64(1)

	ws.Put(int64(1), &oldW)
	ws.Put(int64(2), &oldW)
	ws.Put(int64(3), &newW)
	ws.Put(int64(4), &newW)
	ws.Put(int64(5), &newW)

	windows := ws.PopOlderThan(3)

	for _, window := range windows {
		assert.Equal(t, oldW["hello"], (*window)["hello"])
	}

	assert.Len(t, windows, 2)
	assert.Equal(t, 3, ws.Len())
}

func TestWindowStore_PopAll(t *testing.T) {
	ws := NewWindowStore()
	w := make(Window)

	ws.Put(int64(1), &w)
	ws.Put(int64(2), &w)
	ws.Put(int64(3), &w)
	ws.Put(int64(4), &w)
	ws.Put(int64(5), &w)

	windows := ws.PopAll()

	for _, window := range windows {
		assert.Equal(t, *window, w)
	}

	assert.Len(t, windows, 5)
	assert.Equal(t, 0, ws.Len())
}
