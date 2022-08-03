package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

type RedisClientMock struct {
	mock.Mock
}

func (m *RedisClientMock) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func TestGetLongURL(t *testing.T) {
	ctx := context.Background()
	request := &http.Request{URL: &url.URL{Scheme: "https", Host: "en.wikipedia.org", Path: "/wiki/Main_Page"}}
	domain := "localhost"

	testCases := []struct {
		name             string
		mock             *RedisClientMock
		mockMethodInputs struct {
			ctx context.Context
			key string
		}
		mockMethodOutput *redis.StringCmd
	}{
		{
			name: "Value is empty",
			mock: new(RedisClientMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: context.Background(),
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
		},
		{
			name: "Key does not exist",
			mock: new(RedisClientMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
		},
		{
			name: "Get failed",
			mock: new(RedisClientMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
		},
		{
			name: "full value",
			mock: new(RedisClientMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
		},
	}

	for i, tc := range testCases {
		switch i {
		case 0:
			tc.mockMethodOutput.SetVal("")
		case 1:
			tc.mockMethodOutput.SetErr(redis.Nil)
		case 2:
			tc.mockMethodOutput.SetErr(errors.New("Get failed"))
		case 3:
			tc.mockMethodOutput.SetVal("https://en.wikipedia.org/wiki/Main_Page")
		}

		tc.mock.On("Get", tc.mockMethodInputs.ctx, tc.mockMethodInputs.key).Return(tc.mockMethodOutput)
		getLongURL(ctx, request, tc.mock, domain)
		tc.mock.AssertExpectations(t)
	}
}
