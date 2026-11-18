package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"lab13/order-service/internal/common/exceptions"
)

func TestExists(t *testing.T) {
	for _, tt := range []struct {
		name string
		code int
		body string
		want error
	}{
		{"found", 200, `{"data":{"id":"user"}}`, nil},
		{"missing", 404, `{}`, exceptions.ErrInvalidUser},
		{"server failure", 500, `{}`, exceptions.ErrUnavailable},
		{"invalid JSON", 200, `oops`, exceptions.ErrUnavailable},
		{"wrong user", 200, `{"data":{"id":"other"}}`, exceptions.ErrUnavailable},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/users/user" {
					t.Errorf("path %s", r.URL.Path)
				}
				w.WriteHeader(tt.code)
				w.Write([]byte(tt.body))
			}))
			defer s.Close()
			if got := New(s.Client(), s.URL, s.URL).Exists(context.Background(), "user"); got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkFailureAndCancellation(t *testing.T) {
	s := httptest.NewServer(http.NotFoundHandler())
	s.Close()
	c := New(s.Client(), s.URL, s.URL)
	if c.Exists(context.Background(), "user") != exceptions.ErrUnavailable {
		t.Fatal("network failure ignored")
	}
	if c.Send(context.Background(), "user", "hello") == nil {
		t.Fatal("notification failure ignored")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if c.Exists(ctx, "user") != exceptions.ErrUnavailable {
		t.Fatal("cancel ignored")
	}
}

func TestNotificationStatus(t *testing.T) {
	for _, status := range []int{201, 400, 500} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
				t.Error("invalid request")
			}
			w.WriteHeader(status)
		}))
		err := New(s.Client(), s.URL, s.URL).Send(context.Background(), "user", "message")
		s.Close()
		if (err == nil) != (status == 201) {
			t.Fatalf("status=%d err=%v", status, err)
		}
	}
}
