package internal

import "runtime"

var _stackPool = NewPool(func() *Stack {
	return &Stack{
		storage: make([]uintptr, 64),
	}
})

// Stack is a captured stack trace.
type Stack struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// Depth specifies how deep of a stack trace should be captured.
type Depth int

const (
	// First captures only the first frame.
	First Depth = iota

	// Full captures the entire call stack, allocating more
	// storage for it if needed.
	Full
)

// Capture captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// Capture.
//
// The caller must call Free on the returned stacktrace after using it.
func Capture(skip int, depth Depth) *Stack {
	stack := _stackPool.Get()

	switch depth {
	case First:
		stack.pcs = stack.storage[:1]
	case Full:
		stack.pcs = stack.storage
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	if depth == Full {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *Stack) Free() {
	st.frames = nil
	st.pcs = nil
	_stackPool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *Stack) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *Stack) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}
