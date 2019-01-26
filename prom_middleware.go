package main

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"strings"
)

type PromEchoInstrumentation struct {
	TotalRequestCounter prometheus.Counter
	registry            *prometheus.Registry
}

func NewPromEchoInstrumentation() *PromEchoInstrumentation {
	res := &PromEchoInstrumentation{}

	res.TotalRequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "echoapi_http_total_count",
	})

	res.registry = prometheus.NewRegistry()
	res.registry.MustRegister(res.TotalRequestCounter)

	return res
}

func (p PromEchoInstrumentation) PrometheusStatsPushMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		p.TotalRequestCounter.Add(1)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

func (p PromEchoInstrumentation) MetricsEndpoint(c echo.Context) error {
	contentType := expfmt.Negotiate(c.Request().Header)
	metrics, err := p.registry.Gather()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	w := &strings.Builder{}
	enc := expfmt.NewEncoder(w, contentType)
	for _, metric := range metrics {
		if err := enc.Encode(metric); err != nil {
			c.Error(err)
			return err
		}
	}
	return c.String(http.StatusOK, w.String())
}
