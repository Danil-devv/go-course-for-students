package adrepo

import (
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

var wrongIdErr = errors.New("id must be non-negative and must be less than repository size")

func New() app.Repository {
	return &adRepo{repo: make([]ads.Ad, 0)}
}

type adRepo struct {
	repo []ads.Ad
}

func (r *adRepo) checkID(id int64) error {
	if id >= 0 && id <= int64(len(r.repo)-1) {
		return nil
	}

	return wrongIdErr
}

func (r *adRepo) AddAd(ad ads.Ad) int64 {
	r.repo = append(r.repo, ad)
	return int64(len(r.repo) - 1)
}

func (r *adRepo) GetById(id int64) (ads.Ad, error) {
	if err := r.checkID(id); err != nil {
		return ads.Ad{}, err
	}

	return r.repo[id], nil
}

func (r *adRepo) ReplaceByID(id int64, ad ads.Ad) error {
	if err := r.checkID(id); err != nil {
		return err
	}

	r.repo[id] = ad
	return nil
}

func (r *adRepo) GetSize() int64 {
	return int64(len(r.repo))
}
