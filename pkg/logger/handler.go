package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Color constants
var (
	colorReset  = []byte("\033[0m")
	colorGray   = []byte("\033[90m")
	colorCyan   = []byte("\033[36m")
	colorYellow = []byte("\033[93m")
	colorRed    = []byte("\033[91m")
	colorWhite  = []byte("\033[97m")
)

// Pre-allocated level strings to avoid repeated concatenation
var levelStrings = map[slog.Level][]byte{
	slog.LevelDebug: append(append(colorGray, []byte("DEBUG:")...), colorReset...),
	slog.LevelInfo:  append(append(colorCyan, []byte("INFO:")...), colorReset...),
	slog.LevelWarn:  append(append(colorYellow, []byte("WARN:")...), colorReset...),
	slog.LevelError: append(append(colorRed, []byte("ERROR:")...), colorReset...),
}

const (
	FormatJSON = "json"
	FormatText = "text"
)

type Handler struct {
	opts       *slog.HandlerOptions
	out        io.Writer
	format     string
	bufferPool sync.Pool
	attrsPool  sync.Pool
	baseAttrs  []slog.Attr
	encoder    *json.Encoder
	mutex      sync.Mutex
}

func NewHandler(opts *slog.HandlerOptions, format string) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	h := &Handler{
		opts:   opts,
		out:    os.Stdout,
		format: format,
		bufferPool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 0, 1024))
			},
		},
		attrsPool: sync.Pool{
			New: func() any {
				return make(map[string]any, 10)
			},
		},
	}

	h.encoder = json.NewEncoder(h.out)
	return h
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newOpts := *h.opts
	newHandler := NewHandler(&newOpts, h.format)
	newHandler.baseAttrs = append(h.baseAttrs, attrs...)
	return newHandler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h // Groups are handled in Handle method
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	switch h.format {
	case FormatJSON:
		return h.handleJSON(r)
	default:
		return h.handleText(r)
	}
}

func (h *Handler) handleJSON(r slog.Record) error {
	attrs := h.attrsPool.Get().(map[string]any)
	clear(attrs)
	defer h.attrsPool.Put(attrs)

	// Apply baseAttrs first
	for _, a := range h.baseAttrs {
		val := a.Value.Any()
		if errVal, ok := val.(error); ok {
			attrs[a.Key] = errVal.Error()
		} else {
			attrs[a.Key] = val
		}
	}

	// Add standard fields
	attrs["time"] = r.Time.Format(time.RFC3339)
	attrs["level"] = r.Level.String()
	attrs["msg"] = r.Message

	// Add record attrs
	r.Attrs(func(a slog.Attr) bool {
		val := a.Value.Any()
		if errVal, ok := val.(error); ok {
			attrs[a.Key] = errVal.Error()
		} else {
			attrs[a.Key] = val
		}
		return true
	})

	return h.encoder.Encode(attrs)
}

func (h *Handler) handleText(r slog.Record) error {
	buf := h.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer h.bufferPool.Put(buf)

	// Write time
	timeStr := r.Time.Format(time.RFC822)
	if _, err := buf.Write(colorGray); err != nil {
		return fmt.Errorf("failed to write color: %w", err)
	}
	if _, err := buf.WriteString(timeStr); err != nil {
		return fmt.Errorf("failed to write time: %w", err)
	}
	if _, err := buf.Write(colorReset); err != nil {
		return fmt.Errorf("failed to write color reset: %w", err)
	}
	buf.WriteByte(' ')

	// Write level
	levelBytes, ok := levelStrings[r.Level]
	if !ok {
		levelBytes = []byte(r.Level.String())
	}
	if _, err := buf.Write(levelBytes); err != nil {
		return fmt.Errorf("failed to write level: %w", err)
	}
	buf.WriteByte(' ')

	// Write message
	if _, err := buf.Write(colorWhite); err != nil {
		return fmt.Errorf("failed to write color: %w", err)
	}
	if _, err := buf.WriteString(r.Message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if _, err := buf.Write(colorReset); err != nil {
		return fmt.Errorf("failed to write color reset: %w", err)
	}
	buf.WriteByte(' ')

	// Write attributes
	attrs := h.attrsPool.Get().(map[string]any)
	clear(attrs)
	defer h.attrsPool.Put(attrs)

	for _, a := range h.baseAttrs {
		val := a.Value.Any()
		if errVal, ok := val.(error); ok {
			attrs[a.Key] = errVal.Error()
		} else {
			attrs[a.Key] = val
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return true
		}
		val := a.Value.Any()
		if errVal, ok := val.(error); ok {
			attrs[a.Key] = errVal.Error()
		} else {
			attrs[a.Key] = val
		}
		return true
	})

	if len(attrs) > 0 {
		if _, err := buf.Write(colorGray); err != nil {
			return fmt.Errorf("failed to write color: %w", err)
		}
		jsonBytes, err := json.Marshal(attrs)
		if err != nil {
			return fmt.Errorf("failed to marshal attributes: %w", err)
		}
		if _, err := buf.Write(jsonBytes); err != nil {
			return fmt.Errorf("failed to write attributes: %w", err)
		}
		if _, err := buf.Write(colorReset); err != nil {
			return fmt.Errorf("failed to write color reset: %w", err)
		}
	}

	buf.WriteByte('\n')

	_, err := h.out.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}

	return nil
}

// Helper function to clear a map
func clear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}
