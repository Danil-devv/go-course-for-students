package tests

import (
	"homework10/internal/ads"
	"homework10/internal/ports/httpgin"
	mockApp "homework10/internal/tests/mocks/app"
	"net/http/httptest"
	"testing"
)

func getTestMockClient(t *testing.T) *testClient {
	a := mockApp.NewApp(t)

	a.On("CreateAd", "hello", "world", int64(123)).Return(ads.Ad{
		ID:        int64(0),
		Title:     "hello",
		Text:      "world",
		Published: false,
		AuthorID:  int64(123),
	}, nil)

	a.On("GetAds").Return([]ads.Ad{{
		ID:        int64(0),
		Title:     "hello",
		Text:      "world",
		Published: false,
		AuthorID:  int64(123),
	}}, nil)

	a.On("GetAd", int64(0)).Return(ads.Ad{
		ID:        int64(0),
		Title:     "hello",
		Text:      "world",
		Published: false,
		AuthorID:  int64(123),
	}, nil)

	server := httpgin.NewHTTPServer(":18080", a)
	testServer := httptest.NewServer(server.Handler())

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}
