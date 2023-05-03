package httpgin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func handleErr(err error) int {
	switch err {
	case app.ValidationErr:
		return http.StatusBadRequest
	case app.AccessErr:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод получения всех опубликованных объявлений
func getAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		ads, err := a.GetAds()
		fmt.Println(ads)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(&ads))
	}
}

// Метод получения объявления по его ID
func getAdByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.GetAd(int64(adID))

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для получения объявлений по названию
func getAdsByTitle(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		title := c.Param("title")

		ads, err := a.GetAdsByTitle(title)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(&ads))
	}
}

// Метод для получения отфильтрованного списка объявлений
func getFilteredAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			published = 0
			authorID  = -1
			date      = ""
			err       error
		)

		if c.Query("published") != "" {
			published, err = strconv.Atoi(c.Query("published"))
			if err != nil {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			}
		}

		if c.Query("author") != "" {
			authorID, err = strconv.Atoi(c.Query("author"))
			if err != nil {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			}
		}

		if c.Query("date") != "" {
			d, err := time.Parse(time.DateOnly, c.Query("date"))
			if err != nil {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			}
			date = d.Format(time.DateOnly)
		}

		ads, err := a.GetFilteredAds(published, int64(authorID), date)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(&ads))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для создания пользователя (user)
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		user, err := a.CreateUser(reqBody.UserID, reqBody.Nickname, reqBody.Email)

		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

// метод для получения пользователя по id
func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.GetUser(int64(id))
		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

// метод для обновления пользователя по id
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		id, err := strconv.Atoi(c.Query("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.UpdateUser(int64(id), reqBody.Nickname, reqBody.Email)
		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

// Метод для удаления пользователя (user)
func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.DeleteUser(int64(id))
		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

// Метод для удаления объявления (ad)
func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdResponse
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.DeleteAd(reqBody.ID, reqBody.AuthorID)
		if err != nil {
			c.JSON(handleErr(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}
