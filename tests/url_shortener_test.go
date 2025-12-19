//go:build integration
// +build integration

package tests

import (
	"net/http"

	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/noctusha/url-shortener/internal/transport/http/shortenerhandler"
)

const (
	host      = "localhost:8080"
	basicUser = "username"
	basicPass = "password"
)

func newExpect(t *testing.T) *httpexpect.Expect {
	t.Helper()

	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  u.String(),
		Reporter: httpexpect.NewAssertReporter(t),
		// отключение стандартного редиректа из дефолтного клиента
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	})
}

func TestURLShortener_HappyPath(t *testing.T) {
	e := newExpect(t)

	alias := gofakeit.LetterN(10)
	originalURL := gofakeit.URL()

	createResp := e.POST("/url").
		WithBasicAuth(basicUser, basicPass).
		WithJSON(shortenerhandler.SaveRequest{
			URL:   originalURL,
			Alias: alias,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	createResp.Value("alias").String().IsEqual(alias)

	e.GET("/" + alias).
		Expect().
		Status(http.StatusFound).
		Header("Location").IsEqual(originalURL)

	e.DELETE("/url/"+alias).
		WithBasicAuth(basicUser, basicPass).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("alias").String().IsEqual(alias)

	e.GET("/" + alias).
		Expect().
		Status(http.StatusNotFound)
}

func TestURLShortener_Unauthorized(t *testing.T) {
	e := newExpect(t)

	e.POST("/url").
		WithJSON(shortenerhandler.SaveRequest{
			URL:   gofakeit.URL(),
			Alias: gofakeit.LetterN(10),
		}).
		Expect().
		Status(http.StatusUnauthorized)
}

func TestURLShortener_NotFound(t *testing.T) {
	e := newExpect(t)

	alias := gofakeit.LetterN(15)

	e.GET("/" + alias).
		Expect().
		Status(http.StatusNotFound)

	e.DELETE("/url/"+alias).
		WithBasicAuth(basicUser, basicPass).
		Expect().
		Status(http.StatusNotFound)
}
