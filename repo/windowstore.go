package repo

type Window map[string]int64

type WindowStore struct {
	windows map[int64]*Window
}

func NewWindowStore() *WindowStore {
	return &WindowStore{
		windows: make(map[int64]*Window),
	}
}

func (w *WindowStore) Put(windowKey int64, item *Window) {
	w.windows[windowKey] = item
}

func (w *WindowStore) Get(windowKey int64) *Window {
	value, ok := w.windows[windowKey]
	if !ok {
		newWindow := make(Window)
		w.windows[windowKey] = &newWindow
		return &newWindow
	}

	return value
}

func (w *WindowStore) PopOlderThan(maxWindow int64) map[int64]*Window {
	out := make(map[int64]*Window)
	kept := make(map[int64]*Window)

	for k, w := range w.windows {
		if k < maxWindow {
			out[k] = w
		} else {
			kept[k] = w
		}
	}

	w.windows = kept

	return out
}

func (w *WindowStore) PopAll() map[int64]*Window {
	out := w.windows
	w.windows = make(map[int64]*Window)
	return out
}

func (w *WindowStore) Len() int {
	return len(w.windows)
}
