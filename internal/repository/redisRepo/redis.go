package redisRepo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/pkg/apiError"
	"github.com/redis/go-redis/v9"
)

var (
	SessionPrefix      = "session"
	UserSessionsPrefix = "user_sessions"
)

type SessionRepo struct {
	rdb *redis.Client
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		rdb: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func GenerateRandomSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *SessionRepo) CreateSession(ctx context.Context, s *entity.Session) (*entity.Session, error) {
	id, err := GenerateRandomSessionID()
	if err != nil {
		return nil, err
	}

	s.ID = id

	// create "session <session_id>{"user_id": int, "user_ip": "string", "user_agent": "string"}""
	key_session := fmt.Sprintf("%s:%s", SessionPrefix, s.ID)
	err = r.rdb.HSet(ctx, key_session, s).Err()
	if err != nil {
		return nil, fmt.Errorf("redis create session: %w", err)
	}

	// Установка времени жизни для сессии
	err = r.rdb.Expire(ctx, key_session, 2*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("redis set session TTL: %w", err)
	}

	// create "user_sessions <user_id>["session_id",...]""
	key_user := fmt.Sprintf("%s:%d", UserSessionsPrefix, s.UserID)
	err = r.rdb.SAdd(ctx, key_user, s.ID).Err()
	if err != nil {
		return nil, fmt.Errorf("redis add session in user: %w", err)
	}
	return s, nil
}

func (r *SessionRepo) GetSession(ctx context.Context, id string) (*entity.Session, error) {
	res := new(entity.Session)
	key := fmt.Sprintf("%s:%s", SessionPrefix, id)

	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get session TTL: %w", err)
	}

	if ttl <= 0 {
		return nil, fmt.Errorf("redis get session: %w", apiError.ErrSessionExpired)

	}

	if err := r.rdb.HGetAll(ctx, key).Scan(res); err != nil {
		return nil, fmt.Errorf("redis get session: %w", err)
	}

	return res, nil
}

// delete "session <session_id>{"user_id": int, "user_ip": "string", "user_agent": "string"}""
func (r *SessionRepo) DeleteSession(ctx context.Context, id string) error {
	key := fmt.Sprintf("%s:%s", SessionPrefix, id)
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete session: %w", err)
	}
	return nil
}

// delete from "user_sessions <user_id>["session_id",...]""
func (r *SessionRepo) DeleteSessionFromUser(ctx context.Context, id string, user_id int) error {
	key := fmt.Sprintf("%s:%d", UserSessionsPrefix, user_id)
	if err := r.rdb.SRem(ctx, key, id).Err(); err != nil {
		return fmt.Errorf("redis delete session from user: %w", err)
	}
	return nil
}
func (r *SessionRepo) GetListSession(ctx context.Context, user_id int) ([]*entity.Session, error) {
	sessions := []*entity.Session{}
	key_user := fmt.Sprintf("%s:%d", UserSessionsPrefix, user_id)
	sessions_id, err := r.rdb.SMembers(ctx, key_user).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get list sessions: %w", err)
	}
	for _, session_id := range sessions_id {
		session, err := r.GetSession(ctx, session_id)
		if errors.Is(err, apiError.ErrSessionExpired) {
			r.DeleteSessionFromUser(ctx, session_id, user_id)

		}
		if err != nil && !errors.Is(err, apiError.ErrSessionExpired) {
			return nil, fmt.Errorf("redis get list sessions: %w", err)
		}
		if session != nil {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (r *SessionRepo) DeleteListSession(ctx context.Context, user_id int) error {
	key_user := fmt.Sprintf("%s:%d", UserSessionsPrefix, user_id)
	sessions_id, err := r.rdb.SMembers(ctx, key_user).Result()
	if err != nil {
		return fmt.Errorf("redis delete list sessions: %w", err)
	}

	for _, session_id := range sessions_id {
		if err := r.DeleteSession(ctx, session_id); err != nil {
			return fmt.Errorf("redis delete list session: %w", err)
		}
		if err := r.DeleteSessionFromUser(ctx, session_id, user_id); err != nil {
			return fmt.Errorf("redis delete list session: %w", err)
		}
	}
	return nil
}
