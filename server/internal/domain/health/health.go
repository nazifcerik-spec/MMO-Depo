// Package health serves /live and /ready probes. /live is process-up only; /ready
// checks DB + Redis (see infra/reverse-proxy.md — LB promotes to healthy only after
// /ready 200).
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Service holds dependencies the /ready probe checks.
type Service struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewService(db *pgxpool.Pool, rdb *redis.Client) *Service {
	return &Service{db: db, rdb: rdb}
}

// Live returns 200 as long as the process is up. No dependency checks.
func Live(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready checks Postgres + Redis; returns 503 if either is degraded.
func (s *Service) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resp := map[string]string{"status": "ok"}
	status := http.StatusOK

	if err := s.db.Ping(ctx); err != nil {
		resp["db"] = "fail: " + err.Error()
		status = http.StatusServiceUnavailable
	} else {
		resp["db"] = "ok"
	}

	if err := s.rdb.Ping(ctx).Err(); err != nil {
		resp["redis"] = "fail: " + err.Error()
		status = http.StatusServiceUnavailable
	} else {
		resp["redis"] = "ok"
	}

	if status != http.StatusOK {
		resp["status"] = "degraded"
	}
	JSON(w, status, resp)
}

// JSON is the shared response helper — every handler serializes through it so
// content-type and error handling stay consistent.
func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
