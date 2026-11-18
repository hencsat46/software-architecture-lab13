package service

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"lab13/order-service/internal/common/exceptions"
	"lab13/order-service/internal/logger"
	domain "lab13/order-service/internal/orders"
)

type repoStub struct {
	IPostgres
	created *domain.Entity
	err     error
}

func (r *repoStub) Create(_ context.Context, e *domain.Entity) (*domain.Entity, error) {
	r.created = e
	return e, r.err
}

type userStub struct{ err error }

func (u userStub) Exists(context.Context, string) error { return u.err }

type notificationStub struct {
	err   error
	calls int
}

func (n *notificationStub) Send(context.Context, string, string) error { n.calls++; return n.err }

func TestCreateDependencyFailures(t *testing.T) {
	for _, tt := range []struct {
		name                      string
		userErr, dbErr, notifyErr error
		wantSaved, wantSent       bool
	}{
		{"success", nil, nil, nil, true, true},
		{"missing user", exceptions.ErrInvalidUser, nil, nil, false, false},
		{"user unavailable", exceptions.ErrUnavailable, nil, nil, false, false},
		{"database failure", nil, errors.New("db"), nil, false, false},
		{"notification unavailable", nil, nil, errors.New("offline"), true, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repoStub{err: tt.dbErr}
			n := &notificationStub{err: tt.notifyErr}
			s := NewService(repo, userStub{tt.userErr}, n, &logger.Logger{Logger: zap.NewNop()})
			e, sent, err := s.Create(context.Background(), &domain.Entity{UserID: "user", Product: "book", Quantity: 1})
			if (err == nil) != tt.wantSaved || sent != tt.wantSent {
				t.Fatalf("entity=%v sent=%v err=%v", e, sent, err)
			}
			if tt.wantSaved && (e.ID == "" || e.CreatedAt.IsZero()) {
				t.Fatal("metadata missing")
			}
			if tt.userErr != nil && repo.created != nil {
				t.Fatal("saved without validated user")
			}
			if !tt.wantSaved && n.calls != 0 {
				t.Fatal("notified before save")
			}
		})
	}
}
