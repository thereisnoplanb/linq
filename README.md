# linq

LINQ-style sequence operations for Go, built on Go 1.27 iterator sequences and generics inspired by .NET LINQ.

## Installation

```powershell
go get github.com/thereisnoplanb/linq
```

```go
import "github.com/thereisnoplanb/linq"
```

Examples that format output additionally require `import "fmt"`.

## Core Types

### Enumerable[TSource]

`Enumerable[TSource]` represents a sequence of values. Operations are composed as deferred sequences and are executed when the result is iterated or materialized with `ToSlice`, `ToMap`, or another terminal operation.

```go
numbers := linq.FromSlice([]int{1, 2, 3})
result := numbers.Select(func(value int) int { return value * 2 }).ToSlice()
// result: []int{2, 4, 6}
```

### OrderedEnumerable[TSource]

`OrderedEnumerable[TSource]` represents an ordered sequence returned by `OrderBy` or `OrderDescendingBy`. Use `ThenBy` and `ThenDescendingBy` to add secondary ordering criteria.

### Grouping[TKey, TSource]

`Grouping[TKey, TSource]` is a group produced by `GroupBy`. `Key` contains the group key and the embedded sequence contains the associated source elements.

### KeyValuePair[TKey, TValue]

`KeyValuePair[TKey, TValue]` stores a map entry in `Key` and `Value` fields. It is the element type produced by `FromMap`.

## Creating Sequences

### FromSlice

```go
func FromSlice[TSlice ~[]TSource, TSource any](source TSlice) Enumerable[TSource]
```

Creates an `Enumerable` from a slice. Element order is preserved.

```go
numbers := linq.FromSlice([]int{3, 1, 2})
values := numbers.ToSlice()
// values: []int{3, 1, 2}
```

### FromMap

```go
func FromMap[TMap ~map[TKey]TValue, TKey comparable, TValue any](source TMap) Enumerable[KeyValuePair[TKey, TValue]]
```

Creates an `Enumerable` of key-value pairs from a map. Go does not specify map iteration order.

```go
entries := linq.FromMap(map[string]int{"a": 1})
entry := entries.ToSlice()[0]
// entry.Key == "a", entry.Value == 1
```

### FromString

```go
func FromString(source string) Enumerable[rune]
```

Creates a sequence of runes from a string. Invalid UTF-8 is decoded according to Go string-range rules.

```go
letters := linq.FromString("Go").ToSlice()
// letters: []rune{'G', 'o'}
```

### Repeat

```go
func Repeat[TSource any](element TSource, count int) Enumerable[TSource]
```

Creates a sequence containing `element` exactly `count` times. A non-positive count produces an empty sequence.

```go
values := linq.Repeat("x", 3).ToSlice()
// values: []string{"x", "x", "x"}
```

### Range

```go
func Range(start, count int, step ...int) Enumerable[int]
```

Creates `count` integers beginning at `start`. `step` defaults to `1`; a zero step repeats the same value. A non-positive count produces an empty sequence.

```go
values := linq.Range(2, 4, 2).ToSlice()
// values: []int{2, 4, 6, 8}
```

### Infinite

```go
func Infinite[TSource any](ctx context.Context, seed TSource, next Value[TSource, TSource]) Enumerable[TSource]
```

Creates a sequence that starts with `seed` and generates each subsequent value by applying `next` to the previous value. The sequence stops when `ctx` is canceled or when the consumer stops iterating. Since the sequence is infinite until one of these conditions occurs, use a limiting operation such as `Take` when materializing it.

```go
ctx := context.Background()
values := linq.Infinite(ctx, 1, func(value int) int { return value + 1 }).Take(4).ToSlice()
// values: []int{1, 2, 3, 4}
```

### FromChannel

```go
func FromChannel[TSource any](ctx context.Context, source <-chan TSource) Enumerable[TSource]
```

Creates an `Enumerable` from a channel. Iteration ends when the channel closes or when the supplied context is canceled.

