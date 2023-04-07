package app

import (
	"errors"
	validator "github.com/Danil-devv/structValidator"
	"homework6/internal/ads"
)

var (
	ValidationErr = errors.New("some fields does not pass the validation")
	AccessErr     = errors.New("user can only change his ads")
)

type App interface {
	CreateAd(title string, text string, authorID int64) (ads.Ad, error)
	ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error)
}

type Repository interface {
	AddAd(ad ads.Ad) int64
	GetById(id int64) (ads.Ad, error)
	ReplaceByID(id int64, ad ads.Ad) error
	GetSize() int64
}

func NewApp(repo Repository) App {
	return &app{repo: repo}
}

type app struct {
	repo Repository
}

func (a *app) CreateAd(title string, text string, authorID int64) (ads.Ad, error) {
	ad := ads.Ad{ID: a.repo.GetSize(), Title: title, Text: text, AuthorID: authorID}

	if err := validator.Validate(ad); err != nil {
		return ads.Ad{}, ValidationErr
	}

	a.repo.AddAd(ad)
	return ad, nil
}

func (a *app) ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error) {
	ad, err := a.repo.GetById(adID)
	if err != nil {
		return ads.Ad{}, err
	}

	if ad.AuthorID != userID {
		return ads.Ad{}, AccessErr
	}

	ad.Published = published
	return ad, a.repo.ReplaceByID(adID, ad)
}

func (a *app) UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error) {
	ad, err := a.repo.GetById(adID)
	if err != nil {
		return ads.Ad{}, err
	}

	if ad.AuthorID != userID {
		return ads.Ad{}, AccessErr
	}

	ad.Title, ad.Text = title, text

	if err := validator.Validate(ad); err != nil {
		return ads.Ad{}, ValidationErr
	}

	return ad, a.repo.ReplaceByID(adID, ad)
}
