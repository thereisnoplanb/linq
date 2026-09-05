package linq_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/thereisnoplanb/linq"
)

func TestFromSlice(t *testing.T) {
	input := []int{3, 1, 2}
	result := linq.FromSlice(input).ToSlice()
	if !reflect.DeepEqual(result, input) {
		t.Fatalf("FromSlice() = %v, want %v", result, input)
	}
}

func TestFromMap(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2}
	result := linq.FromMap(input).ToSlice()
	got := make(map[string]int, len(result))
	for _, pair := range result {
		got[pair.Key] = pair.Value
	}
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("FromMap() = %v, want %v", got, input)
	}
}

func TestFromString(t *testing.T) {
	result := linq.FromString("Go \u0182").ToSlice()
	want := []rune{'G', 'o', ' ', '\u0182'}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("FromString() = %v, want %v", result, want)
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  []string
	}{
		{name: "positive count", count: 3, want: []string{"x", "x", "x"}},
		{name: "zero count", count: 0, want: []string{}},
		{name: "negative count", count: -1, want: []string{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := linq.Repeat("x", test.count).ToSlice()
			if !reflect.DeepEqual(result, test.want) {
				t.Fatalf("Repeat() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name  string
		start int
		count int
		step  []int
		want  []int
	}{
		{name: "default step", start: 2, count: 4, want: []int{2, 3, 4, 5}},
		{name: "custom step", start: 2, count: 4, step: []int{2}, want: []int{2, 4, 6, 8}},
		{name: "zero step", start: 7, count: 3, step: []int{0}, want: []int{7, 7, 7}},
		{name: "non-positive count", start: 2, count: 0, want: []int{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := linq.Range(test.start, test.count, test.step...).ToSlice()
			if !reflect.DeepEqual(result, test.want) {
				t.Fatalf("Range() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestInfinite(t *testing.T) {
	result := linq.Infinite(context.Background(), 1, func(value int) int {
		return value + 1
	}).Take(4).ToSlice()
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Infinite() = %v, want %v", result, want)
	}
}

func TestInfinite_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := linq.Infinite(ctx, 1, func(value int) int {
		return value + 1
	}).ToSlice()
	if len(result) != 0 {
		t.Fatalf("Infinite() with canceled context = %v, want empty sequence", result)
	}
}

func TestFromChannel(t *testing.T) {
	channel := make(chan int, 3)
	channel <- 1
	channel <- 2
	channel <- 3
	close(channel)

	result := linq.FromChannel(context.Background(), channel).ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("FromChannel() = %v, want %v", result, want)
	}
}

func TestFromChannel_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	channel := make(chan int)

	result := linq.FromChannel(ctx, channel).ToSlice()
	if len(result) != 0 {
		t.Fatalf("FromChannel() with canceled context = %v, want empty sequence", result)
	}
}

func TestFromIterator(t *testing.T) {
	iterator := func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			if !yield(value) {
				return
			}
		}
	}
	result := linq.FromIterator(iterator).ToSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("FromIterator() = %v, want %v", result, want)
	}
}

func TestFromKeyValueIterator(t *testing.T) {
	iterator := func(yield func(string, int) bool) {
		for _, pair := range []struct {
			key   string
			value int
		}{
			{key: "a", value: 1},
			{key: "b", value: 2},
		} {
			if !yield(pair.key, pair.value) {
				return
			}
		}
	}
	result := linq.FromKeyValueIterator(iterator).ToSlice()
	want := []linq.KeyValuePair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("FromKeyValueIterator() = %v, want %v", result, want)
	}
}
