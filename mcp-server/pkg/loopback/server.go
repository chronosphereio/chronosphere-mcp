package loopback

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
)

type PrometheusV1Server struct {
	promQLEngine *promql.Engine

	queryable *fakePromStorageQueryable
	querier   *fakePromStorageQuerier
}

func NewPrometheusV1Server() *PrometheusV1Server {
	querier := newFakePromStorageQuerier()
	return &PrometheusV1Server{
		promQLEngine: promql.NewEngine(promql.EngineOpts{
			MaxSamples: 100000,
			Timeout:    120 * time.Second,
		}),
		querier:   querier,
		queryable: newFakePromStorageQueryable(querier),
	}
}

func (p *PrometheusV1Server) QueryRange(ctx context.Context, query string, r v1.Range, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	q, err := p.promQLEngine.NewRangeQuery(ctx, p.queryable, nil, query, r.Start, r.End, r.Step)
	if err != nil {
		return nil, v1.Warnings{}, err
	}
	res := q.Exec(ctx)
	if res.Err != nil {
		return nil, v1.Warnings{}, res.Err
	}

	m, err := toMatrix(res)
	if err != nil {
		return nil, v1.Warnings{}, err
	}
	return m, v1.Warnings{}, nil
}

func (p *PrometheusV1Server) LabelNames(ctx context.Context, matches []string, startTime, endTime time.Time, opts ...v1.Option) ([]string, v1.Warnings, error) {
	var labelNames []string
	if len(matches) > 0 {
		for _, match := range matches {
			res, _, err := p.QueryRange(ctx, match, v1.Range{Start: startTime, End: endTime, Step: 1 * time.Minute})
			if err != nil {
				return nil, v1.Warnings{}, err
			}
			mat, ok := res.(model.Matrix)
			if !ok {
				return nil, v1.Warnings{}, fmt.Errorf("not a matrix")
			}
			for _, s := range mat {
				for labelName := range maps.Keys(s.Metric) {
					labelNames = append(labelNames, string(labelName))
				}
			}
		}
	} else {
		for _, m := range p.querier.allMetrics {
			for _, label := range m.labels {
				labelNames = append(labelNames, label.Name)
			}
		}
	}

	labelNames = dedupe(labelNames)
	slices.Sort(labelNames)
	return labelNames, v1.Warnings{}, nil
}

func (p *PrometheusV1Server) LabelValues(ctx context.Context, label string, matches []string, startTime, endTime time.Time, opts ...v1.Option) (model.LabelValues, v1.Warnings, error) {
	var labelValues []model.LabelValue
	if len(matches) > 0 {
		for _, match := range matches {
			res, _, err := p.QueryRange(ctx, match, v1.Range{Start: startTime, End: endTime, Step: 1 * time.Minute})
			if err != nil {
				return nil, v1.Warnings{}, err
			}
			mat, ok := res.(model.Matrix)
			if !ok {
				return nil, v1.Warnings{}, fmt.Errorf("not a matrix")
			}
			for _, s := range mat {
				labelValue, ok := s.Metric[model.LabelName(label)]
				if ok {
					labelValues = append(labelValues, labelValue)
					break
				}
			}
		}
	} else {
		for _, m := range p.querier.allMetrics {
			for _, l := range m.labels {
				if l.Name == label {
					labelValues = append(labelValues, model.LabelValue(l.Value))
					break
				}
			}
		}
	}

	labelValues = dedupe(labelValues)
	slices.Sort(labelValues)
	return labelValues, v1.Warnings{}, nil
}

func (p *PrometheusV1Server) Query(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	q, err := p.promQLEngine.NewInstantQuery(ctx, p.queryable, nil, query, ts)
	if err != nil {
		return nil, v1.Warnings{}, err
	}
	res := q.Exec(ctx)
	if res.Err != nil {
		return nil, v1.Warnings{}, err
	}

	vec, err := toVector(res)
	if err != nil {
		return nil, v1.Warnings{}, err
	}
	return vec, v1.Warnings{}, nil
}

