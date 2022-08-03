package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestShortenUrl(t *testing.T) {
	var testcases = []struct {
		name string
		in   struct {
			url        url.URL
			domain     string
			uuidString string
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
			}{
				url:        url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "wiki/Main_Page/123456"},
				domain:     "",
				uuidString: "123e45",
			},
			want:    "",
			wantErr: false},

		{
			name: "error when uuid is greater than 10",
			in: struct {
				url        url.URL
				domain     string
				uuidString string
			}{
				url:        url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "wiki/Main_Page/123456"},
				domain:     "",
				uuidString: "123e4567-e89b-12d3-a456-426614174000",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		got, err := ShortenURL(tc.in.url, tc.in.domain, tc.in.uuidString)
		var gotErr bool
		if err != nil {
			gotErr = true
		}

		if gotErr != tc.wantErr {
			t.Errorf("%v: ShortenURL(%+v)==%v, want %v", tc.name, tc.in, got, tc.want)
		}
	}
}

func TestURLFromBody(t *testing.T) {
	testCases := []struct {
		name    string
		in      io.ReadCloser
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "successful when long url is correct",
			in:      io.NopCloser(strings.NewReader(`{"LongURL":"https://en.wikipedia.org/wiki/Main_Page"}`)),
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name:    "successful when scheme is missing",
			in:      io.NopCloser(strings.NewReader(`{"LongURL":"//en.wikipedia.org/wiki/Main_Page"}`)),
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name:    "no error when host is missing",
			in:      io.NopCloser(strings.NewReader(`{"LongURL":"https:///wiki/Main_Page"}`)),
			want:    &url.URL{Scheme: "https", Host: "", Path: "/wiki/Main_Page"},
			wantErr: false,
		},
		{
			name:    "no error when path is missing",
			in:      io.NopCloser(strings.NewReader(`{"LongURL":"https://en.wikipedia.org/"}`)),
			want:    &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: ""},
			wantErr: false,
		},
		{
			name:    "invalid request, long url cannot be blank",
			in:      io.NopCloser(strings.NewReader("")),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		got, err := URLFromBody(tc.in)

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
