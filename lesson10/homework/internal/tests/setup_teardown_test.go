package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestSuite struct {
	suite.Suite
	client *testClient
}

func (suite *TestSuite) SetupSuite() {
	suite.client = getTestClient()

	response, err := suite.client.createAd(123, "hello", "world")
	suite.Suite.NoError(err)
	_, err = suite.client.changeAdStatus(123, response.Data.ID, true)
	suite.Suite.NoError(err)

	response, err = suite.client.createAd(123, "best cat", "not for sale")
	suite.Suite.NoError(err)
	_, err = suite.client.changeAdStatus(123, response.Data.ID, true)
	suite.Suite.NoError(err)

	response, err = suite.client.createAd(123, "best dog", "not for sale")
	suite.Suite.NoError(err)
	_, err = suite.client.changeAdStatus(123, response.Data.ID, true)
	suite.Suite.NoError(err)

	response, err = suite.client.createAd(123, "some title", "some text")
	suite.Suite.NoError(err)
	_, err = suite.client.changeAdStatus(123, response.Data.ID, true)
	suite.Suite.NoError(err)
}

func (suite *TestSuite) TearDownSuite() {
	for i := 0; i < 4; i++ {
		_, err := suite.client.deleteAd(0, 123)
		assert.NoError(suite.T(), err)
	}

	suite.client.client.CloseIdleConnections()
}

func (suite *TestSuite) TestGetAd() {
	ad, err := suite.client.getAdByID(2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), ad.Data.Text, "not for sale")
	assert.Equal(suite.T(), ad.Data.Title, "best dog")
	assert.Equal(suite.T(), ad.Data.AuthorID, int64(123))
	assert.True(suite.T(), ad.Data.Published)
}

func (suite *TestSuite) TestListAds() {
	ad, err := suite.client.listAds()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), ad.Data, 4)
	assert.Equal(suite.T(), ad.Data[2].Text, "not for sale")
	assert.Equal(suite.T(), ad.Data[2].Title, "best dog")
	assert.Equal(suite.T(), ad.Data[2].AuthorID, int64(123))
	assert.True(suite.T(), ad.Data[2].Published)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
