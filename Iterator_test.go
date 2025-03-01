package linq

import (
	"reflect"
	"testing"
	"time"

	"github.com/thereisnoplanb/generic"
)

func TestIterator_Aggregate(t *testing.T) {
	type args struct {
		source      Iterator[int]
		seed        int
		accumulator generic.Accumulator[int, int]
		convert     []func(int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Aggregate empty source without result conversion",
			args: args{
				source: FromSlice([]int{}),
				seed:   100,
				accumulator: func(accumulator, object int) (result int) {
					return accumulator + object
				},
				convert: nil,
			},
			want: 100,
		},
		{
			name: "Aggregate empty source with result conversion",
			args: args{
				source: FromSlice([]int{}),
				seed:   100,
				accumulator: func(accumulator, object int) (result int) {
					return accumulator + object
				},
				convert: []func(int) int{
					func(accumulator int) int {
						return accumulator / 2
					},
				},
			},
			want: 50,
		},
		{
			name: "Aggregate source without result conversion",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				seed:   100,
				accumulator: func(accumulator, object int) (result int) {
					return accumulator + object
				},
				convert: nil,
			},
			want: 128,
		},
		{
			name: "Aggregate source with result conversion",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				seed:   100,
				accumulator: func(accumulator, object int) (result int) {
					return accumulator + object
				},
				convert: []func(int) int{
					func(accumulator int) int {
						return accumulator / 2
					},
				},
			},
			want: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Aggregate(tt.args.seed, tt.args.accumulator, tt.args.convert...); got != tt.want {
				t.Errorf("Iterator.Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_All(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "All empty source",
			args: args{
				source: FromSlice([]int{}),
				predicate: func(object int) (result bool) {
					return object > 0
				},
			},
			want: true,
		},
		{
			name: "All source, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 0
				},
			},
			want: true,
		},
		{
			name: "All source, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 4
				},
			},
			want: false,
		},
		{
			name: "All source, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 8
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.All(tt.args.predicate); got != tt.want {
				t.Errorf("Iterator.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Any(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Any empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: false,
		},
		{
			name: "Any source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: true,
		},
		{
			name: "Any source without predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: true,
		},
		{
			name: "Any empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: false,
		},
		{
			name: "Any source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: true,
		},
		{
			name: "Any source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: true,
		},
		{
			name: "Any source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Any(tt.args.predicate...); got != tt.want {
				t.Errorf("Iterator.Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Append(t *testing.T) {
	type args struct {
		source Iterator[int]
		values []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Append empty source without values",
			args: args{
				source: FromSlice([]int{}),
				values: []int{},
			},
			want: []int{},
		},
		{
			name: "Append source without values",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				values: []int{},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Append empty source with values",
			args: args{
				source: FromSlice([]int{}),
				values: []int{8, 9, 0},
			},
			want: []int{8, 9, 0},
		},
		{
			name: "Append source with values",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				values: []int{8, 9, 0},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Append(tt.args.values...).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Concat(t *testing.T) {
	type args struct {
		source   Iterator[int]
		sequence Iterator[int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Concat empty source with empty sequence",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{}),
			},
			want: []int{},
		},
		{
			name: "Concat source with empty sequence",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{}),
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Concat empty source with sequence",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{8, 9, 0}),
			},
			want: []int{8, 9, 0},
		},
		{
			name: "Concat source with sequence",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{8, 9, 0}),
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Concat(tt.args.sequence).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Concat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Contains(t *testing.T) {
	type args struct {
		source   Iterator[int]
		value    int
		comparer []generic.Equality[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Contains empty source",
			args: args{
				source: FromSlice([]int{}),
				value:  4,
			},
			want: false,
		},
		{
			name: "Contains source with value",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				value:  4,
			},
			want: true,
		},
		{
			name: "Contains source without value",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				value:  8,
			},
			want: false,
		},
		{
			name: "Contains empty source with equality comparer",
			args: args{
				source: FromSlice([]int{}),
				value:  4,
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: false,
		},
		{
			name: "Contains source with value with equality comparer",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				value:  4,
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: true,
		},
		{
			name: "Contains source without value with equality comparer",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				value:  8,
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Contains(tt.args.value, tt.args.comparer...); got != tt.want {
				t.Errorf("Iterator.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Count(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Count empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "Count source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: 7,
		},
		{
			name: "Count source without predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: 7,
		},
		{
			name: "Count empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "Count source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 7,
		},
		{
			name: "Count source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: 3,
		},
		{
			name: "Count source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Count(tt.args.predicate...); got != tt.want {
				t.Errorf("Iterator.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Distinct(t *testing.T) {
	type args struct {
		source   Iterator[int]
		comparer []generic.Equality[int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Distinct empty source",
			args: args{
				source: FromSlice([]int{}),
			},
			want: []int{},
		},
		{
			name: "Distinct source with distinct values",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Distinct source with not distinct values",
			args: args{
				source: FromSlice([]int{1, 7, 2, 6, 3, 5, 4, 4, 5, 3, 6, 2, 7, 1}),
			},
			want: []int{1, 7, 2, 6, 3, 5, 4},
		},
		{
			name: "Distinct empty source with equality comparer",
			args: args{
				source: FromSlice([]int{}),
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: []int{},
		},
		{
			name: "Distinct source with distinct values with equality comparer",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Distinct source with not distinct values with equality comparer",
			args: args{
				source: FromSlice([]int{1, 7, 2, 6, 3, 5, 4, 4, 5, 3, 6, 2, 7, 1}),
				comparer: []generic.Equality[int]{
					func(first, second int) (result bool) {
						return first == second
					},
				},
			},
			want: []int{1, 7, 2, 6, 3, 5, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Distinct(tt.args.comparer...).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_ElementAt(t *testing.T) {
	type args struct {
		source Iterator[int]
		index  int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "ElementAt, empty source",
			args: args{
				source: FromSlice([]int{}),
				index:  4,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "ElementAt, source with values, index in range",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:  4,
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "ElementAt, source with values, index out of range",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:  8,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "ElementAt, source with values, index out of range (below 0)",
			args: args{
				source: FromSlice([]int{}),
				index:  -1,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.source.ElementAt(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Iterator.First() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Iterator.ElementAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_ElementAtOrDefault(t *testing.T) {
	type args struct {
		source Iterator[int]
		index  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ElementAtOrDefault, empty source",
			args: args{
				source: FromSlice([]int{}),
				index:  4,
			},
			want: 0,
		},
		{
			name: "ElementAtOrDefault, source with values, index in range",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:  4,
			},
			want: 5,
		},
		{
			name: "ElementAtOrDefault, source with values, index out of range",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:  8,
			},
			want: 0,
		},
		{
			name: "ElementAtOrDefault, source with values, index out of range (below 0)",
			args: args{
				source: FromSlice([]int{}),
				index:  -1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.ElementAtOrDefault(tt.args.index); got != tt.want {
				t.Errorf("Iterator.ElementAtOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_ElementAtOrFallback(t *testing.T) {
	type args struct {
		source   Iterator[int]
		index    int
		fallback int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ElementAtOrFallback, empty source",
			args: args{
				source:   FromSlice([]int{}),
				index:    4,
				fallback: 100,
			},
			want: 100,
		},
		{
			name: "ElementAtOrFallback, source with values, index in range",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:    4,
				fallback: 100,
			},
			want: 5,
		},
		{
			name: "ElementAtOrFallback, source with values, index out of range",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				index:    8,
				fallback: 100,
			},
			want: 100,
		},
		{
			name: "ElementAtOrFallback, source with values, index out of range (below 0)",
			args: args{
				source:   FromSlice([]int{}),
				index:    -1,
				fallback: 100,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.ElementAtOrFallback(tt.args.index, tt.args.fallback); got != tt.want {
				t.Errorf("Iterator.ElementAtOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_First(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name      string
		args      args
		want      int
		wantError bool
	}{
		{
			name: "First empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "First source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want:      1,
			wantError: false,
		},
		{
			name: "First source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want:      1,
			wantError: false,
		},
		{
			name: "First empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "First source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      1,
			wantError: false,
		},
		{
			name: "First source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want:      5,
			wantError: false,
		},
		{
			name: "First source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want:      0,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.source.First(tt.args.predicate...)
			if (err != nil) != tt.wantError {
				t.Errorf("Iterator.First() error = %v, wantErr %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("Iterator.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_FirstOrDefault(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "FirstOrDefault empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "FirstOrDefault source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: 1,
		},
		{
			name: "FirstOrDefault source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: 1,
		},
		{
			name: "FirstOrDefault empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "FirstOrDefault source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 1,
		},
		{
			name: "FirstOrDefault source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: 5,
		},
		{
			name: "FirstOrDefault source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.FirstOrDefault(tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.FirstOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_FirstOrFallback(t *testing.T) {
	type args struct {
		source    Iterator[int]
		fallback  int
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "FirstOrFallback empty source without predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "FirstOrFallback source without predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
			},
			want: 1,
		},
		{
			name: "FirstOrFallback source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback:  100,
				predicate: []generic.Predicate[int]{nil},
			},
			want: 1,
		},
		{
			name: "FirstOrFallback empty source with predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "FirstOrFallback source with predicate, all values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 1,
		},
		{
			name: "FirstOrFallback source with predicate, only some values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: 5,
		},
		{
			name: "FirstOrFallback source with predicate, none value satisfies predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.FirstOrFallback(tt.args.fallback, tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.FirstOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Last(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name      string
		args      args
		want      int
		wantError bool
	}{
		{
			name: "Last empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Last source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want:      7,
			wantError: false,
		},
		{
			name: "Last source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want:      7,
			wantError: false,
		},
		{
			name: "Last empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Last source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 8
					},
				},
			},
			want:      7,
			wantError: false,
		},
		{
			name: "Last source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 4
					},
				},
			},
			want:      3,
			wantError: false,
		},
		{
			name: "Last source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.source.Last(tt.args.predicate...)
			if (err != nil) != tt.wantError {
				t.Errorf("Iterator.Last() error = %v, wantErr %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("Iterator.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_LastOrDefault(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "LastOrDefault empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "LastOrDefault source without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: 7,
		},
		{
			name: "LastOrDefault source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: 7,
		},
		{
			name: "LastOrDefault empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "LastOrDefault source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 8
					},
				},
			},
			want: 7,
		},
		{
			name: "LastOrDefault source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 4
					},
				},
			},
			want: 3,
		},
		{
			name: "LastOrDefault source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 0
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.LastOrDefault(tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.LastOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_LastOrFallback(t *testing.T) {
	type args struct {
		source    Iterator[int]
		fallback  int
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "LastOrFallback empty source without predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "LastOrFallback source without predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
			},
			want: 7,
		},
		{
			name: "LastOrFallback source with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback:  100,
				predicate: []generic.Predicate[int]{nil},
			},
			want: 7,
		},
		{
			name: "LastOrFallback empty source with predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "LastOrFallback source with predicate, all values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 8
					},
				},
			},
			want: 7,
		},
		{
			name: "LastOrFallback source with predicate, only some values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 4
					},
				},
			},
			want: 3,
		},
		{
			name: "LastOrFallback source with predicate, none value satisfies predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object < 0
					},
				},
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.LastOrFallback(tt.args.fallback, tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.LastOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Single(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name      string
		args      args
		want      int
		wantError bool
	}{
		{
			name: "Single empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single source with one item without predicate",
			args: args{
				source: FromSlice([]int{1}),
			},
			want:      1,
			wantError: false,
		},
		{
			name: "Single source with more than one item without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single source with one item with predicate nil",
			args: args{
				source:    FromSlice([]int{1}),
				predicate: []generic.Predicate[int]{nil},
			},
			want:      1,
			wantError: false,
		},
		{
			name: "Single source with more than one item with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single source with predicate, more than one value satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Single source with predicate, only one value satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object == 4
					},
				},
			},
			want:      4,
			wantError: false,
		},
		{
			name: "Single source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want:      0,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.source.Single(tt.args.predicate...)
			if (err != nil) != tt.wantError {
				t.Errorf("Iterator.Single() error = %v, wantErr %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("Iterator.Single() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SingleOrDefault(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "SingleOrDefault empty source without predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "SingleOrDefault source with only one value without predicate",
			args: args{
				source: FromSlice([]int{1}),
			},
			want: 1,
		},
		{
			name: "SingleOrDefault source with more than one value without predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: 0,
		},
		{
			name: "SingleOrDefault source with only one value with predicate nil",
			args: args{
				source:    FromSlice([]int{1}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: 1,
		},
		{
			name: "SingleOrDefault source with more than one value with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{nil},
			},
			want: 0,
		},
		{
			name: "SingleOrDefault empty source with predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "SingleOrDefault source with predicate, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 0,
		},
		{
			name: "SingleOrDefault source with predicate, only some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: 0,
		},
		{
			name: "SingleOrDefault source with predicate, only one value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object == 4
					},
				},
			},
			want: 4,
		},
		{
			name: "SingleOrDefault source with predicate, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.SingleOrDefault(tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.SingleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SingleOrFallback(t *testing.T) {
	type args struct {
		source    Iterator[int]
		fallback  int
		predicate []generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "SingleOrFallback empty source without predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "SingleOrFallback source with only one value without predicate",
			args: args{
				source:   FromSlice([]int{1}),
				fallback: 100,
			},
			want: 1,
		},
		{
			name: "SingleOrFallback source with more than one value without predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
			},
			want: 100,
		},
		{
			name: "SingleOrFallback source with only one value with predicate nil",
			args: args{
				source:    FromSlice([]int{1}),
				fallback:  100,
				predicate: []generic.Predicate[int]{nil},
			},
			want: 1,
		},
		{
			name: "SingleOrFallback source with more than one value with predicate nil",
			args: args{
				source:    FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback:  100,
				predicate: []generic.Predicate[int]{nil},
			},
			want: 100,
		},
		{
			name: "SingleOrFallback empty source with predicate",
			args: args{
				source:   FromSlice([]int{}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "SingleOrFallback source with predicate, all values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 0
					},
				},
			},
			want: 100,
		},
		{
			name: "SingleOrFallback source with predicate, only some values satisfy predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 4
					},
				},
			},
			want: 100,
		},
		{
			name: "SingleOrFallback source with predicate, only one value satisfies predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object == 4
					},
				},
			},
			want: 4,
		},
		{
			name: "SingleOrFallback source with predicate, none value satisfies predicate",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				fallback: 100,
				predicate: []generic.Predicate[int]{
					func(object int) (result bool) {
						return object > 8
					},
				},
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.source.SingleOrFallback(tt.args.fallback, tt.args.predicate...)
			if got != tt.want {
				t.Errorf("Iterator.SingleOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Where(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Where empty source",
			args: args{
				source: FromSlice([]int{}),
				predicate: func(object int) (result bool) {
					return object == 4
				},
			},
			want: []int{},
		},
		{
			name: "Where source, all values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 0
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Where source, some values satisfy predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 4
				},
			},
			want: []int{5, 6, 7},
		},
		{
			name: "Where source, none value satisfies predicate",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object > 8
				},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Where(tt.args.predicate).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Take(t *testing.T) {
	type args struct {
		source Iterator[int]
		count  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Take empty source, count -1",
			args: args{
				source: FromSlice([]int{}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "Take empty source, count 0",
			args: args{
				source: FromSlice([]int{}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "Take empty source, count 4",
			args: args{
				source: FromSlice([]int{}),
				count:  4,
			},
			want: []int{},
		},
		{
			name: "Take source, count -1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "Take source, count 0",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "Take source, count 4",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  4,
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "Take source, count > len",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  100,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Take(tt.args.count).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_TakeLast(t *testing.T) {
	type args struct {
		source Iterator[int]
		count  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "TakeLast empty source, count -1",
			args: args{
				source: FromSlice([]int{}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "TakeLast empty source, count 0",
			args: args{
				source: FromSlice([]int{}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "TakeLast empty source, count 4",
			args: args{
				source: FromSlice([]int{}),
				count:  4,
			},
			want: []int{},
		},
		{
			name: "TakeLast source, count -1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "TakeLast source, count 0",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "TakeLast source, count 4",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  4,
			},
			want: []int{4, 5, 6, 7},
		},
		{
			name: "TakeLast source, count > len",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  100,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.TakeLast(tt.args.count).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.TakeLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_TakeWhile(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "TakeWhile empty source, predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: func(object int) (result bool) {
					return object%3 == 0
				},
			},
			want: []int{},
		},
		{
			name: "TakeWhile source, predicate 1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object%3 == 0
				},
			},
			want: []int{1, 2},
		},
		{
			name: "TakeWhile source, count 0",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object%13 == 0
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.TakeWhile(tt.args.predicate).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.TakeWhile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_ToSlice(t *testing.T) {
	type args struct {
		iterator Iterator[int]
	}
	tests := []struct {
		name       string
		args       args
		wantResult []int
	}{
		{
			name: "ToSlice",
			args: args{
				iterator: FromSlice([]int{
					1,
					2,
					3,
				}),
			},
			wantResult: []int{
				1,
				2,
				3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := tt.args.iterator.ToSlice(); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToSlice() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_ToMap(t *testing.T) {
	type args struct {
		iterator      Iterator[int]
		keySelector   generic.KeySelector[int, int]
		valueSelector generic.ValueSelector[int, int]
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[int]int
	}{
		{
			name: "ToMap",
			args: args{
				iterator: FromSlice([]int{
					1,
					2,
					3,
				}),
				keySelector: func(object int) (key int) {
					return object
				},
				valueSelector: func(object int) (value int) {
					return object
				},
			},
			wantResult: map[int]int{
				1: 1,
				2: 2,
				3: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ToMap(tt.args.iterator, tt.args.keySelector, tt.args.valueSelector); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToMap() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestIterator_Skip(t *testing.T) {
	type args struct {
		source Iterator[int]
		count  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Skip empty source, count -1",
			args: args{
				source: FromSlice([]int{}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "Skip empty source, count 0",
			args: args{
				source: FromSlice([]int{}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "Skip empty source, count 4",
			args: args{
				source: FromSlice([]int{}),
				count:  4,
			},
			want: []int{},
		},
		{
			name: "Skip source, count -1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  -1,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Skip source, count 0",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  0,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "Skip source, count 4",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  4,
			},
			want: []int{5, 6, 7},
		},
		{
			name: "Skip source, count > len",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  100,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.Skip(tt.args.count).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SkipLast(t *testing.T) {
	type args struct {
		source Iterator[int]
		count  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SkipLast empty source, count -1",
			args: args{
				source: FromSlice([]int{}),
				count:  -1,
			},
			want: []int{},
		},
		{
			name: "SkipLast empty source, count 0",
			args: args{
				source: FromSlice([]int{}),
				count:  0,
			},
			want: []int{},
		},
		{
			name: "SkipLast empty source, count 4",
			args: args{
				source: FromSlice([]int{}),
				count:  4,
			},
			want: []int{},
		},
		{
			name: "SkipLast source, count -1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  -1,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "SkipLast source, count 0",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  0,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "SkipLast source, count 4",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  4,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "SkipLast source, count > len",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				count:  100,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.SkipLast(tt.args.count).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.SkipLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SkipWhile(t *testing.T) {
	type args struct {
		source    Iterator[int]
		predicate generic.Predicate[int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SkipWhile empty source, predicate",
			args: args{
				source: FromSlice([]int{}),
				predicate: func(object int) (result bool) {
					return object%3 == 0
				},
			},
			want: []int{},
		},
		{
			name: "SkipWhile source, predicate 1",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object%3 == 0
				},
			},
			want: []int{3, 4, 5, 6, 7},
		},
		{
			name: "SkipWhile source, predicate 2",
			args: args{
				source: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				predicate: func(object int) (result bool) {
					return object%13 == 0
				},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.SkipWhile(tt.args.predicate).ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.SkipWhile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SequenceEqual(t *testing.T) {
	type args struct {
		source   Iterator[int]
		sequence Iterator[int]
		comparer []generic.Equality[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "SequenceEqual empty source with empty sequence",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{}),
			},
			want: true,
		},
		{
			name: "SequenceEqual source with empty sequence",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{}),
			},
			want: false,
		},
		{
			name: "SequenceEqual empty source with sequence",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{1, 2, 3}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source:   FromSlice([]int{1, 2, 3}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 0}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - equal",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
			want: true,
		},

		//-----

		{
			name: "SequenceEqual empty source with empty sequence, with comparer",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: true,
		},
		{
			name: "SequenceEqual source with empty sequence, with comparer",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual empty source with sequence, with comparer",
			args: args{
				source:   FromSlice([]int{}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{1, 2, 3}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source:   FromSlice([]int{1, 2, 3}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 0}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - equal, with comparer",
			args: args{
				source:   FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				sequence: FromSlice([]int{1, 2, 3, 4, 5, 6, 7}),
				comparer: []generic.Equality[int]{
					func(x, y int) bool {
						return x == y
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.SequenceEqual(tt.args.sequence, tt.args.comparer...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.SequenceEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SequenceEqual_IEquatable(t *testing.T) {
	type args struct {
		source   Iterator[time.Time]
		sequence Iterator[time.Time]
		comparer []generic.Equality[time.Time]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "SequenceEqual empty source with empty sequence",
			args: args{
				source:   FromSlice([]time.Time{}),
				sequence: FromSlice([]time.Time{}),
			},
			want: true,
		},
		{
			name: "SequenceEqual source with empty sequence",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{}),
			},
			want: false,
		},
		{
			name: "SequenceEqual empty source with sequence",
			args: args{
				source: FromSlice([]time.Time{}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 1, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - equal",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			want: true,
		},

		//-----

		{
			name: "SequenceEqual empty source with empty sequence, with comparer",
			args: args{
				source:   FromSlice([]time.Time{}),
				sequence: FromSlice([]time.Time{}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: true,
		},
		{
			name: "SequenceEqual source with empty sequence, with comparer",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual empty source with sequence, with comparer",
			args: args{
				source: FromSlice([]time.Time{}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - not equal, with comparer",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 1, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: false,
		},
		{
			name: "SequenceEqual source with sequence - equal, with comparer",
			args: args{
				source: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				sequence: FromSlice([]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
				comparer: []generic.Equality[time.Time]{
					func(x, y time.Time) bool {
						return x.Equal(y)
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.source.SequenceEqual(tt.args.sequence, tt.args.comparer...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.SequenceEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
