package postgres_test

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/mock/gomock"
)

func Test_New_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockPoolConnector := postgres.NewMockPoolConnector(ctrl)

	ctx := logger.WithLogger(context.Background(), mockLogger)
	cfg := postgres.Config{}

	mockLogger.EXPECT().Info("connected to database")

	gomock.InOrder(
		mockPoolConnector.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(&pgxpool.Pool{}, nil),
		mockPoolConnector.EXPECT().Ping(gomock.Any(), gomock.Any()).Return(nil),
	)

	pool, err := postgres.New(ctx, cfg, mockPoolConnector)

	require.NoError(t, err)
	require.NotNil(t, pool)
}

func Test_New_ConnectionFailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockPoolConnector := postgres.NewMockPoolConnector(ctrl)

	ctx := logger.WithLogger(context.Background(), mockLogger)
	cfg := postgres.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(0)
	mockPoolConnector.EXPECT().
		Connect(gomock.Any(), gomock.Any()).
		Return(nil, postgres.ErrConnectionFailed)

	pool, err := postgres.New(ctx, cfg, mockPoolConnector)

	require.Nil(t, pool)
	require.Error(t, err)
	require.ErrorIs(t, err, postgres.ErrConnectionFailed)
}

func Test_New_PingFailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockPoolConnector := postgres.NewMockPoolConnector(ctrl)

	ctx := logger.WithLogger(context.Background(), mockLogger)
	cfg := postgres.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(0)
	gomock.InOrder(
		mockPoolConnector.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(&pgxpool.Pool{}, nil),
		mockPoolConnector.EXPECT().
			Ping(gomock.Any(), gomock.Any()).
			Return(postgres.ErrConnectionFailed),
		mockPoolConnector.EXPECT().Close(&pgxpool.Pool{}),
	)

	pool, err := postgres.New(ctx, cfg, mockPoolConnector)

	require.Nil(t, pool)
	require.Error(t, err)
	require.ErrorIs(t, err, postgres.ErrConnectionFailed)
}
