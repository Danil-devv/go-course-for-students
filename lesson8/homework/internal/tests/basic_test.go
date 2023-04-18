package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(123, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(123, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(123, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestGetAdByID(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	ad, err := client.getAdByID(0)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, publishedAd.Data.ID)
	assert.Equal(t, ad.Data.Title, publishedAd.Data.Title)
	assert.Equal(t, ad.Data.Text, publishedAd.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ad.Data.Published)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.getAdByID(1)
	assert.Equal(t, err, ErrForbidden)
}

func TestGetAdsByTitle(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "best", "world")
	assert.NoError(t, err)

	_, err = client.createAd(123, "qq", "danil")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.getAdsByTitle("best")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)

	ads, err = client.getAdsByTitle("some title")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 0)

	ads, err = client.getAdsByTitle("qq")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 0)
}

func TestFilteredAds(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(123, "best", "world")
	assert.NoError(t, err)

	_, err = client.createAd(123, "qq", "danil")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(5, "test ad", "some text")
	assert.NoError(t, err)

	r, err := client.getFilteredAds(1, -1, "")
	assert.NoError(t, err)
	assert.Len(t, r.Data, 1)

	r, err = client.getFilteredAds(-1, -1, "")
	assert.NoError(t, err)
	assert.Len(t, r.Data, 4)

	r, err = client.getFilteredAds(-1, -1, time.Now().Format(time.DateOnly))
	assert.NoError(t, err)
	assert.Len(t, r.Data, 4)

	r, err = client.getFilteredAds(-1, -1, "2004-04-23")
	assert.NoError(t, err)
	assert.Len(t, r.Data, 0)

	r, err = client.getFilteredAds(-1, 5, "")
	assert.NoError(t, err)
	assert.Len(t, r.Data, 1)
}

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser(123, "danil", "mail@example.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "danil")
	assert.Equal(t, response.Data.Email, "mail@example.com")
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(123, "danil", "mail@example.com")
	assert.NoError(t, err)

	response, err := client.getUser(123)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Email, "mail@example.com")
	assert.Equal(t, response.Data.Nickname, "danil")
}
