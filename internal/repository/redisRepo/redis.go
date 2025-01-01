package redisRepo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/redis/go-redis/v9"
)

var (
	SessionPrefix = "session"
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

// func Check() {

// 	ctx := context.Background()

// 	s := NewSessionRepo()

// id, err := s.CreateSession(ctx, &entity.Session{UserID: 3, UserAgent: "test", UserIP: "test1"})
// fmt.Println(id, err)

// res, err := s.GetSession(ctx, id, 3)
// fmt.Println(res, err)

// err = s.DeleteSession(ctx, id, 3)
// fmt.Println(err)

// id, err = s.CreateSession(ctx, &entity.Session{UserID: 1, UserAgent: "test", UserIP: "test1"})
// fmt.Print(id, err)
// id, err = s.CreateSession(ctx, &entity.Session{UserID: 1, UserAgent: "test2", UserIP: "test2"})
// fmt.Print(id, err)

// 	list, err := s.GetListSession(ctx, 1)
// 	fmt.Println(list, err)

// 	err = s.DeleteListSession(ctx, 3)
// 	fmt.Println(err)

//		list, err = s.GetListSession(ctx, 1)
//		fmt.Println(list, err)
//	}
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
	key := fmt.Sprintf("%s:%d:%s", SessionPrefix, s.UserID, s.ID)

	err = r.rdb.HSet(ctx, key, s).Err()
	if err != nil {
		return nil, fmt.Errorf("redis create session: %w", err)
	}
	return s, nil
}

func (r *SessionRepo) GetSession(ctx context.Context, id string, user_id int) (*entity.Session, error) {
	res := new(entity.Session)
	key := fmt.Sprintf("%s:%d:%s", SessionPrefix, user_id, id)

	if err := r.rdb.HGetAll(ctx, key).Scan(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, id string, user_id int) error {
	key := fmt.Sprintf("%s:%d:%s", SessionPrefix, user_id, id)
	return r.rdb.Del(ctx, key).Err()
}

func (r *SessionRepo) GetListSession(ctx context.Context, user_id int) ([]*entity.Session, error) {
	sessions := []*entity.Session{}
	pattern := fmt.Sprintf("%s:%d:*", SessionPrefix, user_id)
	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		session := new(entity.Session)
		if err := r.rdb.HGetAll(ctx, key).Scan(session); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *SessionRepo) DeleteListSession(ctx context.Context, user_id int) error {
	pattern := fmt.Sprintf("%s:%d:*", SessionPrefix, user_id)
	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	return r.rdb.Del(ctx, keys...).Err()
}
