package internal

import (
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
		uCopy := *u
		if err != nil {
			return
		}

		var gotErr bool
		got, err := ShortenURL(u, tc.in.domain)
		if err != nil {
			gotErr = true
		}

		if gotErr {
			t.Errorf("testcase %v, ShortenURL(%v), error=%v", tc.name, uCopy.String(), err)
		} else if !strings.HasPrefix(got, tc.want) {
			t.Errorf("testcase %v, ShortenURL(%v)=%v, want %v", tc.name, uCopy.String(), got, tc.want)
		}

	}
}
