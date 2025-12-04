package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ridKey struct{}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-ID")
			if rid == "" {
				rid = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", rid)
			ctx := context.WithValue(r.Context(), ridKey{}, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RIDFrom returns the request ID from context (or empty string if absent).
func RIDFrom(ctx context.Context) string {
	if v := ctx.Value(ridKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// --- Passthrough optional interfaces ---

// Ensure SSE works: expose Flush if the underlying writer can flush.
func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// (Optional) WebSockets/HTTP2 support
// Only compile these if you have the imports; otherwise you can omit them.

// import ( "bufio"; "fmt"; "net" )

// func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// 	h, ok := w.ResponseWriter.(http.Hijacker)
// 	if !ok {
// 		return nil, nil, fmt.Errorf("hijacker not supported")
// 	}
// 	return h.Hijack()
// }

// func (w *statusWriter) Push(target string, opts *http.PushOptions) error {
// 	if p, ok := w.ResponseWriter.(http.Pusher); ok {
// 		return p.Push(target, opts)
// 	}
// 	return http.ErrNotSupported
// }


func Recover(log zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error().Interface("panic", rec).Msg("panic recovered")
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func AccessLog(log zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w, status: 200}
			start := time.Now()
			next.ServeHTTP(sw, r)
			log.Info().
				Str("route", r.URL.Path).
				Str("method", r.Method).
				Int("status", sw.status).
				Int64("latency_ms", time.Since(start).Milliseconds()).
				Msg("req")
		})
	}
}