```go
ctx := context.Background()
ch := make(chan int, 3)
ch <- 1
ch <- 2
close(ch)

values := linq.FromChannel(ctx, ch).ToSlice()
// values: []int{1, 2}
```

### FromIterator

```go
func FromIterator[TSource any](source iter.Seq[TSource]) Enumerable[TSource]
```

Creates an `Enumerable` from a Go iterator sequence. This lets you bridge existing `iter.Seq` APIs into the helper set without materializing the sequence first.

```go
seq := func(yield func(int) bool) {
	for _, value := range []int{1, 2, 3} {
		if !yield(value) {
			return
		}
	}
}

values := linq.FromIterator(seq).ToSlice()
// values: []int{1, 2, 3}
```

### FromKeyValueIterator

```go
func FromKeyValueIterator[TSource any, TKey any](source iter.Seq2[TKey, TSource]) Enumerable[KeyValuePair[TKey, TSource]]
```

Creates an `Enumerable` of `KeyValuePair` values from a Go key-value iterator. This is useful when adapting maps or iterators that yield `(key, value)` pairs.

```go
pairs := linq.FromKeyValueIterator(func(yield func(string, int) bool) {
	for key, value := range map[string]int{"a": 1, "b": 2} {
		if !yield(key, value) {
			return
		}
	}
}).ToSlice()
// pairs: []linq.KeyValuePair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}}
```

## Delegates and Constraints

The package defines reusable function types:

| Type | Purpose |
| --- | --- |
| `Predicate[T]` | Tests an element and returns `bool`. |
| `Action[T]` | Performs an action without returning a value. |
| `Accumulate[TAccumulator, T]` | Combines an accumulator with an element. |
| `Key[T, TKey]` | Selects a grouping or ordering key. |
| `Value[T, TValue]` | Selects a projected value. |
| `Equal[T]` | Compares two values for equality. |
| `Compare[T]` | Returns a negative, zero, or positive ordering result. |

`iEquatable[T]` and `iComparable[T]` are constraints used by the `Equatable` and `Comparable` method variants.

## Aggregation

### Aggregate

```go
func (source Enumerable[TSource]) Aggregate[TResult any](seed TResult, accumulate Accumulate[TResult, TSource], resultSelect ...func(accumulator TResult, count int) TResult) TResult
```

Folds the sequence into an accumulator. The optional result selector transforms the final accumulator; only its first value is used.

```go
sum := linq.FromSlice([]int{1, 2, 3}).Aggregate(0, func(total, value int) int { return total + value })
// sum: 6
```

### All and Any

```go
func (source Enumerable[TSource]) All(predicate Predicate[TSource]) bool
func (source Enumerable[TSource]) Any(predicate ...Predicate[TSource]) bool
```

`All` returns true when every element satisfies the predicate; it also returns true for an empty sequence. `Any` tests whether a sequence contains an element, or whether at least one element satisfies its optional predicate. Both methods stop as soon as their result is known.

```go
source := linq.FromSlice([]int{2, 4, 6})
allEven := source.All(func(value int) bool { return value%2 == 0 })
hasLarge := source.Any(func(value int) bool { return value > 5 })
// allEven: true, hasLarge: true
```

### Count

```go
func (source Enumerable[TSource]) Count(predicate ...Predicate[TSource]) int
```

Counts all elements, or only elements satisfying the predicate.

```go
count := linq.FromSlice([]int{1, 2, 3, 4}).Count(func(value int) bool { return value%2 == 0 })
// count: 2
```

### Sum and SumFunc

```go
func (source Enumerable[TSource]) Sum[T Number | ~string]() (T, error)
func (source Enumerable[TSource]) SumFunc[T Number | ~string](value Value[TSource, T]) T
```

`Sum` adds numeric elements or concatenates strings. `SumFunc` first maps each element to a number or string. `Sum` returns `ErrTypeIsNotNumberOrString` when the source type is incompatible.

```go
total, err := linq.FromSlice([]int{2, 3}).Sum[int]()
// total: 5, err: nil
```

### Average

