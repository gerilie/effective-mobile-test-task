package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"github.com/yushafro/effective-mobile-tz/pkg/ratelimiter"
	"github.com/yushafro/effective-mobile-tz/pkg/test"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func Test_RateLimiter_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockIPRateLimiter := ratelimiter.NewMockIPRateLimiter(ctrl)

	mockNext := new(test.MockHandler)
	mockNext.On("ServeHTTP", mock.Anything, mock.Anything)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	limiter := rate.NewLimiter(10, 100)
	mockIPRateLimiter.EXPECT().GetLimiter(httputil.GetClientIP(req)).Return(limiter)

	handler := middleware.RateLimiter(mockNext, mockIPRateLimiter)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "100", w.Header().Get("X-Ratelimit-Limit"))
	require.NotEmpty(t, w.Header().Get("X-Ratelimit-Remaining"))
}

func Test_RateLimiter_LimitExceeded(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockIPRateLimiter := ratelimiter.NewMockIPRateLimiter(ctrl)
	mockNext := new(test.MockHandler)

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Error("rate limit exceeded", zap.Error(middleware.ErrRateLimit))

	w := httptest.NewRecorder()

	ctx := logger.WithLogger(context.Background(), mockLogger)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	limiter := rate.NewLimiter(0, 0)
	mockIPRateLimiter.EXPECT().GetLimiter(httputil.GetClientIP(req)).Return(limiter)

	handler := middleware.RateLimiter(mockNext, mockIPRateLimiter)
	handler.ServeHTTP(w, req)

	mockNext.AssertNotCalled(t, "ServeHTTP")
	require.Equal(t, "0", w.Header().Get("X-Ratelimit-Limit"))
	require.NotEmpty(t, w.Header().Get("X-Ratelimit-Remaining"))

	require.Equal(t, http.StatusTooManyRequests, w.Code)
	require.Contains(t, w.Body.String(), middleware.ErrRateLimit.Error())
}

func Test_RateLimiter_GetClientIPFailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockIPRateLimiter := ratelimiter.NewMockIPRateLimiter(ctrl)
	mockNext := new(test.MockHandler)

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Error("internal middleware error", zap.Error(middleware.ErrGetClientIP))

	w := httptest.NewRecorder()
	ctx := logger.WithLogger(context.Background(), mockLogger)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.RemoteAddr = ""

	handler := middleware.RateLimiter(mockNext, mockIPRateLimiter)
	handler.ServeHTTP(w, req)

	mockNext.AssertNotCalled(t, "ServeHTTP")

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Contains(t, w.Body.String(), httputil.InternalServerErrorMsg)
}

func Test_RateLimiter_GetLimiterFailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockNext := new(test.MockHandler)

	mockIPRateLimiter := ratelimiter.NewMockIPRateLimiter(ctrl)
	mockIPRateLimiter.EXPECT().GetLimiter(gomock.Any()).Return(nil)

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Error("internal middleware error", zap.Error(middleware.ErrGetLimiter))

	w := httptest.NewRecorder()
	ctx := logger.WithLogger(context.Background(), mockLogger)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	handler := middleware.RateLimiter(mockNext, mockIPRateLimiter)
	handler.ServeHTTP(w, req)

	mockNext.AssertNotCalled(t, "ServeHTTP")

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Contains(t, w.Body.String(), httputil.InternalServerErrorMsg)
}
