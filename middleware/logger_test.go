package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

// capturingLogger captures log output for assertions.
type capturingLogger struct {
	mu       sync.Mutex
	messages []string
}

func (c *capturingLogger) Print(v ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	parts := make([]string, len(v))
	for i, val := range v {
		parts[i] = val.(string)
	}
	c.messages = append(c.messages, strings.Join(parts, " "))
}

func (c *capturingLogger) messageCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.messages)
}

func newTestLogger(level middleware.Level) (*capturingLogger, func(http.Handler) http.Handler) {
	logger := &capturingLogger{}
	formatter := &middleware.DefaultLogFormatter{
		Logger:   logger,
		NoColor:  true,
		LogLevel: level,
	}
	return logger, middleware.RequestLogger(formatter)
}

func handlerWithStatus(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}

func TestStatusLevel_Mapping(t *testing.T) {
	tests := []struct {
		name      string
		level     middleware.Level
		status    int
		shouldLog bool
	}{
		// Log level (0) — logs everything
		{"Log level sees 200", middleware.Log, 200, true},
		{"Log level sees 404", middleware.Log, 404, true},
		{"Log level sees 500", middleware.Log, 500, true},

		// Info level (1) — logs everything (all statuses map to Info+)
		{"Info level sees 200", middleware.Info, 200, true},
		{"Info level sees 301", middleware.Info, 301, true},
		{"Info level sees 404", middleware.Info, 404, true},
		{"Info level sees 500", middleware.Info, 500, true},

		// Error level (2) — only 4xx and 5xx
		{"Error level skips 200", middleware.Error, 200, false},
		{"Error level skips 201", middleware.Error, 201, false},
		{"Error level skips 301", middleware.Error, 301, false},
		{"Error level sees 400", middleware.Error, 400, true},
		{"Error level sees 404", middleware.Error, 404, true},
		{"Error level sees 500", middleware.Error, 500, true},
		{"Error level sees 503", middleware.Error, 503, true},

		// Fatal level (3) — only 5xx
		{"Fatal level skips 200", middleware.Fatal, 200, false},
		{"Fatal level skips 404", middleware.Fatal, 404, false},
		{"Fatal level skips 499", middleware.Fatal, 499, false},
		{"Fatal level sees 500", middleware.Fatal, 500, true},
		{"Fatal level sees 502", middleware.Fatal, 502, true},
		{"Fatal level sees 503", middleware.Fatal, 503, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, logMiddleware := newTestLogger(tt.level)

			handler := logMiddleware(handlerWithStatus(tt.status))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if tt.shouldLog {
				assert.Equal(t, 1, logger.messageCount(), "expected log output")
			} else {
				assert.Equal(t, 0, logger.messageCount(), "expected no log output")
			}
		})
	}
}

func TestLogger_200Success_LoggedAtInfoLevel(t *testing.T) {
	logger, logMiddleware := newTestLogger(middleware.Info)

	handler := logMiddleware(handlerWithStatus(http.StatusOK))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 1, logger.messageCount())
	assert.Contains(t, logger.messages[0], "200")
}

func TestLogger_201Created_LoggedAtInfoLevel(t *testing.T) {
	logger, logMiddleware := newTestLogger(middleware.Info)

	handler := logMiddleware(handlerWithStatus(http.StatusCreated))
	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 1, logger.messageCount())
	assert.Contains(t, logger.messages[0], "201")
}

func TestLogger_500Error_LoggedAtAllLevels(t *testing.T) {
	levels := []middleware.Level{middleware.Log, middleware.Info, middleware.Error, middleware.Fatal}

	for _, level := range levels {
		t.Run("level_"+levelName(level), func(t *testing.T) {
			logger, logMiddleware := newTestLogger(level)

			handler := logMiddleware(handlerWithStatus(http.StatusInternalServerError))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, 1, logger.messageCount(), "500 should be logged at level %d", level)
		})
	}
}

func TestLogger_200_NotLoggedAtErrorLevel(t *testing.T) {
	logger, logMiddleware := newTestLogger(middleware.Error)

	handler := logMiddleware(handlerWithStatus(http.StatusOK))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 0, logger.messageCount(), "200 should not be logged at Error level")
}

func TestLogger_404_NotLoggedAtFatalLevel(t *testing.T) {
	logger, logMiddleware := newTestLogger(middleware.Fatal)

	handler := logMiddleware(handlerWithStatus(http.StatusNotFound))
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 0, logger.messageCount(), "404 should not be logged at Fatal level")
}

func TestLogger_IncludesMethod(t *testing.T) {
	logger, logMiddleware := newTestLogger(middleware.Log)

	handler := logMiddleware(handlerWithStatus(http.StatusOK))
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 1, logger.messageCount())
	assert.Contains(t, logger.messages[0], "POST")
}

func TestSetupLogger(t *testing.T) {
	ls := middleware.SetupLogger(middleware.Error)
	assert.NotNil(t, ls)
}

func TestGetLogEntry_NilContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	entry := middleware.GetLogEntry(req)
	assert.Nil(t, entry)
}

func TestWithLogEntry_RoundTrip(t *testing.T) {
	logger, _ := newTestLogger(middleware.Log)
	formatter := &middleware.DefaultLogFormatter{
		Logger:  logger,
		NoColor: true,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	entry := formatter.NewLogEntry(req)

	req = middleware.WithLogEntry(req, entry)
	retrieved := middleware.GetLogEntry(req)

	assert.NotNil(t, retrieved)
}

func levelName(l middleware.Level) string {
	switch l {
	case middleware.Log:
		return "Log"
	case middleware.Info:
		return "Info"
	case middleware.Error:
		return "Error"
	case middleware.Fatal:
		return "Fatal"
	default:
		return "Unknown"
	}
}