func (p *PrometheusV1Server) Series(ctx context.Context, matches []string, startTime, endTime time.Time, opts ...v1.Option) ([]model.LabelSet, v1.Warnings, error) {
	labelSets := []model.LabelSet{}

	for _, match := range matches {
		res, _, err := p.QueryRange(ctx, match, v1.Range{Start: startTime, End: endTime, Step: 1 * time.Minute})
		if err != nil {
			return nil, v1.Warnings{}, err
		}
		mat, ok := res.(model.Matrix)
		if !ok {
			return nil, v1.Warnings{}, fmt.Errorf("not a matrix")
		}
		for _, s := range mat {
			labelSets = append(labelSets, model.LabelSet(s.Metric))
		}
	}

	// Sort for some determinism
	slices.SortFunc(labelSets, func(a, b model.LabelSet) int {
		return cmp.Compare(a.FastFingerprint(), b.FastFingerprint())
	})
	return labelSets, v1.Warnings{}, nil
}

func (p *PrometheusV1Server) Metadata(context.Context, string, string) (map[string][]v1.Metadata, error) {
	return map[string][]v1.Metadata{}, nil
}

func toMatrix(res *promql.Result) (model.Matrix, error) {
	matrix, err := res.Matrix()
	if err != nil {
		return nil, err
	}
	ret := make([]*model.SampleStream, 0, matrix.Len())
	for _, s := range matrix {
		m := model.Metric{}
		var values []model.SamplePair
		for _, label := range s.Metric {
			m[model.LabelName(label.Name)] = model.LabelValue(label.Value)
		}
		for _, fp := range s.Floats {
			values = append(values, model.SamplePair{
				Timestamp: model.Time(fp.T),
				Value:     model.SampleValue(fp.F),
			})
		}
		ret = append(ret, &model.SampleStream{
			Metric: m,
			Values: values,
		})
	}
	return ret, nil
}

func toVector(res *promql.Result) (model.Vector, error) {
	vec, err := res.Vector()
	if err != nil {
		return nil, err
	}
	ret := make([]*model.Sample, 0, len(vec))
	for _, s := range vec {
		m := model.Metric{}
		for _, label := range s.Metric {
			m[model.LabelName(label.Name)] = model.LabelValue(label.Value)
		}
		ret = append(ret, &model.Sample{
			Metric: m,
			Value:  model.SampleValue(s.F),
		})
	}
	return ret, nil
}

// UNIMPLEMENTED BELOW

func (p *PrometheusV1Server) Rules(context.Context) (v1.RulesResult, error) {
	return v1.RulesResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) Alerts(context.Context) (v1.AlertsResult, error) {
	return v1.AlertsResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) AlertManagers(context.Context) (v1.AlertManagersResult, error) {
	return v1.AlertManagersResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) CleanTombstones(context.Context) error {
	return errors.New("unimplemented")
}

func (p *PrometheusV1Server) Config(context.Context) (v1.ConfigResult, error) {
	return v1.ConfigResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) DeleteSeries(context.Context, []string, time.Time, time.Time) error {
	return errors.New("unimplemented")
}

func (p *PrometheusV1Server) Flags(context.Context) (v1.FlagsResult, error) {
	return v1.FlagsResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) QueryExemplars(context.Context, string, time.Time, time.Time) ([]v1.ExemplarQueryResult, error) {
	return nil, errors.New("unimplemented")
}

func (p *PrometheusV1Server) Buildinfo(context.Context) (v1.BuildinfoResult, error) {
	return v1.BuildinfoResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) Runtimeinfo(context.Context) (v1.RuntimeinfoResult, error) {
	return v1.RuntimeinfoResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) Snapshot(context.Context, bool) (v1.SnapshotResult, error) {
	return v1.SnapshotResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) Targets(context.Context) (v1.TargetsResult, error) {
	return v1.TargetsResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) TargetsMetadata(context.Context, string, string, string) ([]v1.MetricMetadata, error) {
	return nil, errors.New("unimplemented")
}

func (p *PrometheusV1Server) TSDB(context.Context, ...v1.Option) (v1.TSDBResult, error) {
	return v1.TSDBResult{}, errors.New("unimplemented")
}

func (p *PrometheusV1Server) WalReplay(context.Context) (v1.WalReplayStatus, error) {
	return v1.WalReplayStatus{}, errors.New("unimplemented")
}

var _ v1.API = (*PrometheusV1Server)(nil)
