package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
)

func AppRouter(r *gin.Engine, a app.App) {
	r.POST("/api/v1/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.GET("/api/v1/ads", getAds(a))                       // Метод для получения списка всех объявлений (ad)
	r.PUT("/api/v1/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/api/v1/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("api/v1/ads/:ad_id", getAdByID(a))              // Метод для получения объявления по его ID
	r.GET("api/v1/ads/find/:title", getAdsByTitle(a))     // Метод для получения списка объявлений по их заголовку
	r.GET("api/v1/ads/filter", getFilteredAds(a))         // Метод для получения списка отфильтрованных объявлений
	r.POST("/api/v1/users", createUser(a))                // Метод для создания пользователя (user)
	r.GET("/api/v1/users/:user_id", getUser(a))           // Метод для получения пользователя по id (user)
	r.POST("/api/v1/users/:user_id", updateUser(a))       // Метод для изменения пользователя по id (user)
	r.DELETE("/api/v1/users/:user_id", deleteUser(a))     // Метод для удаления пользователя по id (user)
	r.DELETE("/api/v1/ads/:ad_id", deleteAd(a))           // Метод для удаления объявления по id (ad)
}
