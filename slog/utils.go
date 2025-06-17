package slog

import (
	"runtime"
	"strings"
)

func pkgFilePath(frame *runtime.Frame) string {
	if frame.Function == "" {
		return ""
	}
	idx := strings.IndexByte(frame.Function, '(')
	if idx != -1 {
		return frame.Function[:idx-1]
	}
	idx = strings.LastIndexByte(frame.Function, '.')
	if idx != -1 {
		return frame.Function[:idx]
	}
	return frame.Function
}
