package usersrepo

import (
	"errors"
	"homework8/internal/users"
)

var wrongIdErr = errors.New("id must be non-negative and must be less than repository size")

func New() users.Repository {
	return &userRepo{repo: make(map[int64]users.User)}
}

type userRepo struct {
	repo map[int64]users.User
}

func (r *userRepo) checkID(id int64) error {
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

	r.repo[u.ID] = u
	return nil
}

func (r *userRepo) GetById(id int64) (users.User, error) {
	if err := r.checkID(id); err != nil {
		return users.User{}, err
	}

	return r.repo[id], nil
}

func (r *userRepo) ReplaceByID(id int64, u users.User) error {
	if err := r.checkID(id); err != nil {
		return err
	}

	r.repo[id] = u
	return nil
}
