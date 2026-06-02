package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/middleware"
	"github.com/yushafro/effective-mobile-tz/pkg/test"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_Logging_WithID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockNext := new(test.MockHandler)

	mockNext.On("ServeHTTP", mock.Anything, mock.MatchedBy(func(r *http.Request) bool {
		id := logger.RequestIDFromContext(r.Context())
		if id == "" {
			return false
		}

		_, err := uuid.Parse(id)

		return err == nil
	}))

	mockLogger.EXPECT().Info("empty request id, creating new one")
	mockLogger.EXPECT().Zap().Return(zap.NewNop())
	mockLogger.EXPECT().Stop()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler := middleware.Logging(mockNext, mockLogger)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
}

func Test_Logging_WithoutID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockLogger := logger.NewMockLogger(ctrl)
	mockNext := new(test.MockHandler)

	headerID := uuid.NewString()

	mockNext.On("ServeHTTP", mock.Anything, mock.MatchedBy(func(r *http.Request) bool {
		ctxID := logger.RequestIDFromContext(r.Context())
		if ctxID == "" {
			return false
		}

		return ctxID == headerID
	}))
	mockLogger.EXPECT().Zap().Return(zap.NewNop())
	mockLogger.EXPECT().Stop()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(httputil.RequestID, headerID)

	handler := middleware.Logging(mockNext, mockLogger)
	handler.ServeHTTP(w, req)

	mockNext.AssertExpectations(t)
}
