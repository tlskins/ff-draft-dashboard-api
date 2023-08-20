package store

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tryhungry/hungry-games-api/api"
	t "github.com/tryhungry/hungry-games-api/types"
)

var colUser = "users"

func (s *Store) GetUser(id string) (out *t.User, err error) {
	sess, c := s.C(colUser)
	defer sess.Close()

	out = &t.User{}
	err = api.FindOne(c, out, api.M{"_id": id})
	return
}

func (s *Store) GetUserByLowerEmail(lowerEmail string) (out *t.User, err error) {
	sess, c := s.C(colUser)
	defer sess.Close()

	out = &t.User{}
	err = api.FindOne(c, out, api.M{"lowEmail": lowerEmail})
	return
}

func (s *Store) CreateUser(user *t.User) (out *t.User, err error) {
	sess, c := s.C(colUser)
	defer sess.Close()

	if user.Id == "" {
		user.Id = uuid.NewV4().String()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	out = &t.User{}
	err = api.Upsert(c, &out, M{"_id": user.Id}, M{"$set": user})
	return
}
