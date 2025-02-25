package loopback

import (
	"iter"
	"math/rand"

	"github.com/prometheus/prometheus/tsdb/chunks"
)

func dedupe[T comparable](input []T) []T {
	seen := map[T]struct{}{}
	for _, inp := range input {
		seen[inp] = struct{}{}
	}
	out := make([]T, 0, len(seen))
	for val, _ := range seen {
		out = append(out, val)
	}
	return out
}

func choice[T any](rnd *rand.Rand, options []T) T {
	return options[rnd.Intn(len(options))]
}

func product[T any](iterables ...[]T) iter.Seq[[]T] {
	n := 1
	for _, iterable := range iterables {
		n *= len(iterable)
	}
	return func(yield func([]T) bool) {
		for i := range n {
			vals := make([]T, 0, len(iterables))
			for _, iterable := range iterables {
				iterableIdx := i % len(iterable)
				i /= len(iterable)
				vals = append(vals, iterable[iterableIdx])
			}
			if !yield(vals) {
				return
			}
		}
	}
}

func cumsumInplace(input chunks.SampleSlice) {
	prev := 0.0
	for _, sample := range input {
		fs, ok := sample.(*floatSample)
		if !ok {
			panic("could not convert to floatSample")
		}
		if prev < 0 {
			continue
		}
		fs.value += prev
		prev = fs.value
	}
}
