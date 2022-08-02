package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
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
		got, err := getURL(tc.in)

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

type redisClientMock struct {
}

func (r *redisClientMock) Get(context context.Context, key string) *redis.StringCmd {
	return &redis.StringCmd{}
}

func TestGetLongURL(t *testing.T) {
	wants := []string{"https://en.wikipedia.org/wiki/Main_Page"}
	testcases := []struct {
		name string
		in   struct {
			context     context.Context
			request     *http.Request
			redisClient RedisClient
			domain      string
		}
		want      *string
		wantError bool
	}{
		{
			name: "successfully gets a value",
			in: struct {
				context     context.Context
				request     *http.Request
				redisClient RedisClient
				domain      string
			}{
				context:     context.Background(),
				request:     &http.Request{ /*Host: "localhost:3000", */ URL: &url.URL{Scheme: "http", Host: "localhost:3000", Path: "/73a546"}},
				redisClient: &redisClientMock{},
				domain:      "localhost",
			},
			want:      &wants[0],
			wantError: false,
		},
	}

	for _, tc := range testcases {
		var gotError bool
		got, err := getLongURL(tc.in.context, tc.in.request, tc.in.redisClient, tc.in.domain)
		if err != nil {
			gotError = true
		}

		if gotError != tc.wantError {
			t.Errorf("%v: getLongURL(%+v)==%v, want %v", tc.name, tc.in, got, tc.want)
		}
	}
}
