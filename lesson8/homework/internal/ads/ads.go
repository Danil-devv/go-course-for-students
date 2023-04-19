package ads

import "time"

type Repository interface {
	AddAd(ad Ad) int64
	GetById(id int64) (Ad, error)
	ReplaceByID(id int64, ad Ad) error
	GetSize() int64
}

type Ad struct {
	ID         int64
	Title      string `validate:"min:1;max:99"`
	Text       string `validate:"min:1;max:499"`
	AuthorID   int64
	Published  bool
	CreateDate time.Time
	LastUpdate time.Time
}