```go
func (source Enumerable[TSource]) Average[T Real]() (float64, error)
func (source Enumerable[TSource]) AverageFunc[T Real](value Value[TSource, T]) (float64, error)
func (source Enumerable[TSource]) AverageComplex[T constraints.Complex]() (complex128, error)
func (source Enumerable[TSource]) AverageComplexFunc[T constraints.Complex](value Value[TSource, T]) (complex128, error)
```

`Average` computes the arithmetic mean of numeric source elements. `AverageFunc` first projects each source element to an integer or floating-point value, then computes the mean of those projected values. `AverageComplex` computes the arithmetic mean of complex source elements, while `AverageComplexFunc` first projects each source element to a complex value. Empty sequences return `ErrSourceContainsNoElements`; incompatible source types return `ErrTypeIsNotNumber` for the direct variants. Selector variants do not return `ErrTypeIsNotNumber` because their result type is constrained to the appropriate numeric type.

```go
average, err := linq.FromSlice([]int{2, 4}).Average[int]()
// average: 3, err: nil
```

### Min, Max, and MinMax

Each operation has direct, `Comparable`, `Func`, `By`, `ByComparable`, and `ByFunc` variants where applicable:

```go
min, err := linq.FromSlice([]int{4, 1, 3}).Min[int]()
max, err := linq.FromSlice([]int{4, 1, 3}).Max[int]()
minValue, maxValue, err := linq.FromSlice([]int{4, 1, 3}).MinMax[int]()
```

| Family | Comparison | Result |
| --- | --- | --- |
| `Min`, `Max`, `MinMax` | Operators on `cmp.Ordered` values | Value(s) from the sequence |
| `MinComparable`, `MaxComparable`, `MinMaxComparable` | `Compare` method | Element(s) from the sequence |
| `MinFunc`, `MaxFunc`, `MinMaxFunc` | `Compare[TSource]` function | Element(s) from the sequence |
| `MinBy`, `MaxBy`, `MinMaxBy` | Operators on selected values | Source element(s) |
| `MinByComparable`, `MaxByComparable`, `MinMaxByComparable` | `Compare` on selected values | Source element(s) |
| `MinByFunc`, `MaxByFunc`, `MinMaxByFunc` | Custom comparison on selected values | Source element(s) |

All variants return `ErrSourceContainsNoElements` for an empty source. `By` variants return source elements, not selected key values.

## Filtering and Projection

### Where

```go
func (source Enumerable[TSource]) Where(predicate Predicate[TSource]) Enumerable[TSource]
```

Returns elements for which `predicate` is true, preserving source order.

```go
values := linq.FromSlice([]int{1, 2, 3, 4}).Where(func(value int) bool { return value%2 == 0 }).ToSlice()
// values: []int{2, 4}
```

### Select

```go
func (source Enumerable[TSource]) Select[TResult any](value Value[TSource, TResult]) Enumerable[TResult]
```

Projects every source element into one result element.

```go
values := linq.FromSlice([]int{1, 2, 3}).Select(func(value int) string { return fmt.Sprint(value) }).ToSlice()
// values: []string{"1", "2", "3"}
```

### SelectMany

```go
func (source Enumerable[TSource]) SelectMany[TResult any](value Value[TSource, []TResult]) Enumerable[TResult]
```

Projects each source element into a slice and concatenates all produced slices in order.

```go
values := linq.FromSlice([]int{1, 2}).SelectMany(func(value int) []int { return []int{value, value * 10} }).ToSlice()
// values: []int{1, 10, 2, 20}
```

### Skip and Take

```go
func (source Enumerable[TSource]) Skip(count int) Enumerable[TSource]
func (source Enumerable[TSource]) SkipLast(count int) Enumerable[TSource]
func (source Enumerable[TSource]) Take(count int) Enumerable[TSource]
func (source Enumerable[TSource]) TakeLast(count int) Enumerable[TSource]
```

`Skip` removes elements from the beginning; `SkipLast` removes elements from the end. `Take` selects elements from the beginning; `TakeLast` selects elements from the end. `Skip` and `SkipLast` return all source elements for a non-positive count; `Take` and `TakeLast` return an empty sequence for a non-positive count.

