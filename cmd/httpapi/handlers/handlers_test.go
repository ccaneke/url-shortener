package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) Print(v ...any) {
	_ = m.Called(v)
}

func (m *LoggerMock) Fatal(v ...any) {
	_ = m.Called(v)
}

func TestShortenUrl(t *testing.T) {
	var testcases = []struct {
		name string
		in   struct {
			url        url.URL
			domain     string
			uuidString string
			loggerMock *LoggerMock
		}
		want    string
		wantErr bool
	}{
		{
			name: "successful when domain is missing",
			in: struct {
				url        url.URL
				domain     string
				uuidString string
				loggerMock *LoggerMock
			}{
				url:        url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "wiki/Main_Page/123456"},
				domain:     "",
				uuidString: "123e45",
				loggerMock: new(LoggerMock),
			},
			want:    "",
			wantErr: false},

		{
			name: "error when uuid is greater than 10",
			in: struct {
				url        url.URL
				domain     string
				uuidString string
				loggerMock *LoggerMock
			}{
				url:        url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "wiki/Main_Page/123456"},
				domain:     "",
				uuidString: "123e4567-e89b-12d3-a456-426614174000",
				loggerMock: new(LoggerMock),
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		if tc.wantErr {
			tc.in.loggerMock.On("Print", mock.Anything).Return()
		}
		got, err := shortenURL(tc.in.url, tc.in.domain, tc.in.uuidString, tc.in.loggerMock)

		var gotErr bool
		if err != nil {
			gotErr = true
		}

		if gotErr != tc.wantErr {
			t.Errorf("%v: ShortenURL(%+v)==%v, want %v", tc.name, tc.in, got, tc.want)
		}

		tc.in.loggerMock.AssertExpectations(t)
	}
}

func TestUrlFromBody(t *testing.T) {
	testCases := []struct {
		name string
		in   struct {
			body   io.ReadCloser
			logger *LoggerMock
		}
		want    *url.URL
		wantErr bool
	}{
		{
			name: "valid url",
			in: struct {
				body   io.ReadCloser
				logger *LoggerMock
			}{
				body:   io.NopCloser(strings.NewReader(`{"LongURL":"https://en.wikipedia.org/wiki/Main_Page"}`)),
				logger: new(LoggerMock),
			},
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name: "invalid request, long url cannot be blank",
			in: struct {
				body   io.ReadCloser
				logger *LoggerMock
			}{
				body:   io.NopCloser(strings.NewReader("")),
				logger: new(LoggerMock),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		if tc.wantErr {
			tc.in.logger.On("Print", mock.Anything).Return()
		}
		got, err := urlFromBody(tc.in.body, tc.in.logger)

		var gotErr bool
		if err != nil {
			gotErr = true
		}

		if gotErr != tc.wantErr {
			fmt.Println(tc.name)
			t.Errorf("%v: getURL(%v)==%v, want %v", tc.name, tc.in, *got, *tc.want)
		}

		tc.in.logger.AssertExpectations(t)
	}
}

func TestGetDomain(t *testing.T) {
	testCases := []struct {
		name string
		in   *http.Request
		want string
	}{
		{
			name: "host includes domain and port number",
			in:   &http.Request{Host: "localhost:3000"},
			want: "localhost",
		},
		{
			name: "host does not include port number",
			in:   &http.Request{Host: "localhost"},
			want: "localhost"},
		{
			name: "host is empty",
			in:   &http.Request{Host: ""},
			want: "",
		},
	}

	for _, tc := range testCases {
		got := getDomain(tc.in)

		if got != tc.want {
			t.Errorf("%v: GetDomain(%v)==%v, want %v", tc.name, tc.in, got, tc.want)
		}
	}
}
