package loopback

import (
	"math"
	"math/rand"
	randv2 "math/rand/v2"
	"slices"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"gonum.org/v1/gonum/stat/distuv"
)

func metrics() []metric {
	var allMetrics []metric

	for ls := range product(
		[]string{"staging", "production", "rc"},
		[]string{"us-central", "us-west"},
	) {
		namespace := ls[0]
		cluster := ls[1]
		for _, service := range []string{"orders", "payments"} {

			allMetrics = append(
				allMetrics,
				counterWithLabels(
					"container_cpu_usage_seconds_total",
					labels.Label{Name: "service", Value: service},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				),
				gaugeWithLabels(

					"container_memory_usage_bytes",
					labels.Label{Name: "service", Value: service},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				),
			)
		}
		// From
		// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/providers/prometheus/server_metrics.go
		// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/providers/prometheus/client_metrics.go
		// orders#OrdersService.CreateOrder => payments#PaymentsService.CreatePayment
		for _, code := range []string{"OK", "DeadlineExceeded"} {
			var errorRateParams sampleParams
			if code == "DeadlineExceeded" {
				// Per sec
				errorRateParams = &poissonSampleParams{rate: 0.1}
			}

			allMetrics = append(
				allMetrics,
				counterWithLabels(
					"grpc_server_handled_total",
					labels.Label{Name: "service", Value: "orders"},
					// ListOrders has no errors, whereas CreateOrder can have errors
					labels.Label{Name: "grpc_method", Value: "OrdersService.ListOrders"},
					labels.Label{Name: "grpc_code", Value: code},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				),
				counterWithLabels(
					"grpc_server_handled_total",
					labels.Label{Name: "service", Value: "orders"},
					labels.Label{Name: "grpc_method", Value: "OrdersService.CreateOrder"},
					labels.Label{Name: "grpc_code", Value: code},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				).withSampleParams(errorRateParams),
				counterWithLabels(
					"grpc_client_handled_total",
					labels.Label{Name: "service", Value: "orders"},
					labels.Label{Name: "grpc_method", Value: "PaymentsService.CreatePayment"},
					labels.Label{Name: "grpc_code", Value: code},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				).withSampleParams(errorRateParams),
				counterWithLabels(
					"grpc_server_handled_total",
					labels.Label{Name: "service", Value: "payments"},
					labels.Label{Name: "grpc_method", Value: "PaymentsService.CreatePayment"},
					labels.Label{Name: "grpc_code", Value: code},
					labels.Label{Name: "namespace", Value: namespace},
					labels.Label{Name: "cluster", Value: cluster},
				).withSampleParams(errorRateParams),
			)
		}
	}

	return allMetrics
}

func gaugeWithLabels(name string, additionalLabels ...labels.Label) metric {
	ls := labels.New(slices.Concat(labels.New(
		labels.Label{
			Name:  labels.MetricName,
			Value: name,
		}), additionalLabels)...)
	return metric{
		isGauge:      true,
		labels:       ls,
		sampleParams: randomSineSampleParams(ls),
	}
}

func counterWithLabels(name string, additionalLabels ...labels.Label) metric {
	ls := labels.New(slices.Concat(labels.New(
		labels.Label{
			Name:  labels.MetricName,
			Value: name,
		}), additionalLabels)...)
	return metric{
		labels:       ls,
		sampleParams: randomSineSampleParams(ls),
	}
}

type sampleParams interface {
	value(t time.Time) float64
}

// sineSampleParams encapsulate parameters to generate a fake time series that looks like a sine wave.
// Formula is scale * sin((2pi/period) * t + phaseDeg) + bias + scale/2
type sineSampleParams struct {
	scale    float64
	phaseDeg float64
	period   time.Duration
	bias     float64
}

func (p sineSampleParams) value(t time.Time) float64 {
	return p.scale*math.Sin(2*math.Pi/float64(p.period.Seconds())*float64(t.Unix())+(p.phaseDeg*math.Pi/180.0)) + p.bias + p.scale
}

// poissonSampleParams encapsulate parameters to generate a fake time series that looks like a sine wave.
// Formula is scale * sin((2pi/period) * t + phaseDeg) + bias + scale/2
type poissonSampleParams struct {
	rate float64
}

func (p poissonSampleParams) value(t time.Time) float64 {
	rnd := distuv.Poisson{
		Lambda: p.rate,
		Src:    randv2.NewPCG(uint64(t.Unix()), 0),
	}
	return rnd.Rand()
}

type metric struct {
	labels       labels.Labels
	sampleParams sampleParams
	isGauge      bool
}

func (m metric) withSampleParams(params sampleParams) metric {
	if params == nil {
		return m
	}
	m.sampleParams = params
	return m
}

func randomSineSampleParams(ls labels.Labels) sampleParams {
	rnd := rand.New(rand.NewSource(int64(ls.Hash())))
	return &sineSampleParams{
		scale:    choice(rnd, []float64{1.0, 2.0, 3.0, 10.0}),
		phaseDeg: choice(rnd, []float64{30.0, 60.0, 120.0}),
		period:   choice(rnd, []time.Duration{30 * time.Minute, 5 * time.Minute, 150 * time.Minute}),
		bias:     choice(rnd, []float64{0.0, 0.5, 4.0}),
	}
}