```go
head := linq.FromSlice([]int{1, 2, 3, 4}).Take(2).ToSlice()
tail := linq.FromSlice([]int{1, 2, 3, 4}).TakeLast(2).ToSlice()
// head: []int{1, 2}, tail: []int{3, 4}
```

### SkipWhile and TakeWhile

```go
func (source Enumerable[TSource]) SkipWhile(predicate func(TSource) bool) Enumerable[TSource]
func (source Enumerable[TSource]) TakeWhile(predicate func(TSource) bool) Enumerable[TSource]
```

`SkipWhile` skips the consecutive prefix that satisfies the predicate, then returns the rest. `TakeWhile` returns the consecutive prefix that satisfies the predicate, then stops.

```go
skipped := linq.FromSlice([]int{1, 2, 3, 1}).SkipWhile(func(value int) bool { return value < 3 }).ToSlice()
taken := linq.FromSlice([]int{1, 2, 3, 1}).TakeWhile(func(value int) bool { return value < 3 }).ToSlice()
// skipped: []int{3, 1}, taken: []int{1, 2}
```

### Chunk

```go
func (source Enumerable[TSource]) Chunk[TResult []TSource](size int) Enumerable[TResult]
```

Splits a sequence into chunks of at most `size` elements. The final chunk may be smaller. `size < 1` causes a panic with `ErrSizeIsBelowOne`.

```go
chunks := linq.FromSlice([]int{1, 2, 3}).Chunk[[]int](2).ToSlice()
// chunks contains {1, 2}, then {3}
```

### Cast

```go
func (source Enumerable[TSource]) Cast[TResult any]() Enumerable[TResult]
```

Casts every element to `TResult`. A value that cannot be asserted as `TResult` causes a panic during enumeration.

```go
values := linq.FromSlice([]any{1, 2}).Cast[int]().ToSlice()
// values: []int{1, 2}
```

### Join, LeftJoin, RightJoin, and GroupJoin

```go
func (source Enumerable[TSource]) Join[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) JoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) JoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, equal Equal[TKey], join func(TSource, TInner) TResult) Enumerable[TResult]

func (source Enumerable[TSource]) LeftJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, *TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) RightJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(*TSource, TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) LeftJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, *TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) RightJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(*TSource, TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) LeftJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, equal Equal[TKey], join func(TSource, *TInner) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) RightJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, equal Equal[TKey], join func(*TSource, TInner) TResult) Enumerable[TResult]

func (source Enumerable[TSource]) GroupJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, Enumerable[TInner]) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) GroupJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, Enumerable[TInner]) TResult) Enumerable[TResult]
func (source Enumerable[TSource]) GroupJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey func(TSource) TKey, innerKey func(TInner) TKey, join func(TSource, Enumerable[TInner]) TResult, equal Equal[TKey]) Enumerable[TResult]
```

`Join` correlates two sequences by matching keys and produces a result for each matching pair. `LeftJoin` and `RightJoin` preserve all elements from the left or right side respectively, while `GroupJoin` returns a grouped sequence of matches for each outer element. The `Equatable` variants compare keys by calling the `Equal` method, while the `Func` variants use the supplied `Equal[T]` function instead of the default comparable comparison.

```go
customers := linq.FromSlice([]string{"A", "B"})
orders := linq.FromSlice([]struct{ Customer string; Amount int }{{"A", 10}, {"A", 25}, {"C", 5}})

result := customers.Join(
	orders,
	func(customer string) string { return customer },
	func(order struct{ Customer string; Amount int }) string { return order.Customer },
	func(customer string, order struct{ Customer string; Amount int }) string {
		return customer + ":" + fmt.Sprint(order.Amount)
	},
).ToSlice()
// result: []string{"A:10", "A:25"}
```

## Set Operations

Each set operation has three comparison variants: the default operator variant, an `Equatable` variant using the `Equal` method, and a `Func` variant using a supplied `Equal[T]` function.

### Contains

