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

func (w *WindowStore) Put(windowTimeBin int64, item *Window) {
	w.windows[windowTimeBin] = item
}

func (w *WindowStore) Get(windowTimeBin int64) *Window {
	value, ok := w.windows[windowTimeBin]
	if !ok {
		newWindow := make(Window)
		w.windows[windowTimeBin] = &newWindow
		return &newWindow
	}

	return value
}

func (w *WindowStore) PopOlderThan(maxWindowTimeBin int64) map[int64]*Window {
	out := make(map[int64]*Window)
	kept := make(map[int64]*Window)

	for tb, w := range w.windows {
		if tb < maxWindowTimeBin {
			out[tb] = w
		} else {
			kept[tb] = w
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
