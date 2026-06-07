package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	postHttpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blog_post_http_requests_total", // Note: job and instance are added automatically
			Help: "Number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	postHttpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "my_app_ns",
			Subsystem: "my_blog_post",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests [Blog Post]",
			Buckets:   []float64{0.005, 0.05, 0.1, 0.2},
		},
		[]string{"path", "method", "status", "post_id"},
	)

	//prometheus.MustRegister(httpRequestsTotal)
	//prometheus.MustRegister(httpRequestDuration)
	if err := prometheus.Register(postHttpRequestsTotal); err != nil {
		log.Fatal(err)
	}
	if err := prometheus.Register(postHttpRequestDuration); err != nil {
		log.Fatal(err)
	}

	srv := echo.New()
	srv.Use(middleware.RequestLogger())

	srv.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem: "my_app",
		Skipper: func(c *echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/blog/posts")
		},
		//AfterNext: func(c *echo.Context, err error) {
		//	path := c.Path()
		//	method := c.Request().Method
		//	status := http.StatusOK
		//	if err != nil {
		//		if he, ok := err.(*echo.HTTPError); ok {
		//			status = he.Code
		//		} else {
		//			status = http.StatusInternalServerError
		//		}
		//	}
		//	httpRequestsTotal.WithLabelValues(path, method, http.StatusText(status)).Inc()
		//},
	}))

	blogRouter := srv.Group("/blog/posts")
	// customs middleware for Prometheus metrics
	blogRouter.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()
			path := c.Path()
			method := c.Request().Method
			status := http.StatusOK
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else {
					status = http.StatusInternalServerError
				}
			}

			postHttpRequestsTotal.WithLabelValues(path, method, http.StatusText(status)).Inc()
			postHttpRequestDuration.WithLabelValues(path, method, http.StatusText(status), c.Param("id")).Observe(duration)

			return err
		}
	})
	blogRouter.GET("/:id", getPostItem)

	pRouter := srv.Group("/products")
	pRouter.GET("/:id", getProductItem)

	// expose metrics for scraping
	srv.GET("/metrics", echoprometheus.NewHandler())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, srv); err != nil {
		srv.Logger.Error("failed to start server", "error", err)
	}
}

func getProductItem(c *echo.Context) error {
	// do something heavy
	d := rand.Intn(500)
	if d > 490 {
		return echo.ErrNotFound
	}
	time.Sleep(time.Duration(d) * time.Millisecond)

	id := c.Param("id")
	priceExtra, _ := strconv.Atoi(id)
	return c.JSON(200, map[string]any{
		"id":   id,
		"name": "Product Name " + id,
		"price": map[string]any{
			"amount":   100 + priceExtra,
			"currency": "EUR",
		},
	})
}

func getPostItem(c *echo.Context) error {
	// do something heavy
	d := rand.Intn(500)
	if d > 495 {
		return echo.ErrNotFound
	}
	time.Sleep(time.Duration(d) * time.Millisecond)

	id := c.Param("id")
	return c.JSON(200, map[string]any{
		"id":    id,
		"title": "Post " + id,
	})
}
