package internal

import (
	"net/http"
	"net/url"
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
			name: "domain missing",
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
