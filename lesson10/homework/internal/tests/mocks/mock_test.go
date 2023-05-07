package mocks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGetAd(t *testing.T) {
	client := getTestClient(t)

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)

	response, err = client.getAdByID(0)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(123))
	assert.False(t, response.Data.Published)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Zero(t, ads.Data[0].ID)
	assert.Equal(t, ads.Data[0].Title, "hello")
	assert.Equal(t, ads.Data[0].Text, "world")
	assert.Equal(t, ads.Data[0].AuthorID, int64(123))

}
