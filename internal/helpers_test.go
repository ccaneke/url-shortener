package internal

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestShortenUrl(t *testing.T) {
	var testcases = []struct {
		name string
		in   struct {
			rawUrl string
			domain string
		}
		want    string
		wantErr bool
	}{
		{
			name: "Should return no error when there's no query",
			in: struct {
				rawUrl string
				domain string
			}{
				rawUrl: "https://en.wikipedia.org/wiki/Main_Page/123456",
				domain: "short.io",
			},
			want: "https://short.io/"},
		{
			name: "Should return no error when there is a query",
			in: struct {
				rawUrl string
				domain string
			}{
				rawUrl: "https://en.wikipedia.org/wiki/Main_Page/123456?q=golang",
				domain: "short.io",
			},
			want: "https://short.io/",
		},
	}

	for _, tc := range testcases {
		u, err := url.Parse(tc.in.rawUrl)
		copy := *u
		if err != nil {
			return
		}

		var gotErr bool
		got, err := ShortenURL(copy, tc.in.domain)
		if err != nil {
			gotErr = true
		}

		if gotErr {
			t.Errorf("testcase %v, ShortenURL(%v), error=%v", tc.name, copy.String(), err)
		} else if !strings.HasPrefix(got, tc.want) {
			t.Errorf("testcase %v, ShortenURL(%v)=%v, want %v", tc.name, copy.String(), got, tc.want)
		}

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
		got := GetDomain(tc.in)

		if got != tc.want {
			t.Errorf("%v: GetDomain(%v)==%v, want %v", tc.name, tc.in, got, tc.want)
		}
	}
}
