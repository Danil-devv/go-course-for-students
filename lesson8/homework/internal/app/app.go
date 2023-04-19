package app

import (
	"errors"
	validator "github.com/Danil-devv/structValidator"
	"homework8/internal/ads"
	"homework8/internal/users"
	"net/mail"
	"strings"
	"time"
)

var (
	ValidationErr = errors.New("some fields does not pass the validation")
	AccessErr     = errors.New("user can only change his ads")
)

type App interface {
	CreateAd(title string, text string, authorID int64) (ads.Ad, error)
	GetAds() ([]ads.Ad, error)
	ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error)
	UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error)
	GetAd(adID int64) (ads.Ad, error)
	GetAdsByTitle(title string) ([]ads.Ad, error)
	GetFilteredAds(published int, authorID int64, date string) ([]ads.Ad, error)
	CreateUser(id int64, nickname string, email string) (users.User, error)
	GetUser(id int64) (users.User, error)
	UpdateUser(id int64, nickname string, email string) (users.User, error)
}

func NewApp(adRepo ads.Repository, usersRepo users.Repository) App {
	return &app{adRepo: adRepo,
		usersRepo: usersRepo}
}

type app struct {
	adRepo    ads.Repository
	usersRepo users.Repository
}

func (a *app) CreateUser(id int64, nickname string, email string) (users.User, error) {
	_, err := mail.ParseAddress(email)
	if id < 0 || len(nickname) == 0 || err != nil {
		return users.User{}, ValidationErr
	}

	u := users.User{
		ID:       id,
		Nickname: nickname,
		Email:    email,
	}
	err = a.usersRepo.AddUser(u)
	if err != nil {
		return users.User{}, err
	}

	return u, nil
}

func (a *app) GetUser(id int64) (users.User, error) {
	u, err := a.usersRepo.GetById(id)
	if err != nil {
		return users.User{}, err
	}

	return u, nil
}

func (a *app) UpdateUser(id int64, nickname string, email string) (users.User, error) {
	u, err := a.usersRepo.GetById(id)
	if err != nil {
		return users.User{}, err
	}

	if nickname != "" {
		u.Nickname = nickname
	}

	if email != "" {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return users.User{}, ValidationErr
		}

		u.Email = email
	}

	return u, a.usersRepo.ReplaceByID(id, u)
}

func (a *app) CreateAd(title string, text string, authorID int64) (ads.Ad, error) {
	t := time.Now().UTC()
	ad := ads.Ad{ID: a.adRepo.GetSize(), Title: title, Text: text, AuthorID: authorID,
		CreateDate: t, LastUpdate: t}

	if err := validator.Validate(ad); err != nil {
		return ads.Ad{}, ValidationErr
	}

	a.adRepo.AddAd(ad)
	return ad, nil
}

func (a *app) GetAds() ([]ads.Ad, error) {
	res := make([]ads.Ad, 0)
	for i := int64(0); i < a.adRepo.GetSize(); i++ {
		r, err := a.adRepo.GetById(i)
		if err != nil {
			return []ads.Ad{}, err
		}
		if r.Published {
			res = append(res, r)
		}
	}
	return res, nil
}

func (a *app) GetAd(adID int64) (ads.Ad, error) {
	ad, err := a.adRepo.GetById(adID)

	if err != nil {
		return ads.Ad{}, err
	}

	if !ad.Published {
		return ads.Ad{}, AccessErr
	}
	return ad, nil
}

func (a *app) ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error) {
	t := time.Now().UTC()
	ad, err := a.adRepo.GetById(adID)
	if err != nil {
		return ads.Ad{}, err
	}

	if ad.AuthorID != userID {
		return ads.Ad{}, AccessErr
	}

	ad.Published, ad.LastUpdate = published, t
	return ad, a.adRepo.ReplaceByID(adID, ad)
}

func (a *app) UpdateAd(adID int64, userID int64, title string, text string) (ads.Ad, error) {
	t := time.Now().UTC()
	ad, err := a.adRepo.GetById(adID)
	if err != nil {
		return ads.Ad{}, err
	}

	if ad.AuthorID != userID {
		return ads.Ad{}, AccessErr
	}

	ad.Title, ad.Text, ad.LastUpdate = title, text, t

	if err := validator.Validate(ad); err != nil {
		return ads.Ad{}, ValidationErr
	}

	return ad, a.adRepo.ReplaceByID(adID, ad)
}

func (a *app) GetAdsByTitle(title string) ([]ads.Ad, error) {
	res := make([]ads.Ad, 0)

	for i := int64(0); i < a.adRepo.GetSize(); i++ {
		ad, _ := a.adRepo.GetById(i)
		if strings.Contains(ad.Title, title) && ad.Published {
			res = append(res, ad)
		}
	}

	return res, nil
}

func (a *app) GetFilteredAds(published int, authorID int64, date string) ([]ads.Ad, error) {
	res := make([]ads.Ad, 0)
	for i := int64(0); i < a.adRepo.GetSize(); i++ {
		r, err := a.adRepo.GetById(i)
		if err != nil {
			return []ads.Ad{}, err
		}

		if published == 1 && !r.Published {
			continue
		}

		if authorID != -1 && r.AuthorID != authorID {
			continue
		}

		if date != "" && date != r.CreateDate.Format(time.DateOnly) {
			continue
		}

		res = append(res, r)
	}
	return res, nil
}
