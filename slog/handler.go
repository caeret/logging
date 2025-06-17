package slog

import (
	"context"
	"log/slog"

	"github.com/caeret/logging/internal"
)

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	h    slog.Handler
	opts *HandlerOptions
}

func NewHandler(handler slog.Handler, opts ...OptionFunc) slog.Handler {
	o := HandlerOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Handler{
		h:    handler,
		opts: &o,
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	if h.opts.addCaller {
		stack := internal.Capture(0, internal.Full)
		defer stack.Free()
		for frame, ok := stack.Next(); ok; frame, ok = stack.Next() {
			pkg := pkgFilePath(&frame)
			if h.opts.Match(pkg) {
				if h.opts.callerFormater != nil {
					record.AddAttrs(slog.Any("caller", h.opts.callerFormater(frame.File, frame.Function, frame.Line)))
				} else {
					record.AddAttrs(slog.Any("caller", slog.Source{
						Function: frame.Function,
						File:     frame.File,
						Line:     frame.Line,
					}))
				}
				break
			}
		}
	}
	return h.h.Handle(ctx, record)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), opts: h.opts}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), opts: h.opts}
}
