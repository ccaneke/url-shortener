package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ccaneke/url-shortner-poc/internal"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

func TestGetURL(t *testing.T) {
	testCases := []struct {
		name    string
		in      io.ReadCloser
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "no error",
			in:      io.NopCloser(strings.NewReader("https://en.wikipedia.org/wiki/Main_Page")),
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name:    "body is empty",
			in:      io.NopCloser(strings.NewReader("")),
			want:    &url.URL{},
			wantErr: false,
		},
		{
			name:    "error scheme is missing",
			in:      io.NopCloser(strings.NewReader("://en.wikipedia.org/wiki/Main_Page")),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no error when host is missing",
			in:      io.NopCloser(strings.NewReader("https:///wiki/Main_Page")),
			want:    &url.URL{Scheme: "https", Host: "", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name:    "no error when path is missing",
			in:      io.NopCloser(strings.NewReader("https://en.wikipedia.org/")),
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: ""},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		got, err := internal.GetURL(tc.in)

		var gotErr bool
		if err != nil {
			gotErr = true
		}

		if gotErr != tc.wantErr {
			fmt.Println(tc.name)
			t.Errorf("%v: getURL(%v)==%v, want %v", tc.name, tc.in, *got, *tc.want)
		}
	}
}

type RedisClientMock struct {
	mock.Mock
}

func (m *RedisClientMock) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func TestGetLongURL(t *testing.T) {
	testObj := new(RedisClientMock)
	ctx := context.Background()
	request := &http.Request{URL: &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"}}
	domain := "localhost"
	testObj.On("Get", ctx, "https://localhost/wiki/Main_Page").Return(&redis.StringCmd{})

	getLongURL(ctx, request, testObj, domain)
	testObj.AssertExpectations(t)

}
