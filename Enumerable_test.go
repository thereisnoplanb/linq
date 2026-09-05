package linq_test

import (
	"cmp"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/thereisnoplanb/linq"
)

func TestEnumerable_Aggregate(t *testing.T) {
	// Example test for Aggregate
	source := linq.FromSlice([]int{1, 2, 3, 4})
	result := source.Aggregate(0, func(accumulator, item int) int {
		return accumulator + item
	}, func(accumulator int, count int) int {
		return accumulator * count
	})
	want := 40
	if result != want {
		t.Fatalf("Aggregate() = %v, want %v", result, want)
	}
}

func TestEnumerable_All(t *testing.T) {
	tests := []struct {
		source []int
		want   bool
	}{
		{source: []int{2, 4, 6, 8}, want: true},
		{source: []int{1, 2, 3, 4}, want: false},
		{source: []int{1, 3, 5, 7}, want: false},
		{source: []int{}, want: true},
	}

	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.All(func(item int) bool {
			return item%2 == 0
		})
		if result != test.want {
			t.Fatalf("All() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Any(t *testing.T) {
	tests := []struct {
		source []int
		want   bool
	}{
		{source: []int{1, 3, 5, 7}, want: false},
		{source: []int{2, 3, 5, 7}, want: true},
		{source: []int{2, 4, 6, 8}, want: true},
		{source: []int{}, want: false},
	}

	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.Any(func(item int) bool {
			return item%2 == 0
		})
		if result != test.want {
			t.Fatalf("Any() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Append(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3})
	result := source.Append(4, 5, 6).ToSlice()
	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Append() = %v, want %v", result, want)
	}
}

func TestEnumerable_Average(t *testing.T) {
	tests := []struct {
		source  []int
		want    float64
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4},
			want:    2.5,
			wantErr: false,
		},
		{
			source:  []int{},
			want:    math.NaN(),
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.Average[int]()
		want := test.want
		if (err != nil) != test.wantErr {
			t.Fatalf("Average() error = %v, wantErr %v", err, test.wantErr)
			if math.IsNaN(result) && math.IsNaN(want) {
				continue
			}
			if result != want {
				t.Fatalf("Average() = %v, want %v", result, want)
			}
			continue
		}
		if math.IsNaN(result) && math.IsNaN(want) {
			continue
		}
		if result != want {
			t.Fatalf("Average() = %v, want %v", result, want)
		}
	}
}

func TestEnumerable_AverageFunc(t *testing.T) {
	tests := []struct {
		source  []int
		value   func(item int) float64
		want    float64
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4},
			value:   func(item int) float64 { return float64(item) },
			want:    2.5,
			wantErr: false,
		},
		{
			source:  []int{},
			value:   func(item int) float64 { return float64(item) },
			want:    math.NaN(),
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.AverageFunc(test.value)
		want := test.want
		if (err != nil) != test.wantErr {
			t.Fatalf("AverageFunc() error = %v, wantErr %v", err, test.wantErr)
			if math.IsNaN(result) && math.IsNaN(want) {
				continue
			}
			if result != want {
				t.Fatalf("AverageFunc() = %v, want %v", result, want)
			}
			continue
		}
		if math.IsNaN(result) && math.IsNaN(want) {
			continue
		}
		if result != want {
			t.Fatalf("AverageFunc() = %v, want %v", result, want)
		}
	}
}

func TestEnumerable_AverageComplex(t *testing.T) {
	tests := []struct {
		source  []complex128
		want    complex128
		wantErr bool
	}{
		{
			source:  []complex128{1 + 2i, 3 + 4i, 5 + 6i},
			want:    3 + 4i,
			wantErr: false,
		},
		{
			source:  []complex128{},
			want:    complex(math.NaN(), math.NaN()),
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.AverageComplex[complex128]()
		want := test.want
		if (err != nil) != test.wantErr {
			t.Fatalf("AverageComplex() error = %v, wantErr %v", err, test.wantErr)
			if result == want || (math.IsNaN(real(result)) && math.IsNaN(real(want)) && math.IsNaN(imag(result)) && math.IsNaN(imag(want))) {
				continue
			}
			if result != want {
				t.Fatalf("AverageComplex() = %v, want %v", result, want)
			}
			continue
		}
		if result == want || (math.IsNaN(real(result)) && math.IsNaN(real(want)) && math.IsNaN(imag(result)) && math.IsNaN(imag(want))) {
			continue
		}
		if result != want {
			t.Fatalf("AverageComplex() = %v, want %v", result, want)
		}
	}
}

func TestEnumerable_AverageComplexFunc(t *testing.T) {
	tests := []struct {
		source  []complex128
		value   func(item complex128) complex128
		want    complex128
		wantErr bool
	}{
		{
			source:  []complex128{1 + 2i, 3 + 4i, 5 + 6i},
			value:   func(item complex128) complex128 { return item },
			want:    3 + 4i,
			wantErr: false,
		},
		{
			source:  []complex128{},
			value:   func(item complex128) complex128 { return item },
			want:    complex(math.NaN(), math.NaN()),
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.AverageComplexFunc(test.value)
		want := test.want
		if (err != nil) != test.wantErr {
			t.Fatalf("AverageComplexFunc() error = %v, wantErr %v", err, test.wantErr)
			if result == want || (math.IsNaN(real(result)) && math.IsNaN(real(want)) && math.IsNaN(imag(result)) && math.IsNaN(imag(want))) {
				continue
			}
			if result != want {
				t.Fatalf("AverageComplexFunc() = %v, want %v", result, want)
			}
			continue
		}
		if result == want || (math.IsNaN(real(result)) && math.IsNaN(real(want)) && math.IsNaN(imag(result)) && math.IsNaN(imag(want))) {
			continue
		}
		if result != want {
			t.Fatalf("AverageComplexFunc() = %v, want %v", result, want)
		}
	}
}

func TestEnumerable_Cast(t *testing.T) {
	source := linq.FromSlice([]any{1, 2, 3})
	result := source.Cast[int]().ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Cast() = %v, want %v", result, want)
	}
}

func TestEnumerable_Chunk(t *testing.T) {
	tests := []struct {
		source []int
		size   int
		want   [][]int
	}{
		{source: []int{1, 2, 3, 4, 5}, size: 1, want: [][]int{{1}, {2}, {3}, {4}, {5}}},
		{source: []int{1, 2, 3, 4, 5}, size: 2, want: [][]int{{1, 2}, {3, 4}, {5}}},
		{source: []int{1, 2, 3, 4, 5}, size: 5, want: [][]int{{1, 2, 3, 4, 5}}},
		{source: []int{1, 2, 3, 4, 5}, size: 10, want: [][]int{{1, 2, 3, 4, 5}}},
		{source: []int{}, size: 3, want: [][]int{}},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.Chunk(test.size).ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("Chunk() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Concat(t *testing.T) {
	source1 := linq.FromSlice([]int{1, 2, 3})
	source2 := linq.FromSlice([]int{4, 5, 6})
	result := source1.Concat(source2).ToSlice()
	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Concat() = %v, want %v", result, want)
	}
}

func TestEnumerable_Contains(t *testing.T) {
	tests := []struct {
		source []int
		value  int
		want   bool
	}{
		{source: []int{1, 2, 3, 4, 5}, value: 3, want: true},
		{source: []int{1, 2, 3, 4, 5}, value: 6, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.Contains[int](test.value)
		if result != test.want {
			t.Fatalf("Contains() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		value  time.Time
		want   bool
	}{
		{source: []time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
		}, value: time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC), want: true},
		{source: []time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
		}, value: time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC), want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsEquatable[time.Time](test.value)
		if result != test.want {
			t.Fatalf("ContainsEquatable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsFunc(t *testing.T) {
	tests := []struct {
		source []int
		value  int
		want   bool
	}{
		{source: []int{1, 2, 3, 4, 5}, value: 3, want: true},
		{source: []int{1, 2, 3, 4, 5}, value: 6, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsFunc(test.value, func(a, b int) bool { return a == b })
		if result != test.want {
			t.Fatalf("ContainsFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAll(t *testing.T) {
	tests := []struct {
		source []int
		values []int
		want   bool
	}{
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 4}, want: true},
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 7}, want: false},
		{source: []int{1, 2, 3, 4, 5}, values: []int{6, 7}, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAll[int](test.values...)
		if result != test.want {
			t.Fatalf("ContainsAll() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAllEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		values []time.Time
		want   bool
	}{
		{source: []time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
		}, values: []time.Time{
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
		}, want: true},
		{source: []time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
		}, values: []time.Time{
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
		}, want: false},
		{source: []time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
		}, values: []time.Time{
			time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
		}, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAllEquatable[time.Time](test.values...)
		if result != test.want {
			t.Fatalf("ContainsAllEquatable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAllFunc(t *testing.T) {
	tests := []struct {
		source []int
		values []int
		want   bool
	}{
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 4}, want: true},
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 6}, want: false},
		{source: []int{1, 2, 3, 4, 5}, values: []int{7, 8}, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAllFunc(func(a, b int) bool { return a == b }, test.values...)
		if result != test.want {
			t.Fatalf("ContainsAllFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAny(t *testing.T) {
	tests := []struct {
		source []int
		values []int
		want   bool
	}{
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 4}, want: true},
		{source: []int{1, 2, 3, 4, 5}, values: []int{2, 3, 6}, want: true},
		{source: []int{1, 2, 3, 4, 5}, values: []int{7, 8}, want: false},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAny[int](test.values...)
		if result != test.want {
			t.Fatalf("ContainsAny() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAnyEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		values []time.Time
		want   bool
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			values: []time.Time{
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			values: []time.Time{
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			values: []time.Time{
				time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAnyEquatable[time.Time](test.values...)
		if result != test.want {
			t.Fatalf("ContainsAnyEquatable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ContainsAnyFunc(t *testing.T) {
	tests := []struct {
		source []int
		values []int
		want   bool
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			values: []int{2, 3, 4},
			want:   true,
		},
		{
			source: []int{1, 2, 3, 4, 5},
			values: []int{5, 6, 7},
			want:   true,
		},
		{
			source: []int{1, 2, 3, 4, 5},
			values: []int{6, 7, 8},
			want:   false,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ContainsAnyFunc(func(a, b int) bool { return a == b }, test.values...)
		if result != test.want {
			t.Fatalf("ContainsAnyFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Count(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			want:      3,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item <= 2 },
			want:      2,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      5,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.Count(test.predicate)
		if result != test.want {
			t.Fatalf("Count() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Distinct(t *testing.T) {
	tests := []struct {
		source []int
		want   []int
	}{
		{
			source: []int{1, 2, 2, 3, 3, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{},
			want:   []int{},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.Distinct[int]().ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("Distinct() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_DistinctEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		want   []time.Time
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{},
			want:   []time.Time{},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.DistinctEquatable[time.Time]().ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("DistinctEquatable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_DistinctFunc(t *testing.T) {
	tests := []struct {
		source []int
		want   []int
	}{
		{
			source: []int{1, 2, 2, 3, 3, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{},
			want:   []int{},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.DistinctFunc(func(a, b int) bool { return a == b }).ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("DistinctFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ElementAt(t *testing.T) {
	tests := []struct {
		source  []int
		index   int
		want    int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			index:   2,
			want:    3,
			wantErr: false,
		},
		{
			source:  []int{1, 2, 3, 4, 5},
			index:   10,
			want:    0,
			wantErr: true,
		},
		{
			source:  []int{1, 2, 3, 4, 5},
			index:   -1,
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.ElementAt(test.index)
		if (err != nil) != test.wantErr {
			t.Fatalf("ElementAt() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("ElementAt() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ElementAtOrDefault(t *testing.T) {
	tests := []struct {
		source []int
		index  int
		want   int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			index:  2,
			want:   3,
		},
		{
			source: []int{1, 2, 3, 4, 5},
			index:  10,
			want:   0,
		},
		{
			source: []int{1, 2, 3, 4, 5},
			index:  -1,
			want:   0,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ElementAtOrDefault(test.index)
		if result != test.want {
			t.Fatalf("ElementAtOrDefault() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ElementAtOrFallback(t *testing.T) {
	tests := []struct {
		source   []int
		index    int
		fallback int
		want     int
	}{
		{
			source:   []int{1, 2, 3, 4, 5},
			index:    2,
			fallback: 42,
			want:     3,
		},
		{
			source:   []int{1, 2, 3, 4, 5},
			index:    10,
			fallback: 42,
			want:     42,
		},
		{
			source:   []int{1, 2, 3, 4, 5},
			index:    -1,
			fallback: 42,
			want:     42,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ElementAtOrFallback(test.index, test.fallback)
		if result != test.want {
			t.Fatalf("ElementAtOrFallback() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ElementAtOrNil(t *testing.T) {
	tests := []struct {
		source []int
		index  int
		want   *int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			index:  2,
			want:   ptr(3),
		},
		{
			source: []int{1, 2, 3, 4, 5},
			index:  10,
			want:   nil,
		},
		{
			source: []int{1, 2, 3, 4, 5},
			index:  -1,
			want:   nil,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.ElementAtOrNil(test.index)
		if result != nil && test.want != nil && *result != *test.want {
			t.Fatalf("ElementAtOrNil() = %v, want %v", *result, *test.want)
		} else if result == nil && test.want == nil {
			// both are nil, test passes
		} else if result == nil && test.want != nil {
			t.Fatalf("ElementAtOrNil() = %v, want %v", result, *test.want)
		} else if result != nil && test.want == nil {
			t.Fatalf("ElementAtOrNil() = %v, want %v", *result, test.want)
		}
	}
}

func TestEnumerable_Except(t *testing.T) {
	tests := []struct {
		source []int
		except []int
		want   []int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{2, 4},
			want:   []int{1, 3, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{1, 2, 3, 4, 5},
			want:   []int{},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{6, 7},
			want:   []int{1, 2, 3, 4, 5},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		except := linq.FromSlice(test.except)
		result := source.Except[int](except).ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("Except() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ExceptEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		except []time.Time
		want   []time.Time
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			except: []time.Time{
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			except: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{},
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			except: []time.Time{},
			want: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			except: []time.Time{
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		except := linq.FromSlice(test.except)
		result := source.ExceptEquatable[time.Time](except).ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("ExceptEquatable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_ExceptFunc(t *testing.T) {
	tests := []struct {
		source []int
		except []int
		want   []int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{2, 4},
			want:   []int{1, 3, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{1, 2, 3, 4, 5},
			want:   []int{},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			source: []int{1, 2, 3, 4, 5},
			except: []int{6, 7},
			want:   []int{1, 2, 3, 4, 5},
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		except := linq.FromSlice(test.except)
		result := source.ExceptFunc(except, func(a, b int) bool { return a == b }).ToSlice()
		if !reflect.DeepEqual(result, test.want) {
			t.Fatalf("ExceptFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_First(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      int
		wantErr   bool
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			want:      3,
			wantErr:   false,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{},
			predicate: func(item int) bool { return item > 0 },
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      1,
			wantErr:   false,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.First(test.predicate)
		if (err != nil) != test.wantErr {
			t.Fatalf("First() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("First() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_FirstOrDefault(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			want:      3,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      0,
		},
		{
			source:    []int{},
			predicate: func(item int) bool { return item > 0 },
			want:      0,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      1,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.FirstOrDefault(test.predicate)
		if result != test.want {
			t.Fatalf("FirstOrDefault() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_FirstOrFallback(t *testing.T) {
	tests := []struct {
		source    []int
		fallback  int
		predicate func(item int) bool
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			fallback:  42,
			want:      3,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{},
			predicate: func(item int) bool { return item > 0 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			fallback:  42,
			want:      1,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.FirstOrFallback(test.fallback, test.predicate)
		if result != test.want {
			t.Fatalf("FirstOrFallback() = %v, want %v", result, test.want)
		}
	}
}

//go:fix inline
func ptr[T any](v T) *T {
	return new(v)
}

func TestEnumerable_FirstOrNil(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      *int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			want:      ptr(3),
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      nil,
		},
		{
			source:    []int{},
			predicate: func(item int) bool { return item > 0 },
			want:      nil,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      ptr(1),
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.FirstOrNil(test.predicate)
		if result != nil && test.want != nil && *result != *test.want {
			t.Fatalf("FirstOrNil() = %v, want %v", *result, *test.want)
		} else if result == nil && test.want != nil {
			t.Fatalf("FirstOrNil() = %v, want %v", result, *test.want)
		} else if result != nil && test.want == nil {
			t.Fatalf("FirstOrNil() = %v, want %v", *result, test.want)
		}
	}
}

func TestEnumerable_ForEach(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	sum := 0
	source.ForEach(func(item int) {
		sum += item
	})
	want := 15
	if sum != want {
		t.Fatalf("ForEach() = %v, want %v", sum, want)
	}
}

func TestEnumerable_GroupBy(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "apricot"})
	result := source.GroupBy(func(item string) rune { return rune(item[0]) }).
		ToMap(
			func(group linq.Grouping[rune, string]) rune { return group.Key },
			func(group linq.Grouping[rune, string]) []string {
				return group.ToSlice()
			})
	want := map[rune][]string{
		'a': {"apple", "apricot"},
		'b': {"banana"},
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("GroupBy() = %v, want %v", result, want)
	}
}

func TestEnumerable_GroupJoin(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})
	result := outer.GroupJoin(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(outerItem string, matches linq.Enumerable[string]) int {
			return len(matches.ToSlice())
		}).ToSlice()

	want := []int{2, 1, 0}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("GroupJoin() = %v, want %v", result, want)
	}
}

func TestEnumerable_GroupJoinEquatable(t *testing.T) {
	outer := linq.FromSlice([]time.Time{
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	inner := linq.FromSlice([]time.Time{
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := outer.GroupJoinEquatable(inner,
		func(item time.Time) time.Time { return item },
		func(item time.Time) time.Time { return item },
		func(outerItem time.Time, matches linq.Enumerable[time.Time]) int {
			return len(matches.ToSlice())
		}).ToSlice()

	want := []int{2, 1, 0}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("GroupJoinEquatable() = %v, want %v", result, want)
	}
}

func TestEnumerable_GroupJoinFunc(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})
	result := outer.GroupJoinFunc(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(outerItem string, matches linq.Enumerable[string]) int {
			return len(matches.ToSlice())
		},
		func(x, y string) bool { return x == y },
	).ToSlice()

	want := []int{2, 1, 0}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("GroupJoinFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Intersect(t *testing.T) {
	tests := []struct {
		source []int
		other  []int
		want   []int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			other:  []int{4, 5, 6, 7, 8},
			want:   []int{4, 5},
		},
		{
			source: []int{1, 2, 3},
			other:  []int{3, 4, 5},
			want:   []int{3},
		},
		{
			source: []int{1, 2, 3},
			other:  []int{6, 7, 8},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		source := linq.FromSlice(tt.source)
		other := linq.FromSlice(tt.other)
		result := source.Intersect[int](other).ToSlice()
		if !reflect.DeepEqual(result, tt.want) {
			t.Fatalf("Intersect() = %v, want %v", result, tt.want)
		}
	}
}

func TestEnumerable_IntersectEquatable(t *testing.T) {
	tests := []struct {
		source []time.Time
		other  []time.Time
		want   []time.Time
	}{
		{
			source: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			other: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			other: []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			source: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			other: []time.Time{
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{},
		},
	}
	for _, tt := range tests {
		source := linq.FromSlice(tt.source)
		other := linq.FromSlice(tt.other)
		result := source.IntersectEquatable[time.Time](other).ToSlice()
		if !reflect.DeepEqual(result, tt.want) {
			t.Fatalf("IntersectEquatable() = %v, want %v", result, tt.want)
		}
	}
}

func TestEnumerable_IntersectFunc(t *testing.T) {
	tests := []struct {
		source []int
		other  []int
		want   []int
	}{
		{
			source: []int{1, 2, 3, 4, 5},
			other:  []int{4, 5, 6, 7, 8},
			want:   []int{4, 5},
		},
		{
			source: []int{1, 2, 3},
			other:  []int{3, 4, 5},
			want:   []int{3},
		},
		{
			source: []int{1, 2, 3},
			other:  []int{4, 5, 6},
			want:   []int{},
		},
	}
	for _, tt := range tests {
		source := linq.FromSlice(tt.source)
		other := linq.FromSlice(tt.other)
		result := source.IntersectFunc(other, func(x, y int) bool { return x == y }).ToSlice()
		if !reflect.DeepEqual(result, tt.want) {
			t.Fatalf("IntersectFunc() = %v, want %v", result, tt.want)
		}
	}
}

func TestEnumerable_Join(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})
	result := outer.Join(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(outerItem string, innerItem string) string {
			return outerItem + innerItem
		}).ToSlice()

	want := []string{"AA", "AA", "BB"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Join() = %v, want %v", result, want)
	}
}

func TestEnumerable_JoinEquatable(t *testing.T) {
	outer := linq.FromSlice([]time.Time{
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	inner := linq.FromSlice([]time.Time{
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := outer.JoinEquatable(inner,
		func(item time.Time) time.Time { return item },
		func(item time.Time) time.Time { return item },
		func(outerItem time.Time, innerItem time.Time) string {
			return outerItem.Format("2006") + innerItem.Format("2006")
		}).ToSlice()

	want := []string{"20202020", "20202020", "20212021"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Join() = %v, want %v", result, want)
	}
}

func TestEnumerable_JoinFunc(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})
	result := outer.JoinFunc(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(x, y string) bool { return x == y },
		func(outerItem string, innerItem string) string {
			return outerItem + innerItem
		}).ToSlice()

	want := []string{"AA", "AA", "BB"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("JoinFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Last(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      int
		wantErr   bool
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 2 },
			want:      5,
			wantErr:   false,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{},
			predicate: nil,
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      5,
			wantErr:   false,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.Last(test.predicate)
		if (err != nil) != test.wantErr {
			t.Fatalf("Last() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("Last() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_LastOrDefault(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      0,
		},
		{
			source:    []int{},
			predicate: nil,
			want:      0,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      5,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.LastOrDefault(test.predicate)
		if result != test.want {
			t.Fatalf("LastOrDefault() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_LastOrFallback(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		fallback  int
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{},
			predicate: func(item int) bool { return item > 0 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			fallback:  42,
			want:      5,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.LastOrFallback(test.fallback, test.predicate)
		if result != test.want {
			t.Fatalf("LastOrFallback() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_LastOrNil(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(item int) bool
		want      *int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 10 },
			want:      nil,
		},
		{
			source:    []int{},
			predicate: nil,
			want:      nil,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: nil,
			want:      ptr(5),
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.LastOrNil(test.predicate)
		if result != nil && test.want != nil && *result != *test.want {
			t.Fatalf("LastOrNil() = %v, want %v", *result, *test.want)
		} else if result == nil && test.want != nil {
			t.Fatalf("LastOrNil() = %v, want %v", result, *test.want)
		} else if result != nil && test.want == nil {
			t.Fatalf("LastOrNil() = %v, want %v", *result, test.want)
		}
	}
}

func TestEnumerable_LeftJoin(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})

	result := outer.LeftJoin(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(outerItem string, innerItem *string) string {
			return outerItem + ValueOrDefault(innerItem)
		}).ToSlice()

	want := []string{"AA", "AA", "BB", "C"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("LeftJoin() = %v, want %v", result, want)
	}
}

func TestEnumerable_LeftJoinEquatable(t *testing.T) {
	outer := linq.FromSlice([]time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	inner := linq.FromSlice([]time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	result := outer.LeftJoinEquatable(inner,
		func(item time.Time) time.Time { return item },
		func(item time.Time) time.Time { return item },
		func(outerItem time.Time, innerItem *time.Time) string {
			return outerItem.Format("2006") + NullValueOrDefault(innerItem, func(value *time.Time) string { return value.Format("2006") })
		}).ToSlice()

	want := []string{"20002000", "20002000", "20012001", "2002"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("LeftJoinEquatable() = %v, want %v", result, want)
	}
}

func TestEnumerable_LeftJoinFunc(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	inner := linq.FromSlice([]string{"A", "A", "B"})

	result := outer.LeftJoinFunc(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(a, b string) bool { return a == b },
		func(outerItem string, innerItem *string) string {
			return outerItem + ValueOrDefault(innerItem)
		}).ToSlice()

	want := []string{"AA", "AA", "BB", "C"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("LeftJoinFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Max(t *testing.T) {
	tests := []struct {
		source  []int
		want    int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			want:    5,
			wantErr: false,
		},
		{
			source:  []int{},
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.Max[int]()
		if (err != nil) != test.wantErr {
			t.Fatalf("Max() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("Max() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MaxComparable(t *testing.T) {
	tests := []struct {
		source  []time.Time
		want    time.Time
		wantErr bool
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			source:  []time.Time{},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MaxComparable[time.Time]()
		if (err != nil) != test.wantErr {
			t.Fatalf("MaxComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MaxComparable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MaxFunc(t *testing.T) {
	tests := []struct {
		source  []int
		want    int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			want:    5,
			wantErr: false,
		},
		{
			source:  []int{},
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MaxFunc(cmp.Compare[int])
		if (err != nil) != test.wantErr {
			t.Fatalf("MaxFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MaxFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MaxBy(t *testing.T) {
	tests := []struct {
		source  []string
		want    string
		wantErr bool
	}{
		{
			source:  []string{"apple", "banana", "cherry"},
			want:    "banana",
			wantErr: false,
		},
		{
			source:  []string{},
			want:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MaxBy(func(item string) int {
			return len(item)
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MaxBy() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MaxBy() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MaxByComparable(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		want    Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want:    Person{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			want:    Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MaxByComparable(func(item Person) time.Time {
			return item.BirthDate
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MaxByComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MaxByComparable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MaxByFunc(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		want    Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want:    Person{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			want:    Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MaxByFunc(
			func(item Person) time.Time {
				return item.BirthDate
			},
			func(a, b time.Time) int {
				return a.Compare(b)
			},
		)
		if (err != nil) != test.wantErr {
			t.Fatalf("MaxByFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MaxByFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_Min(t *testing.T) {
	tests := []struct {
		source  []int
		want    int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			want:    1,
			wantErr: false,
		},
		{
			source:  []int{},
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.Min[int]()
		if (err != nil) != test.wantErr {
			t.Fatalf("Min() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("Min() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinComparable(t *testing.T) {
	tests := []struct {
		source  []time.Time
		want    time.Time
		wantErr bool
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			source:  []time.Time{},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MinComparable[time.Time]()
		if (err != nil) != test.wantErr {
			t.Fatalf("MinComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MinComparable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinFunc(t *testing.T) {
	tests := []struct {
		source  []int
		want    int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			want:    1,
			wantErr: false,
		},
		{
			source:  []int{},
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MinFunc(cmp.Compare[int])
		if (err != nil) != test.wantErr {
			t.Fatalf("MinFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MinFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinBy(t *testing.T) {
	tests := []struct {
		source  []string
		want    string
		wantErr bool
	}{
		{
			source:  []string{"apple", "banana", "cherry"},
			want:    "apple",
			wantErr: false,
		},
		{
			source:  []string{},
			want:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MinBy(func(item string) int {
			return len(item)
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MinBy() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MinBy() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinByComparable(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		want    Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want:    Person{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			want:    Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MinByComparable(func(item Person) time.Time {
			return item.BirthDate
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MinByComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MinByComparable() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinByFunc(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		want    Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want:    Person{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			want:    Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.MinByFunc(
			func(item Person) time.Time {
				return item.BirthDate
			},
			func(a, b time.Time) int {
				return a.Compare(b)
			},
		)
		if (err != nil) != test.wantErr {
			t.Fatalf("MinByFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if result != test.want {
			t.Fatalf("MinByFunc() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_MinMax(t *testing.T) {
	tests := []struct {
		source  []int
		wantMin int
		wantMax int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			wantMin: 1,
			wantMax: 5,
			wantErr: false,
		},
		{
			source:  []int{},
			wantMin: 0,
			wantMax: 0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMax[int]()
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMax() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMax() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_MinMaxComparable(t *testing.T) {
	tests := []struct {
		source  []time.Time
		wantMin time.Time
		wantMax time.Time
		wantErr bool
	}{
		{
			source: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantMin: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			wantMax: time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMaxComparable[time.Time]()
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMaxComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMaxComparable() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_MinMaxFunc(t *testing.T) {
	tests := []struct {
		source  []int
		wantMin int
		wantMax int
		wantErr bool
	}{
		{
			source:  []int{1, 2, 3, 4, 5},
			wantMin: 1,
			wantMax: 5,
			wantErr: false,
		},
		{
			source:  []int{},
			wantMin: 0,
			wantMax: 0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMaxFunc(cmp.Compare[int])
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMaxFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMaxFunc() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_MinMaxBy(t *testing.T) {
	tests := []struct {
		source  []string
		wantMin string
		wantMax string
		wantErr bool
	}{
		{
			source:  []string{"apple", "banana", "cherry"},
			wantMin: "apple",
			wantMax: "banana",
			wantErr: false,
		},
		{
			source:  []string{},
			wantMin: "",
			wantMax: "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMaxBy(func(item string) int {
			return len(item)
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMaxBy() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMaxBy() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_MinMaxByComparable(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		wantMin Person
		wantMax Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantMin: Person{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantMax: Person{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			wantMin: Person{},
			wantMax: Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMaxByComparable(func(item Person) time.Time {
			return item.BirthDate
		})
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMaxByComparable() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMaxByComparable() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_MinMaxByFunc(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	tests := []struct {
		source  []Person
		wantMin Person
		wantMax Person
		wantErr bool
	}{
		{
			source: []Person{
				{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Bob", time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)},
				{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantMin: Person{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantMax: Person{"Charlie", time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			source:  []Person{},
			wantMin: Person{},
			wantMax: Person{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		resultMin, resultMax, err := source.MinMaxByFunc(
			func(item Person) time.Time {
				return item.BirthDate
			},
			func(a, b time.Time) int {
				return a.Compare(b)
			},
		)
		if (err != nil) != test.wantErr {
			t.Fatalf("MinMaxByFunc() error = %v, wantErr %v", err, test.wantErr)
		}
		if resultMin != test.wantMin || resultMax != test.wantMax {
			t.Fatalf("MinMaxByFunc() = (%v, %v), want (%v, %v)", resultMin, resultMax, test.wantMin, test.wantMax)
		}
	}
}

func TestEnumerable_OfType(t *testing.T) {
	source := linq.FromSlice([]any{
		1,
		1.0,
		complex(1, 0),
		"string",
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := source.OfType[time.Time]().ToSlice()
	want := []time.Time{time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OfType() = %v, want %v", result, want)
	}
}

func TestEnumerable_Order(t *testing.T) {
	source := linq.FromSlice([]int{5, 3, 1, 4, 2})
	result := source.Order[int]().ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Order() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderComparable(t *testing.T) {
	source := linq.FromSlice([]time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := source.OrderComparable[time.Time]().ToSlice()
	want := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderComparable() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderFunc(t *testing.T) {
	source := linq.FromSlice([]int{5, 3, 1, 4, 2})
	result := source.OrderFunc(cmp.Compare[int]).ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderBy(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "cherry"})
	result := source.OrderBy(func(item string) int {
		return len(item)
	}).ToSlice()
	want := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderBy() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderByComparable(t *testing.T) {
	source := linq.FromSlice([]time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := source.OrderByComparable(func(item time.Time) time.Time {
		return item
	}).ToSlice()
	want := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderByComparable() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderByFunc(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "cherry"})
	result := source.OrderByFunc(
		func(item string) int {
			return len(item)
		},
		cmp.Compare[int],
	).ToSlice()
	want := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderByFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescending(t *testing.T) {
	source := linq.FromSlice([]int{5, 3, 1, 4, 2})
	result := source.OrderDescending[int]().ToSlice()
	want := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescending() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescendingComparable(t *testing.T) {
	source := linq.FromSlice([]time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := source.OrderDescendingComparable[time.Time]().ToSlice()
	want := []time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescendingComparable() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescendingFunc(t *testing.T) {
	source := linq.FromSlice([]int{5, 3, 1, 4, 2})
	result := source.OrderDescendingFunc(cmp.Compare[int]).ToSlice()
	want := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescendingFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescendingBy(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "plum"})
	result := source.OrderDescendingBy(func(item string) int {
		return len(item)
	}).ToSlice()
	want := []string{"banana", "apple", "plum"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescendingBy() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescendingByComparable(t *testing.T) {
	source := linq.FromSlice([]time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := source.OrderDescendingByComparable(func(item time.Time) time.Time {
		return item
	}).ToSlice()
	want := []time.Time{
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescendingByComparable() = %v, want %v", result, want)
	}
}

func TestEnumerable_OrderDescendingByFunc(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "plum"})
	result := source.OrderDescendingByFunc(
		func(item string) int {
			return len(item)
		},
		cmp.Compare[int],
	).ToSlice()
	want := []string{"banana", "apple", "plum"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("OrderDescendingByFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Prepend(t *testing.T) {
	source := linq.FromSlice([]int{2, 3, 4, 5})
	result := source.Prepend(1).ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Prepend() = %v, want %v", result, want)
	}
}

func TestEnumerable_Reverse(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.Reverse().ToSlice()
	want := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Reverse() = %v, want %v", result, want)
	}
}

func TestEnumerable_RightJoin(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "A", "B"})
	inner := linq.FromSlice([]string{"A", "B", "C"})
	result := outer.RightJoin(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(outerItem *string, innerItem string) string {
			return ValueOrDefault(outerItem) + innerItem
		}).ToSlice()

	want := []string{"AA", "AA", "BB", "C"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("RightJoin() = %v, want %v", result, want)
	}
}

func TestEnumerable_RightJoinEquatable(t *testing.T) {
	outer := linq.FromSlice([]time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	inner := linq.FromSlice([]time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := outer.RightJoinEquatable(inner,
		func(item time.Time) time.Time { return item },
		func(item time.Time) time.Time { return item },
		func(outerItem *time.Time, innerItem time.Time) string {
			return NullValueOrDefault(outerItem, func(value *time.Time) string { return value.Format("2006") }) + innerItem.Format("2006")
		}).ToSlice()
	want := []string{"20002000", "20002000", "20012001", "2002"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("RightJoinEquatable() = %v, want %v", result, want)
	}
}

func TestEnumerable_RightJoinFunc(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "A", "B"})
	inner := linq.FromSlice([]string{"A", "B", "C"})
	result := outer.RightJoinFunc(inner,
		func(item string) string { return item },
		func(item string) string { return item },
		func(a, b string) bool { return a == b },
		func(outerItem *string, innerItem string) string {
			return ValueOrDefault(outerItem) + innerItem
		}).ToSlice()
	want := []string{"AA", "AA", "BB", "C"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("RightJoinFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Select(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	result := outer.Select(func(item string) string {
		return item + item
	}).ToSlice()

	want := []string{"AA", "BB", "CC"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Select() = %v, want %v", result, want)
	}
}

func TestEnumerable_SelectMany(t *testing.T) {
	outer := linq.FromSlice([]string{"A", "B", "C"})
	result := outer.SelectMany(func(item string) []string {
		return []string{item + item, item + item + item}
	}).ToSlice()

	want := []string{"AA", "AAA", "BB", "BBB", "CC", "CCC"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("SelectMany() = %v, want %v", result, want)
	}
}

func TestEnumerable_SequenceEqual(t *testing.T) {
	tests := []struct {
		seq1 []int
		seq2 []int
		want bool
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []int{1, 2, 4}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3, 4}, false},
	}
	for _, test := range tests {
		seq1 := linq.FromSlice(test.seq1)
		seq2 := linq.FromSlice(test.seq2)
		if got := seq1.SequenceEqual[int](seq2); got != test.want {
			t.Fatalf("SequenceEqual() = %v, want %v", got, test.want)
		}
	}
}

func TestEnumerable_SequenceEqualEquatable(t *testing.T) {
	tests := []struct {
		seq1 []time.Time
		seq2 []time.Time
		want bool
	}{
		{
			seq1: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			seq2: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			seq1: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			seq2: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			seq1: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			seq2: []time.Time{
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, test := range tests {
		seq1 := linq.FromSlice(test.seq1)
		seq2 := linq.FromSlice(test.seq2)
		if got := seq1.SequenceEqualEquatable[time.Time](seq2); got != test.want {
			t.Fatalf("SequenceEqualEquatable() = %v, want %v", got, test.want)
		}
	}
}

func TestEnumerable_SequenceEqualFunc(t *testing.T) {
	tests := []struct {
		seq1 []int
		seq2 []int
		want bool
	}{
		{
			seq1: []int{1, 2, 3},
			seq2: []int{1, 2, 3},
			want: true,
		},
		{
			seq1: []int{1, 2, 3},
			seq2: []int{1, 2, 4},
			want: false,
		},
		{
			seq1: []int{1, 2, 3},
			seq2: []int{1, 2, 3, 4},
			want: false,
		},
	}
	for _, test := range tests {
		seq1 := linq.FromSlice(test.seq1)
		seq2 := linq.FromSlice(test.seq2)
		if got := seq1.SequenceEqualFunc(seq2, func(a, b int) bool { return a == b }); got != test.want {
			t.Fatalf("SequenceEqualFunc() = %v, want %v", got, test.want)
		}
	}
}

func TestEnumerable_Shuffle(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	result := source.Shuffle().ToSlice()
	if reflect.DeepEqual(result, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}) {
		t.Fatalf("Shuffle() = %v, want %v", result, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	}
}

func TestEnumerable_Single(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(int) bool
		want      int
		wantErr   bool
	}{
		{
			source:    []int{},
			predicate: nil,
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{1},
			predicate: nil,
			want:      1,
			wantErr:   false,
		},
		{
			source:    []int{1, 2},
			predicate: nil,
			want:      0,
			wantErr:   true,
		},
		{
			source:    []int{1, 2, 3},
			predicate: func(item int) bool { return item%2 == 0 },
			want:      2,
			wantErr:   false,
		},
		{
			source:    []int{1, 2, 3, 4},
			predicate: func(item int) bool { return item%2 == 0 },
			want:      0,
			wantErr:   true,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result, err := source.Single(test.predicate)
		if test.wantErr {
			if err == nil {
				t.Fatalf("Single() = %v, want error", result)
			}
		} else {
			if err != nil {
				t.Fatalf("Single() returned an error: %v", err)
			}
			if result != test.want {
				t.Fatalf("Single() = %v, want %v", result, test.want)
			}
		}
	}
}

func TestEnumerable_SingleOrDefault(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(int) bool
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item == 6 },
			want:      0,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 == 0 },
			want:      0,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 != 0 },
			want:      0,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 3 },
			want:      0,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.SingleOrDefault(test.predicate)
		if result != test.want {
			t.Fatalf("SingleOrDefault() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_SingleOrFallback(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(int) bool
		fallback  int
		want      int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item == 6 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 == 0 },
			fallback:  42,
			want:      42,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 != 0 },
			fallback:  42,
			want:      42,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.SingleOrFallback(test.fallback, test.predicate)
		if result != test.want {
			t.Fatalf("SingleOrFallback() = %v, want %v", result, test.want)
		}
	}
}

func TestEnumerable_SingleOrNil(t *testing.T) {
	tests := []struct {
		source    []int
		predicate func(int) bool
		want      *int
	}{
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item == 6 },
			want:      nil,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 == 0 },
			want:      nil,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item%2 != 0 },
			want:      nil,
		},
		{
			source:    []int{1, 2, 3, 4, 5},
			predicate: func(item int) bool { return item > 3 },
			want:      nil,
		},
	}
	for _, test := range tests {
		source := linq.FromSlice(test.source)
		result := source.SingleOrNil(test.predicate)
		if result != nil && test.want != nil && *result != *test.want {
			t.Fatalf("SingleOrNil() = %v, want %v", *result, *test.want)
		} else if result == nil && test.want != nil {
			t.Fatalf("SingleOrNil() = %v, want %v", result, *test.want)
		} else if result != nil && test.want == nil {
			t.Fatalf("SingleOrNil() = %v, want %v", *result, test.want)
		}
	}
}

func TestEnumerable_Skip(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.Skip(2).ToSlice()
	want := []int{3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Skip() = %v, want %v", result, want)
	}
}

func TestEnumerable_SkipLast(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.SkipLast(2).ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("SkipLast() = %v, want %v", result, want)
	}
}

func TestEnumerable_SkipWhile(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.SkipWhile(func(item int) bool {
		return item < 3
	}).ToSlice()
	want := []int{3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("SkipWhile() = %v, want %v", result, want)
	}
}

func TestEnumerable_Sum(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result, err := source.Sum[int]()
	want := 15
	if err != nil {
		t.Fatalf("Sum() returned an error: %v", err)
	}
	if result != want {
		t.Fatalf("Sum() = %v, want %v", result, want)
	}
}

func TestEnumerable_SumFunc(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "cherry"})
	result := source.SumFunc(func(item string) int {
		return len(item)
	})
	want := 17
	if result != want {
		t.Fatalf("SumFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Take(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.Take(3).ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Take() = %v, want %v", result, want)
	}
}

func TestEnumerable_TakeLast(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.TakeLast(3).ToSlice()
	want := []int{3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("TakeLast() = %v, want %v", result, want)
	}
}

func TestEnumerable_TakeWhile(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.TakeWhile(func(item int) bool {
		return item < 4
	}).ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("TakeWhile() = %v, want %v", result, want)
	}
}

func TestEnumerable_ToChannel(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	channel := make(chan int)
	go source.ToChannel(channel)
	var result []int
	for item := range channel {
		result = append(result, item)
	}
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ToChannel() = %v, want %v", result, want)
	}
}

func TestEnumerable_ToMap(t *testing.T) {
	source := linq.FromSlice([]string{"apple", "banana", "cherry"})
	result := source.ToMap(func(item string) string {
		return item
	}, func(item string) int {
		return len(item)
	})
	want := map[string]int{
		"apple":  5,
		"banana": 6,
		"cherry": 6,
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ToMap() = %v, want %v", result, want)
	}
}

func TestEnumerable_ToSlice(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ToSlice() = %v, want %v", result, want)
	}
}

func TestEnumerable_Union(t *testing.T) {
	seq1 := linq.FromSlice([]int{1, 2, 3})
	seq2 := linq.FromSlice([]int{3, 4, 5})
	result := seq1.Union[int](seq2).ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Union() = %v, want %v", result, want)
	}
}

func TestEnumerable_UnionEquatable(t *testing.T) {
	seq1 := linq.FromSlice([]time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	seq2 := linq.FromSlice([]time.Time{
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	result := seq1.UnionEquatable[time.Time](seq2).ToSlice()
	want := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("UnionEquatable() = %v, want %v", result, want)
	}
}

func TestEnumerable_UnionFunc(t *testing.T) {
	seq1 := linq.FromSlice([]int{1, 2, 3})
	seq2 := linq.FromSlice([]int{3, 4, 5})
	result := seq1.UnionFunc(seq2, func(a, b int) bool { return a == b }).ToSlice()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("UnionFunc() = %v, want %v", result, want)
	}
}

func TestEnumerable_Where(t *testing.T) {
	source := linq.FromSlice([]int{1, 2, 3, 4, 5})
	result := source.Where(func(item int) bool {
		return item%2 == 0
	}).ToSlice()
	want := []int{2, 4}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Where() = %v, want %v", result, want)
	}
}

func TestEnumerable_Zip(t *testing.T) {
	seq1 := linq.FromSlice([]string{"A", "B", "C"})
	seq2 := linq.FromSlice([]string{"1", "2", "3"})
	result := seq1.Zip(seq2, func(item1 string, item2 string) string {
		return item1 + item2
	}).ToSlice()
	want := []string{"A1", "B2", "C3"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Zip() = %v, want %v", result, want)
	}
}

func TestOrderable_ThenBy(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	result := linq.FromSlice([]Person{
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
	}).OrderByComparable(func(p Person) time.Time { return p.BirthDate }).ThenBy(func(p Person) string { return p.Name }).ToSlice()
	want := []Person{
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ThenBy() = %v, want %v", result, want)
	}
}

func TestOrderable_ThenDescendingBy(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	result := linq.FromSlice([]Person{
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
	}).OrderByComparable(func(p Person) time.Time { return p.BirthDate }).ThenDescendingBy(func(p Person) string { return p.Name }).ToSlice()
	want := []Person{
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ThenDescendingBy() = %v, want %v", result, want)
	}
}

func TestOrderable_ThenByComparable(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	result := linq.FromSlice([]Person{
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
	}).OrderBy(func(p Person) string { return p.Name }).ThenByComparable(func(p Person) time.Time { return p.BirthDate }).ToSlice()
	want := []Person{
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ThenByComparable() = %v, want %v", result, want)
	}
}

func TestOrderable_ThenDescendingByComparable(t *testing.T) {
	type Person struct {
		Name      string
		BirthDate time.Time
	}
	result := linq.FromSlice([]Person{
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
	}).OrderBy(func(p Person) string { return p.Name }).ThenDescendingByComparable(func(p Person) time.Time { return p.BirthDate }).ToSlice()
	want := []Person{
		{"Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Alice", time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Bob", time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Charlie", time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("ThenDescendingByComparable() = %v, want %v", result, want)
	}
}

func ValueOrDefault[T any](value *T) T {
	if value != nil {
		return *value
	}
	return *(new(T))
}

func NullValueOrDefault[T any, TResult any](value *T, f linq.Value[*T, TResult]) TResult {
	if value != nil {
		return f(value)
	}
	return *(new(TResult))
}

func BenchmarkTraditional(b *testing.B) {
	x := make([]int, 0, 100000000)
	for i := 0; i < 100000000; i++ {
		x = append(x, i)
	}
	b.ResetTimer() // Resetujemy licznik, ignorując czas przygotowania danych

	//x = append(x, 1)
	x = x[:10000]
	//x = x[len(x)-1:]
	dist := map[int]struct{}{}
	yy := make([]int, 0, len(x))
	for _, item := range x {
		if _, ok := dist[item]; !ok {
			dist[item] = struct{}{}
			yy = append(yy, item)
		}
	}
	sum := 5
	for _, item := range yy {
		sum += item
	}
	_ = sum
}

func BenchmarkLinq(b *testing.B) {
	x := make([]int, 0, 100000000)
	for i := 0; i < 100000000; i++ {
		x = append(x, i)
	}
	b.ResetTimer() // Resetujemy licznik, ignorując czas przygotowania danych

	dd := linq.FromSlice(x). //Append([]int{1}).
					Take(10000).
					Distinct[int]().
					Aggregate(5, func(accumulator int, item int) int {
			return accumulator + item
		})
	_ = dd
}
