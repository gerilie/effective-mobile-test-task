package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"github.com/yushafro/effective-mobile-tz/pkg/test"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_NoBody_WithoutBody(t *testing.T) {
	t.Parallel()

	mockNext := new(test.MockHandler)
	mockNext.On("ServeHTTP", mock.Anything, mock.Anything)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler := middleware.NoBody(mockNext)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
	require.Equal(t, http.StatusOK, w.Code)
}

func Test_NoBody_WithBody(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	ctx := logger.WithLogger(context.Background(), mockLogger)
	mockNext := new(test.MockHandler)

	mockLogger.EXPECT().Error("request body validation failed", zap.Error(middleware.ErrNoBody))

	w := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.ContentLength = 1

	handler := middleware.NoBody(mockNext)
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), middleware.ErrNoBody.Error())
	mockNext.AssertNotCalled(t, "ServeHTTP")
}

func Test_NoBody_ChunkedBody(t *testing.T) {
	t.Parallel()

	mockNext := new(test.MockHandler)
	mockNext.On("ServeHTTP", mock.Anything, mock.Anything)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.ContentLength = -1

	handler := middleware.NoBody(mockNext)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
	require.Equal(t, http.StatusOK, w.Code)
}