```go
func (source Enumerable[TSource]) Contains[T comparable](value TSource) bool
func (source Enumerable[TSource]) ContainsEquatable[T iEquatable[T]](value TSource) bool
func (source Enumerable[TSource]) ContainsFunc(value TSource, equal Equal[TSource]) bool
```

Tests whether one value occurs in the source. Enumeration stops at the first match.

```go
found := linq.FromSlice([]int{1, 2, 3}).Contains[int](2)
// found: true
```

### ContainsAny and ContainsAll

```go
func (source Enumerable[TSource]) ContainsAny[T comparable](values ...TSource) bool
func (source Enumerable[TSource]) ContainsAnyEquatable[T iEquatable[T]](values ...TSource) bool
func (source Enumerable[TSource]) ContainsAnyFunc(equal Equal[TSource], values ...TSource) bool

func (source Enumerable[TSource]) ContainsAll[T comparable](values ...TSource) bool
func (source Enumerable[TSource]) ContainsAllEquatable[T iEquatable[T]](values ...TSource) bool
func (source Enumerable[TSource]) ContainsAllFunc(equal Equal[TSource], values ...TSource) bool
```

`ContainsAny` returns true when at least one requested value occurs. `ContainsAll` returns true when every requested value occurs. Empty `values` therefore makes `ContainsAny` false and `ContainsAll` true.

```go
source := linq.FromSlice([]int{1, 2, 3})
anyValue := source.ContainsAny[int](0, 2)
allValues := source.ContainsAll[int](1, 3)
// anyValue: true, allValues: true
```

### Distinct

```go
func (source Enumerable[TSource]) Distinct[T comparable]() Enumerable[TSource]
func (source Enumerable[TSource]) DistinctEquatable[T iEquatable[T]]() Enumerable[TSource]
func (source Enumerable[TSource]) DistinctFunc(equal Equal[TSource]) Enumerable[TSource]
```

Returns the first occurrence of each value, preserving source order.

```go
values := linq.FromSlice([]int{1, 2, 1, 3, 2}).Distinct[int]().ToSlice()
// values: []int{1, 2, 3}
```

### Except

```go
func (source Enumerable[TSource]) Except[T comparable](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) ExceptEquatable[T iEquatable[T]](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) ExceptFunc(other Enumerable[TSource], equal Equal[TSource]) Enumerable[TSource]
```

Returns source elements that do not occur in `other`. The source order and duplicate source elements are preserved.

```go
values := linq.FromSlice([]int{1, 2, 2, 3}).Except[int](linq.FromSlice([]int{2})).ToSlice()
// values: []int{1, 3}
```

### Intersect

```go
func (source Enumerable[TSource]) Intersect[T comparable](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) IntersectEquatable[T iEquatable[T]](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) IntersectFunc(other Enumerable[TSource], equal Equal[TSource]) Enumerable[TSource]
```

Returns source elements that occur in `other`. The source order and duplicate source elements are preserved.

```go
values := linq.FromSlice([]int{1, 2, 2, 3}).Intersect[int](linq.FromSlice([]int{2, 4})).ToSlice()
// values: []int{2, 2}
```

### Union

```go
func (source Enumerable[TSource]) Union[T comparable](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) UnionEquatable[T iEquatable[T]](other Enumerable[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) UnionFunc(other Enumerable[TSource], equal Equal[TSource]) Enumerable[TSource]
```

Returns distinct elements from both sequences. Source elements are returned first, followed by new elements from `other`.

```go
values := linq.FromSlice([]int{1, 2, 2}).Union[int](linq.FromSlice([]int{2, 3})).ToSlice()
// values: []int{1, 2, 3}
```

### OfType

```go
func (source Enumerable[TSource]) OfType[TResult any]() Enumerable[TResult]
```

Filters a sequence to the elements whose runtime type is `TResult`. The source order is preserved, and values that do not match `TResult` are omitted.

```go
values := linq.FromSlice([]any{1, "2", 3, "4"}).OfType[int]().ToSlice()
// values: []int{1, 3}
```

## Element Selection

### ElementAt

