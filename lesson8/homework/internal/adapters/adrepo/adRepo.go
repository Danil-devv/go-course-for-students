package adrepo

import (
	"errors"
	"homework8/internal/ads"
	"sync"
)

var wrongIdErr = errors.New("id must be non-negative and must be less than repository size")

func New() ads.Repository {
	return &adRepo{repo: make([]ads.Ad, 0)}
}

type adRepo struct {
	repo []ads.Ad
	m    sync.RWMutex
}

func (r *adRepo) checkID(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	if id >= 0 && id <= int64(len(r.repo)-1) {
		return nil
	}

	return wrongIdErr
}

func (r *adRepo) AddAd(ad ads.Ad) int64 {
	r.m.Lock()
	r.repo = append(r.repo, ad)
	r.m.Unlock()

	return int64(len(r.repo) - 1)
}

func (r *adRepo) GetById(id int64) (ads.Ad, error) {
	if err := r.checkID(id); err != nil {
		return ads.Ad{}, err
	}

	r.m.Lock()
	res := r.repo[id]
	r.m.Unlock()

	return res, nil
}

func (r *adRepo) ReplaceByID(id int64, ad ads.Ad) error {
	if err := r.checkID(id); err != nil {
		return err
	}

	r.m.Lock()
	r.repo[id] = ad
	r.m.Unlock()

	return nil
}

func (r *adRepo) GetSize() int64 {
	r.m.Lock()
	size := int64(len(r.repo))
	r.m.Unlock()

	return size
}
