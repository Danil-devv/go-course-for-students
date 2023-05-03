package usersrepo

import (
	"errors"
	"homework9/internal/users"
	"sync"
)

var wrongIdErr = errors.New("id must be non-negative and must be less than repository size")

func New() users.Repository {
	return &userRepo{repo: make(map[int64]users.User)}
}

type userRepo struct {
	repo map[int64]users.User
	m    sync.RWMutex
}

func (r *userRepo) checkID(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.repo[id]; ok {
		return nil
	}

	return wrongIdErr
}

func (r *userRepo) AddUser(u users.User) error {
	// если пользователь с данным ID уже существует
	if r.checkID(u.ID) == nil {
		return wrongIdErr
	}

	r.m.Lock()
	r.repo[u.ID] = u
	r.m.Unlock()

	return nil
}

func (r *userRepo) GetById(id int64) (users.User, error) {
	if err := r.checkID(id); err != nil {
		return users.User{}, err
	}

	r.m.Lock()
	user := r.repo[id]
	r.m.Unlock()

	return user, nil
}

func (r *userRepo) ReplaceByID(id int64, u users.User) error {
	if err := r.checkID(id); err != nil {
		return err
	}

	r.m.Lock()
	r.repo[id] = u
	r.m.Unlock()

	return nil
}