```go
func (source Enumerable[TSource]) ElementAt(index int) (TSource, error)
func (source Enumerable[TSource]) ElementAtOrDefault(index int) TSource
func (source Enumerable[TSource]) ElementAtOrFallback(index int, fallback TSource) TSource
func (source Enumerable[TSource]) ElementAtOrNil(index int) *TSource
```

All indexes are zero-based and count from the beginning. `ElementAt` returns `ErrIndexOutOfRange` for a negative or out-of-range index. The other variants return the zero value, fallback, or `nil` respectively.

```go
value, err := linq.FromSlice([]string{"a", "b"}).ElementAt(1)
// value: "b", err: nil
```

### First

```go
func (source Enumerable[TSource]) First(predicate ...Predicate[TSource]) (TSource, error)
func (source Enumerable[TSource]) FirstOrDefault(predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) FirstOrFallback(fallback TSource, predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) FirstOrNil(predicate ...Predicate[TSource]) *TSource
```

Returns the first source element or the first element satisfying the optional first non-nil predicate. `First` reports `ErrSourceContainsNoElements` for an empty source and `ErrNoElementSatisfiesTheConditionInPredicate` when filtering finds no match. The other variants return a zero value, fallback, or `nil`.

```go
value, err := linq.FromSlice([]int{1, 2, 3}).First(func(value int) bool { return value > 1 })
// value: 2, err: nil
```

### Last

```go
func (source Enumerable[TSource]) Last(predicate ...Predicate[TSource]) (TSource, error)
func (source Enumerable[TSource]) LastOrDefault(predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) LastOrFallback(fallback TSource, predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) LastOrNil(predicate ...Predicate[TSource]) *TSource
```

Returns the last source element or the last element satisfying the optional predicate. `Last` returns the same no-element and no-match errors as `First`; all `Last` variants enumerate the complete source to determine the final match.

```go
value, err := linq.FromSlice([]int{1, 2, 3, 2}).Last(func(value int) bool { return value == 2 })
// value: 2, err: nil
```

### Single

```go
func (source Enumerable[TSource]) Single(predicate ...Predicate[TSource]) (TSource, error)
func (source Enumerable[TSource]) SingleOrDefault(predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) SingleOrFallback(fallback TSource, predicate ...Predicate[TSource]) TSource
func (source Enumerable[TSource]) SingleOrNil(predicate ...Predicate[TSource]) *TSource
```

Returns the only source element or the only element satisfying the optional predicate. `Single` returns `ErrSourceContainsNoElements`, `ErrSourceHasMoreThanOneElement`, `ErrNoElementSatisfiesTheConditionInPredicate`, or `ErrMoreThanOneElementSatisfiesTheConditionInPredicate` as appropriate. The non-error variants return a zero value, fallback, or `nil` when the result is not unique.

```go
value, err := linq.FromSlice([]int{1, 2, 3}).Single(func(value int) bool { return value == 2 })
// value: 2, err: nil
```

## Ordering

### Order and OrderDescending

```go
func (source Enumerable[TSource]) Order[T cmp.Ordered]() Enumerable[TSource]
func (source Enumerable[TSource]) OrderComparable[T iComparable[T]]() Enumerable[TSource]
func (source Enumerable[TSource]) OrderFunc(compare Compare[TSource]) Enumerable[TSource]
func (source Enumerable[TSource]) OrderDescending[T cmp.Ordered]() Enumerable[TSource]
func (source Enumerable[TSource]) OrderDescendingComparable[T iComparable[T]]() Enumerable[TSource]
func (source Enumerable[TSource]) OrderDescendingFunc(compare Compare[TSource]) Enumerable[TSource]
```

Sorts source elements directly in ascending or descending order. The result is an `Enumerable`.

```go
ascending := linq.FromSlice([]int{3, 1, 2}).Order[int]().ToSlice()
descending := linq.FromSlice([]int{3, 1, 2}).OrderDescending[int]().ToSlice()
// ascending: []int{1, 2, 3}, descending: []int{3, 2, 1}
```

