// Package stdslog slog config for get separated output
package stdslog

import (
	"context"
	"log/slog"
)

type SplitHandler struct {
	Out, Err slog.Handler
}

func (h *SplitHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Out.Enabled(ctx, level) || h.Err.Enabled(ctx, level)
}

func (h *SplitHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.Err.Handle(ctx, r)
	}
	return h.Out.Handle(ctx, r)
}

func (h *SplitHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SplitHandler{
		Out: h.Out.WithAttrs(attrs),
		Err: h.Err.WithAttrs(attrs),
	}
}

func (h *SplitHandler) WithGroup(name string) slog.Handler {
	return &SplitHandler{
		Out: h.Out.WithGroup(name),
		Err: h.Err.WithGroup(name),
	}
}
