package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"github.com/yushafro/effective-mobile-tz/pkg/test"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_RequireBody_WithBody(t *testing.T) {
	t.Parallel()

	mockNext := new(test.MockHandler)
	mockNext.On("ServeHTTP", mock.Anything, mock.Anything)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test"))

	handler := middleware.RequireBody(mockNext)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
	require.Equal(t, http.StatusOK, w.Code)
}

func Test_RequireBody_WithoutBody(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockNext := new(test.MockHandler)

	mockLogger.EXPECT().
		Error("request body validation failed", zap.Error(middleware.ErrRequireBody))

	w := httptest.NewRecorder()
	ctx := logger.WithLogger(context.Background(), mockLogger)
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/", nil)

	handler := middleware.RequireBody(mockNext)
	handler.ServeHTTP(w, req)

	mockNext.AssertNotCalled(t, "ServeHTTP")
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), middleware.ErrRequireBody.Error())
}

func Test_RequireBody_ChunkedBody(t *testing.T) {
	t.Parallel()

	mockNext := new(test.MockHandler)
	mockNext.On("ServeHTTP", mock.Anything, mock.Anything)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.ContentLength = -1

	handler := middleware.RequireBody(mockNext)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
	require.Equal(t, http.StatusOK, w.Code)
}
