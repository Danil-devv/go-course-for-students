package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/ads"
	"homework8/internal/users"
)

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserID   int64  `json:"user_id"`
}

type userResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserID   int64  `json:"user_id"`
}

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func UserSuccessResponse(u *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			UserID:   u.ID,
			Email:    u.Email,
			Nickname: u.Nickname,
		},
		"error": nil,
	}
}

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(ads *[]ads.Ad) *gin.H {
	res := make([]adResponse, 0)
	for _, ad := range *ads {
		res = append(res, adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		})
	}
	return &gin.H{
		"data":  res,
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
