package users

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --output=./tests/mocks --name=Repository
type Repository interface {
	AddUser(u User) error
	GetById(id int64) (User, error)
	ReplaceByID(id int64, u User) error
	DeleteByID(id int64) (User, error)
}

type User struct {
	ID       int64
	Nickname string
	Email    string
}
