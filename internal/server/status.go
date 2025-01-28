package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/alfatraining/cloud-infrastructure-engineer/internal/otel"
)

type statusHandler struct {
	flaky      bool // randomizes response time and status code
	statusCode int  // represents http status codes: https://pkg.go.dev/net/http#pkg-constants
	observer   otel.Observer
}

func (h *statusHandler) write(ctx context.Context, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	if _, err := fmt.Fprintln(w, http.StatusText(code)); err != nil {
		h.observer.Error(ctx, fmt.Sprintf("write failed: %s\n", err))
		return
	}

	switch {
	case code >= 400 && code <= 499:
		h.observer.Info(ctx, fmt.Sprintf("request failed: %s", http.StatusText(code)))
		return
	case code >= 500 && code <= 599:
		h.observer.Error(ctx, fmt.Sprintf("server error: %s", http.StatusText(code)))
		return
	}
	h.observer.Debug(ctx, "request successful")
}

// List of valid go http status codes.
var statusCodes = []int{
	100, 101, 102, 103, // 10x
	200, 201, 202, 203, 204, 205, 206, 207, 208, 226, // 20x
	300, 301, 302, 303, 304, 305, 307, 308, // 30x
	400, 401, 402, 403, 404, 405, 406, 407, 408, 409, // 40x
	410, 411, 412, 413, 414, 415, 416, 417, 418, // 41x
	421, 422, 423, 424, 426, 428, 429, // 42x
	431, 451,
	500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511, // 5xx
}

// Returns a random valid http status code.
func randomStatus() int {
	return statusCodes[rand.Intn(len(statusCodes))]
}

// Blocks between 0ms and 500ms.
func randomSleep() {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

// Get returns the current status code of the server.
// Randomizes its output if flakiness is set to 'true'.
func (h *statusHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.write(r.Context(), w, http.StatusMethodNotAllowed)
		return
	}

	statusCode := h.statusCode
	if h.flaky {
		randomSleep()
		h.statusCode = randomStatus()
	}
	h.write(r.Context(), w, statusCode)
}

// Post sets a given status code for the '/status' endpoint.
func (h *statusHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.write(r.Context(), w, http.StatusMethodNotAllowed)
		return
	}

	code, err := strconv.Atoi(r.PathValue("code"))
	if err != nil {
		h.write(r.Context(), w, http.StatusBadRequest)
		return
	}

	if !slices.Contains(statusCodes, code) {
		h.write(r.Context(), w, http.StatusBadRequest)
		return
	}

	hasChange := h.statusCode != code
	h.statusCode = code
	if hasChange {
		h.write(r.Context(), w, http.StatusAccepted)
		return
	}
	h.write(r.Context(), w, http.StatusOK)
}

// Flaky enables or disables randomness of the '/status' endpoint.
// When enabled it takes precedence over any given '/status/{code}'.
func (h *statusHandler) Flaky(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.write(r.Context(), w, http.StatusMethodNotAllowed)
		return
	}

	if h.flaky {
		h.flaky = false
		h.write(r.Context(), w, http.StatusAccepted)
		return
	}
	h.flaky = true
	h.write(r.Context(), w, http.StatusAccepted)
}
