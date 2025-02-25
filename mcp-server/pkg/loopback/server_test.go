package loopback

import (
	"context"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat"
)

func TestQueryRange_CPU(t *testing.T) {
	svr := NewPrometheusV1Server()

	v := queryRange(t, svr, "sum by (cluster) (rate(container_cpu_usage_seconds_total{}[2m]))")
	clusters := extractLabelValues(t, v, "cluster")
	require.ElementsMatch(t, []string{"us-central", "us-west"}, clusters)
}

func TestQueryRange_Memory(t *testing.T) {
	svr := NewPrometheusV1Server()

	v := queryRange(t, svr, "sum by (cluster) (container_memory_usage_bytes)")
	clusters := extractLabelValues(t, v, "cluster")
	require.ElementsMatch(t, []string{"us-central", "us-west"}, clusters)
}

func TestQueryRange_gRPC(t *testing.T) {
	svr := NewPrometheusV1Server()

	v := queryRange(t, svr, `sum by (grpc_method) (rate(grpc_server_handled_total{service="orders"}[2m]))`)
	methods := extractLabelValues(t, v, "grpc_method")
	require.ElementsMatch(t, []string{"OrdersService.ListOrders", "OrdersService.CreateOrder"}, methods)
}

func queryRange(t *testing.T, svr *PrometheusV1Server, query string) model.Value {
	start, err := time.Parse(time.RFC3339, "2025-04-01T00:00:00Z")
	require.NoError(t, err)
	end := start.Add(1 * time.Hour)
	v, warn, err := svr.QueryRange(
		context.Background(),
		query,
		v1.Range{
			Start: start,
			End:   end,
			Step:  time.Minute,
		},
	)
	require.NoError(t, err)
	require.Empty(t, warn)
	return v
}

func extractLabelValues(t *testing.T, v model.Value, labelName string) []string {
	mat, ok := v.(model.Matrix)
	require.True(t, ok)
	require.Len(t, mat, 2)

	var values []string
	for _, entry := range mat {
		require.Len(t, entry.Metric, 1)
		cluster, ok := entry.Metric[model.LabelName(labelName)]
		require.True(t, ok)
		values = append(values, string(cluster))

		// 1 hour with 1 minute spacing and inclusive ranges: 61 samples
		require.Len(t, entry.Values, 61)
	}
	return values
}

func TestLabelNamesValues(t *testing.T) {
	svr := NewPrometheusV1Server()
	start, err := time.Parse(time.RFC3339, "2025-04-01T00:00:00Z")
	require.NoError(t, err)
	end := start.Add(4 * time.Hour)

	names, warn, err := svr.LabelNames(context.Background(), nil, time.Time{}, time.Time{})
	require.NoError(t, err)
	require.Empty(t, warn)
	require.ElementsMatch(t, []string{"__name__", "cluster", "grpc_code", "namespace", "service", "grpc_method"}, names)

	names, warn, err = svr.LabelNames(context.Background(), []string{"container_cpu_usage_seconds_total{}"}, start, end)
	require.NoError(t, err)
	require.Empty(t, warn)
	require.ElementsMatch(t, []string{"__name__", "cluster", "namespace", "service"}, names)

	values, warn, err := svr.LabelValues(context.Background(), "cluster", nil, time.Time{}, time.Time{})
	require.NoError(t, err)
	require.Empty(t, warn)
	require.ElementsMatch(t, model.LabelValues{"us-central", "us-west"}, values)
}

func TestPoisson(t *testing.T) {
	p := &poissonSampleParams{rate: 4.0}
	start, err := time.Parse(time.RFC3339, "2025-04-01T00:00:00Z")
	require.NoError(t, err)
	end := start.Add(4 * time.Hour)
	step := 1 * time.Second

	var vals []float64
	cur := start
	for cur.Before(end) {
		vals = append(vals, p.value(cur))
		cur = cur.Add(step)
	}
	require.InDelta(t, 4.0, stat.Mean(vals, nil), 0.1)
	require.InDelta(t, 2.0, stat.StdDev(vals, nil), 0.1)
}
