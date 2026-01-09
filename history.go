package main

type History struct {
	states  []string
	current int
	maxSize int
}

func NewHistory(maxSize int) *History {
	return &History{
		states:  make([]string, 0, maxSize),
		current: -1,
		maxSize: maxSize,
	}
}

func (h *History) Push(state string) {
	if h.current < len(h.states)-1 {
		h.states = h.states[:h.current+1]
	}
	h.states = append(h.states, state)
	if len(h.states) > h.maxSize {
		h.states = h.states[1:]
	} else {
		h.current++
	}
}

func (h *History) Undo() *string {
	if h.current < 0 || len(h.states) == 0 {
		return nil
	}
	state := h.states[h.current]
	h.current--
	return &state
}

func (h *History) Redo() *string {
	if h.current >= len(h.states)-1 {
		return nil
	}
	h.current++
	state := h.states[h.current]
	return &state
}