### OrderBy and OrderDescendingBy

```go
func (source Enumerable[TSource]) OrderBy[TValue cmp.Ordered](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source Enumerable[TSource]) OrderByComparable[TValue iComparable[TValue]](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source Enumerable[TSource]) OrderByFunc[TValue any](selector Value[TSource, TValue], compare Compare[TValue]) OrderedEnumerable[TSource]
func (source Enumerable[TSource]) OrderDescendingBy[TValue cmp.Ordered](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source Enumerable[TSource]) OrderDescendingByComparable[TValue iComparable[TValue]](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source Enumerable[TSource]) OrderDescendingByFunc[TValue any](selector Value[TSource, TValue], compare Compare[TValue]) OrderedEnumerable[TSource]
```

Orders source elements using a selected key and returns an `OrderedEnumerable`, which can receive secondary criteria through `ThenBy` methods. Ordering is stable for key-based methods.

```go
type person struct { Name string; Age int }
people := linq.FromSlice([]person{{"A", 30}, {"B", 20}})
ordered := people.OrderBy(func(value person) int { return value.Age }).ToSlice()
// ordered: B, then A
```

### ThenBy and ThenDescendingBy

```go
func (source OrderedEnumerable[TSource]) ThenBy[TValue cmp.Ordered](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source OrderedEnumerable[TSource]) ThenByComparable[TValue iComparable[TValue]](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source OrderedEnumerable[TSource]) ThenByFunc[TValue any](selector Value[TSource, TValue], compare Compare[TValue]) OrderedEnumerable[TSource]
func (source OrderedEnumerable[TSource]) ThenDescendingBy[TValue cmp.Ordered](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source OrderedEnumerable[TSource]) ThenDescendingByComparable[TValue iComparable[TValue]](selector Value[TSource, TValue]) OrderedEnumerable[TSource]
func (source OrderedEnumerable[TSource]) ThenDescendingByFunc[TValue any](selector Value[TSource, TValue], compare Compare[TValue]) OrderedEnumerable[TSource]
```

Adds a secondary ordering criterion after an existing `OrderBy` or `OrderDescendingBy` criterion. The `ThenBy` variants sort ascending, while the `ThenDescendingBy` variants sort descending. The `Comparable` variants call `Compare`; `Func` variants call the supplied comparison function.

```go
ordered := people.OrderBy(func(value person) int { return value.Age }).ThenDescendingBy(func(value person) string { return value.Name }).ToSlice()
```

## Sequence Composition

### Append, Prepend, and Concat

```go
func (source Enumerable[TSource]) Append(other ...TSource) Enumerable[TSource]
func (source Enumerable[TSource]) Prepend(other ...TSource) Enumerable[TSource]
func (source Enumerable[TSource]) Concat(other Enumerable[TSource]) Enumerable[TSource]
```

`Append` adds values after the source, `Prepend` adds values before it, and `Concat` appends another sequence. Existing order is preserved.

```go
values := linq.FromSlice([]int{2, 3}).Prepend(1).Append(4).ToSlice()
// values: []int{1, 2, 3, 4}
```

### Reverse

```go
func (source Enumerable[TSource]) Reverse() Enumerable[TSource]
```

Returns the source elements in reverse order. The source is fully materialized when the result is enumerated.

```go
values := linq.FromSlice([]int{1, 2, 3}).Reverse().ToSlice()
// values: []int{3, 2, 1}
```
### Shuffle

```go
func (source Enumerable[TSource]) Shuffle() Enumerable[TSource]
```

Returns a new sequence with the same elements in a randomized order. The operation is fully materialized before enumeration; each call uses Go's random shuffling implementation.

```go
values := linq.FromSlice([]int{1, 2, 3, 4}).Shuffle().ToSlice()
// values is a random permutation of [1, 2, 3, 4]
```

## Grouping

### GroupBy

```go
func (source Enumerable[TSource]) GroupBy[TKey comparable, TResult Grouping[TKey, TSource]](key func(TSource) TKey) Enumerable[TResult]
```

