package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PromPushResource struct {
	VisitCounter prometheus.Counter
	registry     *prometheus.Registry
}

func NewPromPushResource() PromPushResource {
	ppr := PromPushResource{}
	ppr.VisitCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "echoapi_prom_http_count",
		Help: "the of the counter that is pushed to prom gateway",
	})
	ppr.registry = prometheus.NewRegistry()
	ppr.registry.MustRegister(ppr.VisitCounter)

	pusher := push.New("http://localhost:21002", "echoapi").Gatherer(ppr.registry)
	go func() {
		for {
			time.Sleep(time.Second * 10)
			err := pusher.Add()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("metrics pushed to prometheus")
			}
		}
	}()
	return ppr
}

type JSONMessage struct {
	Message string
}

func (ctrl PromPushResource) get(c echo.Context) error {
	ctrl.VisitCounter.Add(1)
	return c.JSON(http.StatusOK, JSONMessage{"Everything ok"})
}
