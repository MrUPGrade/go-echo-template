package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type promEchoInstrumentation struct {
	RequestCounter *RequestCounter
	registry       *prometheus.Registry
}

type RequestCounter struct {
	Count    *prometheus.CounterVec
	Latency  *prometheus.GaugeVec
	BytesIn  *prometheus.GaugeVec
	BytesOut *prometheus.GaugeVec
}

func (r RequestCounter) AddValue(latency, bytesIn, bytesOut int, labels prometheus.Labels) {
	c, err := r.Count.GetMetricWith(labels)
	if err != nil {
		fmt.Println(err)
	}
	c.Inc()

	l, err := r.Latency.GetMetricWith(labels)
	if err != nil {
		fmt.Println(err)
	}
	l.Set(float64(latency))

	bi, err := r.BytesIn.GetMetricWith(labels)
	if err != nil {
		fmt.Println(err)
	}
	bi.Set(float64(bytesIn))

	bo, err := r.BytesOut.GetMetricWith(labels)
	if err != nil {
		fmt.Println(err)
	}
	bo.Set(float64(bytesOut))
}

func newRequestCounters() *RequestCounter {
	rc := &RequestCounter{
		Count: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "echo_http_count",
		}, []string{"method", "uri", "status_code"}),
		Latency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "echo_http_latency",
		}, []string{"method", "uri", "status_code"}),
		BytesIn: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "echo_http_bytes_in",
		}, []string{"method", "uri", "status_code"}),
		BytesOut: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "echo_http_bytes_out",
		}, []string{"method", "uri", "status_code"}),
	}

	return rc
}

func NewPromEchoInstrumentation() *promEchoInstrumentation {
	res := &promEchoInstrumentation{}

	res.RequestCounter = newRequestCounters()

	res.registry = prometheus.NewRegistry()
	res.registry.MustRegister(res.RequestCounter.Count)
	res.registry.MustRegister(res.RequestCounter.Latency)
	res.registry.MustRegister(res.RequestCounter.BytesIn)
	res.registry.MustRegister(res.RequestCounter.BytesOut)

	return res
}

func (p promEchoInstrumentation) PrometheusStatsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		var err error

		if err = next(c); err != nil {
			c.Error(err)
		}

		stop := time.Now()

		req := c.Request()
		res := c.Response()

		latency := stop.Sub(start)

		cl := req.Header.Get(echo.HeaderContentLength)
		var bytesIn int
		if cl == "" {
			bytesIn = 0
		} else {
			bytesIn, _ = strconv.Atoi(cl)
		}

		bytesOut := res.Size

		p.RequestCounter.AddValue(int(latency.Seconds()), bytesIn, int(bytesOut), prometheus.Labels{
			"status_code": fmt.Sprintf("%d", res.Status),
			"uri":         req.RequestURI,
			"method":      req.Method,
		})

		return err
	}
}

func (p promEchoInstrumentation) MetricsEndpoint(c echo.Context) error {
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
