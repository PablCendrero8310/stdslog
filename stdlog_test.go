package stdslog

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestSplitHandlerRoutesLevels(t *testing.T) {
	var stdout, stderr bytes.Buffer

	handler := &SplitHandler{
		Out: slog.NewTextHandler(&stdout, nil),
		Err: slog.NewTextHandler(&stderr, nil),
	}

	logger := slog.New(handler)

	logger.Info("mensaje informativo")
	logger.Warn("mensaje de advertencia")
	logger.Error("mensaje de error")

	if !bytes.Contains(stdout.Bytes(), []byte("mensaje informativo")) {
		t.Error("INFO debería escribirse en stdout")
	}

	if bytes.Contains(stdout.Bytes(), []byte("mensaje de advertencia")) ||
		bytes.Contains(stdout.Bytes(), []byte("mensaje de error")) {
		t.Error("WARN y ERROR no deberían escribirse en stdout")
	}

	if !bytes.Contains(stderr.Bytes(), []byte("mensaje de advertencia")) {
		t.Error("WARN debería escribirse en stderr")
	}

	if !bytes.Contains(stderr.Bytes(), []byte("mensaje de error")) {
		t.Error("ERROR debería escribirse en stderr")
	}
}

func TestSplitHandlerEnabled(t *testing.T) {
	ctx := context.Background()

	handler := &SplitHandler{
		Out: slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
		Err: slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	}

	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("DEBUG no debería estar habilitado")
	}

	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("INFO debería estar habilitado")
	}

	if !handler.Enabled(ctx, slog.LevelWarn) {
		t.Error("WARN debería estar habilitado")
	}
}

func TestSplitHandlerWithAttrsAndGroup(t *testing.T) {
	var stdout, stderr bytes.Buffer

	handler := &SplitHandler{
		Out: slog.NewTextHandler(&stdout, nil),
		Err: slog.NewTextHandler(&stderr, nil),
	}

	logger := slog.New(handler).
		With("module", "mapping").
		WithGroup("test")

	logger.Info("mensaje", "key", "value")

	output := stdout.String()

	if !bytes.Contains([]byte(output), []byte("module=mapping")) {
		t.Errorf("WithAttrs debería conservar los atributos; salida: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("test.key=value")) {
		t.Errorf("WithGroup debería agrupar los atributos; salida: %s", output)
	}
}
