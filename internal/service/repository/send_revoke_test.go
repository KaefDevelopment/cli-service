package repository_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service"
	"github.com/KaefDevelopment/cli-service/internal/service/repository"
)

// These tests exercise the full Send() path against a real sqlite queue and a
// real gorm transaction so we can prove the REVOKED wipe commits OUTSIDE the
// rolled-back send transaction (the core correctness requirement of §3.3).

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "cli.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, model.InitSchema(context.Background(), db))
	return db
}

func seedEvent(t *testing.T, db *gorm.DB, authKey string) {
	t.Helper()

	repo := repository.NewCLIRepository(db)
	err := repo.Create(context.Background(), model.Events{Events: []model.Event{{
		Id:        "event-1",
		CreatedAt: "1",
		Type:      "1",
		Project:   "proj",
		Timezone:  "1",
		AuthKey:   authKey,
		Send:      false,
	}}})
	require.NoError(t, err)
}

func countEvents(t *testing.T, db *gorm.DB) int64 {
	t.Helper()

	var c int64
	require.NoError(t, db.Table("events").Count(&c).Error)
	return c
}

func markedCount(t *testing.T, db *gorm.DB, send int) int64 {
	t.Helper()

	var c int64
	require.NoError(t, db.Table("events").Where("send = ?", send).Count(&c).Error)
	return c
}

func newService(db *gorm.DB, addr string) *service.CLIService {
	repo := repository.NewCLIRepository(db)
	txp := repository.NewTxProvider(db)
	return service.NewCLIService(repo, txp, addr, "test-auth-key")
}

// 410 + {"code":"PLUGIN_REVOKED"} -> queue wiped, ErrRevoked returned (not retried),
// and the wipe survives even though the send transaction rolled back.
func TestSend_Revoked_WipesQueueOutsideRolledBackTx(t *testing.T) {
	db := newTestDB(t)
	seedEvent(t, db, "test-auth-key")
	require.Equal(t, int64(1), countEvents(t, db))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte(`{"code":"PLUGIN_REVOKED","description":"ACCOUNT_DELETED"}`))
	}))
	defer srv.Close()

	err := newService(db, srv.URL).Send(context.Background(), "1.0.0")

	// The error propagates so the caller can emit the stdout revoked marker,
	// but it is ErrRevoked (NOT the generic bad-status error) so the queue is
	// not retried on the next run.
	require.ErrorIs(t, err, service.ErrRevoked)

	// The send tx (which had marked the row send=1) rolled back; the wipe ran on
	// the base repo in a separate tx. If the wipe had been inside the rolled-back
	// tx, the row would still be here. It is gone -> wipe committed independently.
	assert.Equal(t, int64(0), countEvents(t, db))
}

// HTTP 500 -> generic error, queue kept intact and the mark is rolled back
// (send reset to 0) so it retries; NOT treated as revoked.
func TestSend_ServerError_KeepsQueue(t *testing.T) {
	db := newTestDB(t)
	seedEvent(t, db, "test-auth-key")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	err := newService(db, srv.URL).Send(context.Background(), "1.0.0")

	require.Error(t, err)
	assert.NotErrorIs(t, err, service.ErrRevoked)

	// Queue intact and rolled back to send=0 (ready to retry).
	assert.Equal(t, int64(1), countEvents(t, db))
	assert.Equal(t, int64(1), markedCount(t, db, 0))
}

// HTTP 401 (unknown uid) -> queue kept, not revoked. A revoke must be a positive
// 410+PLUGIN_REVOKED signal, never a bare 401.
func TestSend_Unauthorized_KeepsQueue(t *testing.T) {
	db := newTestDB(t)
	seedEvent(t, db, "test-auth-key")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":"UNAUTHORIZED"}`))
	}))
	defer srv.Close()

	err := newService(db, srv.URL).Send(context.Background(), "1.0.0")

	require.Error(t, err)
	assert.NotErrorIs(t, err, service.ErrRevoked)
	assert.Equal(t, int64(1), countEvents(t, db))
}

// 410 WITHOUT the PLUGIN_REVOKED body marker must NOT trigger a wipe (both
// checks required): a spurious 410 from a proxy stays a keep-and-retry error.
func TestSend_GoneWithoutMarker_KeepsQueue(t *testing.T) {
	db := newTestDB(t)
	seedEvent(t, db, "test-auth-key")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte(`<html>410 Gone</html>`))
	}))
	defer srv.Close()

	err := newService(db, srv.URL).Send(context.Background(), "1.0.0")

	require.Error(t, err)
	assert.NotErrorIs(t, err, service.ErrRevoked)
	assert.Equal(t, int64(1), countEvents(t, db))
}