Groups source elements by a comparable key. Elements retain their source order within each group; map iteration means group order is not guaranteed.

```go
groups := linq.FromSlice([]string{"a", "bb", "c"}).GroupBy[int, linq.Grouping[int, string]](func(value string) int { return len(value) })
// group 1 contains "a", "c"; group 2 contains "bb"
```
## Materialization and Helpers

### SequenceEqual

```go
func (source Enumerable[TSource]) SequenceEqual[TValue comparable](other Enumerable[TSource]) bool
func (source Enumerable[TSource]) SequenceEqualEquatable[TValue iEquatable[TValue]](other Enumerable[TSource]) bool
func (source Enumerable[TSource]) SequenceEqualFunc(other Enumerable[TSource], equal Equal[TSource]) bool
```

Compares corresponding elements in order and requires equal sequence lengths. Comparison stops at the first mismatch.

```go
equal := linq.FromSlice([]int{1, 2}).SequenceEqual[int](linq.FromSlice([]int{1, 2}))
// equal: true
```

### Zip

```go
func (source Enumerable[TFirst]) Zip[TSecond, TResult any](other Enumerable[TSecond], zip func(TFirst, TSecond) TResult) Enumerable[TResult]
```

Combines corresponding elements from two sequences. The result has the length of the shorter sequence.

```go
values := linq.FromSlice([]int{1, 2}).Zip(linq.FromSlice([]string{"a", "b"}), func(number int, text string) string {
	return fmt.Sprintf("%d%s", number, text)
}).ToSlice()
// values: []string{"1a", "2b"}
```

### ToChannel

```go
func (source Enumerable[TSource]) ToChannel(result chan<- TSource)
```

Writes the sequence to a channel and closes it after iteration completes. This is useful for streaming values to a consumer without materializing the whole sequence first.

```go
ch := make(chan int, 3)
linq.FromSlice([]int{1, 2, 3}).ToChannel(ch)
// ch receives 1, 2, 3 and then closes
```

### ToSlice and ToMap

```go
func (source Enumerable[TSource]) ToSlice() []TSource
func (source Enumerable[TSource]) ToMap[TKey comparable, TValue any](key func(TSource) TKey, value func(TSource) TValue) map[TKey]TValue
```

Materializes a sequence. `ToSlice` preserves order and returns a non-nil empty slice for an empty source. `ToMap` uses selectors; if a key occurs more than once, the last value replaces the previous value.

```go
result := linq.FromSlice([]string{"go", "rust"}).ToMap(
	func(value string) int { return len(value) },
	func(value string) string { return value },
)
// result[2] == "rust"
```

### ForEach

```go
func (source Enumerable[TSource]) ForEach(action Action[TSource])
```

Invokes `action` for each element in source order. It does nothing for an empty sequence.

```go
linq.FromSlice([]int{1, 2}).ForEach(func(value int) { fmt.Println(value) })
```

## Exceptions and Errors

Go methods return errors rather than throwing exceptions, except for documented input/type failures such as `Chunk` size and `Cast` assertions.

| Error | Meaning |
| --- | --- |
| `ErrSourceContainsNoElements` | The source sequence is empty. |
| `ErrNoElementSatisfiesTheConditionInPredicate` | No element satisfies a predicate. |
| `ErrMoreThanOneElementSatisfiesTheConditionInPredicate` | More than one element satisfies a predicate. |
| `ErrSourceHasMoreThanOneElement` | The source contains more than one element. |
| `ErrSizeIsBelowOne` | A chunk size is less than one. |
| `ErrIndexOutOfRange` | An index is negative or outside the sequence bounds. |
| `ErrTypeIsNotNumberOrString` | A sum source type is neither numeric nor a string type. |
| `ErrTypeIsNotNumber` | An average source type is not numeric. |

## Notes on Deferred Enumeration

Most sequence-producing methods return a new `Enumerable` without immediately consuming the source. Enumerating the result consumes the source. A sequence should therefore be treated as an iterator: repeated enumeration may repeat side effects or observe changed source data.

## License

See [LICENSE](LICENSE).