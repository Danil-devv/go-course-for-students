package tests

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCreateAd(b *testing.B) {
	client := getTestClient()

	for i := 0; i < b.N; i++ {
		response, err := client.createAd(int64(i), "hello", "world")
		assert.NoError(b, err)
		assert.Equal(b, int64(i), response.Data.ID)
		assert.Equal(b, response.Data.Title, "hello")
		assert.Equal(b, response.Data.Text, "world")
		assert.Equal(b, response.Data.AuthorID, int64(i))
		assert.False(b, response.Data.Published)
	}
}

func BenchmarkChangeAdStatus(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		response, err = client.changeAdStatus(123, response.Data.ID, true)
		assert.NoError(b, err)
		assert.True(b, response.Data.Published)
	}
}

func BenchmarkUpdateAd(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		response, err = client.updateAd(123, response.Data.ID, "привет"+s, "мир")
		assert.NoError(b, err)
		assert.Equal(b, response.Data.Title, "привет"+s)
		assert.Equal(b, response.Data.Text, "мир")
	}
}

func BenchmarkListAds(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(int64(11), "hello", "world")
	assert.NoError(b, err)

	_, err = client.changeAdStatus(int64(11), response.Data.ID, true)
	assert.NoError(b, err)

	_, err = client.createAd(int64(11), "best cat", "not for sale")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		ads, err := client.listAds()
		assert.NoError(b, err)
		assert.True(b, ads.Data[0].Published)
	}
}

func BenchmarkGetAdByID(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(b, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		ad, err := client.getAdByID(0)
		assert.NoError(b, err)
		assert.Equal(b, ad.Data.ID, publishedAd.Data.ID)
		assert.Equal(b, ad.Data.Title, publishedAd.Data.Title)
		assert.Equal(b, ad.Data.Text, publishedAd.Data.Text)
		assert.Equal(b, ad.Data.AuthorID, publishedAd.Data.AuthorID)
		assert.True(b, ad.Data.Published)
	}

}

func BenchmarkGetAdsByTitle(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(123, "best", "world")
	assert.NoError(b, err)

	_, err = client.createAd(123, "qq", "danil")
	assert.NoError(b, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(b, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		ads, err := client.getAdsByTitle("best")
		assert.NoError(b, err)
		assert.Len(b, ads.Data, 1)
		assert.Equal(b, ads.Data[0].ID, publishedAd.Data.ID)
		assert.Equal(b, ads.Data[0].Title, publishedAd.Data.Title)
		assert.Equal(b, ads.Data[0].Text, publishedAd.Data.Text)
		assert.Equal(b, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
		assert.True(b, ads.Data[0].Published)
	}
}

func BenchmarkFilteredAds(b *testing.B) {
	client := getTestClient()

	response, err := client.createAd(123, "best", "world")
	assert.NoError(b, err)

	_, err = client.createAd(123, "qq", "danil")
	assert.NoError(b, err)

	_, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(b, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(b, err)

	_, err = client.createAd(5, "test ad", "some text")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		r, err := client.getFilteredAds(-1, -1, time.Now().UTC().Format(time.DateOnly))
		assert.NoError(b, err)
		assert.Len(b, r.Data, 4)
	}
}

func BenchmarkCreateUser(b *testing.B) {
	client := getTestClient()

	for i := 0; i < b.N; i++ {
		response, err := client.createUser(int64(i), "danil", "mail@example.com")
		assert.NoError(b, err)
		assert.Equal(b, response.Data.Nickname, "danil")
		assert.Equal(b, response.Data.Email, "mail@example.com")
	}
}

func BenchmarkGetUser(b *testing.B) {
	client := getTestClient()

	_, err := client.createUser(123, "danil", "mail@example.com")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		response, err := client.getUser(123)
		assert.NoError(b, err)
		assert.Equal(b, response.Data.Email, "mail@example.com")
		assert.Equal(b, response.Data.Nickname, "danil")
	}
}

func BenchmarkDeleteUser(b *testing.B) {
	client := getTestClient()

	for i := 0; i < b.N; i++ {
		_, err := client.deleteUser(int64(i))
		assert.Error(b, err)
	}

}

func BenchmarkDeleteAd(b *testing.B) {
	client := getTestClient()

	_, err := client.createUser(123, "danil", "mail@example.com")
	assert.NoError(b, err)

	_, err = client.createAd(123, "best", "world")
	assert.NoError(b, err)

	for i := 1; i < b.N; i++ {
		_, err = client.deleteAd(int64(i), 123)
		assert.Error(b, err)
	}
}
