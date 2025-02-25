package loopback

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/chunks"
	"github.com/prometheus/prometheus/util/annotations"
)

type series struct {
	metric     metric
	datapoints chunks.SampleSlice
}

type fakePromStorageQuerier struct {
	allMetrics []metric
}

func newFakePromStorageQuerier() *fakePromStorageQuerier {
	return &fakePromStorageQuerier{allMetrics: metrics()}
}

type staticSeriesSet struct {
	sync.RWMutex

	idx           int
	matchedSeries []series
}

func newStaticSeriesSet(matchedSeries []series) *staticSeriesSet {
	return &staticSeriesSet{idx: -1, matchedSeries: matchedSeries}
}

func (s *staticSeriesSet) Next() bool {
	s.idx++
	r := s.idx < len(s.matchedSeries)
	return r
}

func (s *staticSeriesSet) At() storage.Series {
	if s.idx < len(s.matchedSeries) {
		m := s.matchedSeries[s.idx]
		return storage.NewListSeries(m.metric.labels, m.datapoints)
	}
	return nil
}

func (s *staticSeriesSet) Err() error {
	return nil
}

func (s *staticSeriesSet) Warnings() annotations.Annotations {
	return annotations.Annotations{}
}

var _ storage.SeriesSet = (*staticSeriesSet)(nil)

func matchesAll(labels labels.Labels, matchers []*labels.Matcher) bool {
	labelsByName := labels.Map()
	for _, matcher := range matchers {
		labelValue, ok := labelsByName[matcher.Name]
		if !ok {
			return false
		}
		if !matcher.Matches(labelValue) {
			return false
		}
	}
	return true
}
func (f *fakePromStorageQuerier) Select(_ context.Context, sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	var matchedSeries []series
	for _, m := range f.allMetrics {
		if matchesAll(m.labels, matchers) {
			samples := generateSamples(m, time.UnixMilli(hints.Start), time.UnixMilli(hints.End), time.Duration(hints.Step)*time.Millisecond)
			if !m.isGauge {
				cumsumInplace(samples)
			}
			matchedSeries = append(matchedSeries, series{
				metric:     m,
				datapoints: samples,
			})
		}
	}

	return newStaticSeriesSet(matchedSeries)
}

type floatSample struct {
	timestamp time.Time
	value     float64
}

func (f floatSample) T() int64 {
	return f.timestamp.UnixMilli()
}

func (f floatSample) F() float64 {
	return f.value
}

func (f floatSample) H() *histogram.Histogram {
	return nil
}

func (f floatSample) FH() *histogram.FloatHistogram {
	return nil
}

func (f floatSample) Type() chunkenc.ValueType {
	return chunkenc.ValFloat
}

func (f floatSample) Copy() chunks.Sample {
	return f
}

var _ chunks.Sample = (*floatSample)(nil)

func generateSamples(m metric, start, stop time.Time, step time.Duration) chunks.SampleSlice {
	var ret chunks.SampleSlice
	cur := start
	for !cur.After(stop) {
		ret = append(ret, &floatSample{timestamp: cur, value: m.sampleParams.value(cur)})
		cur = cur.Add(step)
	}
	return ret
}

func (f *fakePromStorageQuerier) LabelValues(ctx context.Context, name string, hints *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, annotations.Annotations{}, errors.New("unimplemented")
}

func (f *fakePromStorageQuerier) LabelNames(ctx context.Context, hints *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, annotations.Annotations{}, errors.New("unimplemented")
}

func (f *fakePromStorageQuerier) Close() error {
	return nil
}

var _ storage.Querier = (*fakePromStorageQuerier)(nil)

type fakePromStorageQueryable struct {
	querier *fakePromStorageQuerier
}

func newFakePromStorageQueryable(querier *fakePromStorageQuerier) *fakePromStorageQueryable {
	return &fakePromStorageQueryable{querier: querier}
}

func (f *fakePromStorageQueryable) Querier(int64, int64) (storage.Querier, error) {
	return f.querier, nil
}

var _ storage.Queryable = (*fakePromStorageQueryable)(nil)
