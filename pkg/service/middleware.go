package service

import (
	"context"
	"time"

	"gokit_example/pkg/entity"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, p entity.JSONRequestProduct) (res entity.ResponseHttp, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Create Product", "took", time.Since(begin), "err", res.Message)
	}(time.Now())
	return mw.next.Create(ctx, p)
}

func (mw loggingMiddleware) Show(ctx context.Context) (res entity.ResponseHttp, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Show Products", "took", time.Since(begin), "err", res.Message)
	}(time.Now())
	return mw.next.Show(ctx)
}

func (mw loggingMiddleware) GetProduct(ctx context.Context, uid string) (res entity.ResponseHttp, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Get Product", "took", time.Since(begin), "err", res.Message)
	}(time.Now())
	return mw.next.GetProduct(ctx, uid)
}

func (mw loggingMiddleware) UpdateProduct(ctx context.Context, p entity.JSONRequestUpdateProduct) (res entity.ResponseHttp, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Update Product", "took", time.Since(begin), "err", res.Message)
	}(time.Now())
	return mw.next.UpdateProduct(ctx, p)
}
