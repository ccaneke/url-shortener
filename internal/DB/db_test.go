package db

import (
	"context"
	"errors"
	"testing"
	"time"

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

func (m *RedisClientMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) Print(v ...any) {
	_ = m.Called(v)
}

func (m *LoggerMock) Fatal(v ...any) {
	_ = m.Called(v)
}

func TestGetValue(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name             string
		redisMock        *RedisClientMock
		loggerMock       *LoggerMock
		mockMethodInputs struct {
			ctx context.Context
			key string
		}
		mockMethodOutput *redis.StringCmd
		val              string
		err              error
	}{
		{
			name:       "Value is empty",
			redisMock:  new(RedisClientMock),
			loggerMock: new(LoggerMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: context.Background(),
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
			val:              "",
			err:              nil,
		},
		{
			name:       "Key does not exist",
			redisMock:  new(RedisClientMock),
			loggerMock: new(LoggerMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
			val:              "",
			err:              redis.Nil,
		},
		{
			name:       "Get failed",
			redisMock:  new(RedisClientMock),
			loggerMock: new(LoggerMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
			val:              "",
			err:              errors.New("Get failed"),
		},
		{
			name:       "full value",
			redisMock:  new(RedisClientMock),
			loggerMock: new(LoggerMock),
			mockMethodInputs: struct {
				ctx context.Context
				key string
			}{
				ctx: ctx,
				key: "https://localhost/wiki/Main_Page",
			},
			mockMethodOutput: &redis.StringCmd{},
			val:              "https://en.wikipedia.org/wiki/Main_Page",
			err:              nil,
		},
	}

	for _, tc := range testCases {
		tc.mockMethodOutput.SetVal(tc.val)
		tc.mockMethodOutput.SetErr(tc.err)

		tc.redisMock.On("Get", tc.mockMethodInputs.ctx, tc.mockMethodInputs.key).Return(tc.mockMethodOutput)
		tc.loggerMock.On("Print", mock.Anything).Return()
		GetValue(ctx, tc.mockMethodInputs.key, tc.redisMock, tc.loggerMock)
		tc.redisMock.AssertExpectations(t)
	}
}
