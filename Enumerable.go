package linq

import (
	"cmp"
	"iter"
	"math"
	"math/rand/v2"
	"slices"
	"sort"

	"golang.org/x/exp/constraints"
)

// Enumerable represents a sequence of values of type TSource.
//
// It provides deferred, pull-based iteration and operations for transforming,
// filtering, ordering, and aggregating the sequence.
type Enumerable[TSource any] struct {
	enumerable[TSource]
}

type enumerable[TSource any] iter.Seq[TSource]

// Applies an accumulate function over a sequence. The specified seed value is used as the initial accumulator value.
//
// # Parameters
//
//	seed TResult
//
// The initial accumulator value.
//
//	accumulate Accumulate[TResult, TSource]
//
// An accumulator function to be invoked on each element.
//
//	resultSelect func(accumulator TResult, count int) TResult
//
// A function to transform the final accumulator value and the count of elements processed, into the result value. [OPTIONAL]
//
// # Returns
//
//	result TResult
//
// The final accumulator value.
func (source enumerable[TSource]) Aggregate[TResult any](seed TResult, accumulate Accumulate[TResult, TSource], resultSelect ...func(accumulator TResult, count int) TResult) (result TResult) {
	result = seed
	count := 0
	for item := range source {
		result = accumulate(result, item)
		count++
	}
	if len(resultSelect) > 0 && resultSelect[0] != nil {
		result = resultSelect[0](result, count)
	}
	return result
}

// Determines whether all elements of a sequence satisfy a condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result bool
//
// True if every element of the source sequence passes the test in the specified predicate, or if the sequence is empty; otherwise, false.
func (source enumerable[TSource]) All(predicate Predicate[TSource]) (result bool) {
	for item := range source {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Determines whether any element of a sequence exists or satisfies a condition.
//
// # Parameters
//
//	predicate Predicate[TSource] [OPTIONAL]
//
// A function to test each element for a condition. If omitted or nil, only sequence emptiness is tested. [OPTIONAL]
//
// # Returns
//
//	result bool
//
// True if the source sequence is not empty when no predicate is supplied, or if at least one element passes the test in the specified predicate; otherwise, false.
func (source enumerable[TSource]) Any(predicate ...Predicate[TSource]) (result bool) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				return true
			}
		}
		return false
	}
	for range source {
		return true
	}
	return false
}

// Appends values to the end of the sequence.
//
// # Parameters
//
//	other ...TSource
//
// The values to append to source.
//
// # Returns
//
//	result Enumerable[TSource]
//
// A new sequence containing the source elements followed by other.
func (source enumerable[TSource]) Append(other ...TSource) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if !yield(item) {
					return
				}
			}
			for _, item := range other {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Computes the average of a sequence of numeric values.
//
// # Returns
//
//	result float64
//
// The average of the values in the sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrTypeIsNotNumber
//     When the source element type is not a real numeric type.
//   - linq.ErrSourceContainsNoElements
//     When source contains no elements.
func (source enumerable[TSource]) Average[T Real]() (result float64, err error) {
	if _, ok := (any(*new(TSource))).(T); !ok {
		return result, ErrTypeIsNotNumber
	}

	count := 0
	var sum T

	for value := range source {
		sum += any(value).(T)
		count++
	}
	if count == 0 {
		return math.NaN(), ErrSourceContainsNoElements
	}
	return float64(sum) / float64(count), nil
}

// Computes the average of the values produced by applying a selector function to the elements of a sequence.
//
// # Parameters
//
//	value Value[TSource, T]
//
// A function that transforms each source element to an integer or floating-point value.
//
// # Returns
//
//	result float64
//
// The average of the selected values.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When source contains no elements.
func (source enumerable[TSource]) AverageFunc[T Real](value Value[TSource, T]) (result float64, err error) {
	count := 0
	var sum T

	for item := range source {
		sum += value(item)
		count++
	}
	if count == 0 {
		return math.NaN(), ErrSourceContainsNoElements
	}
	return float64(sum) / float64(count), nil
}

// Computes the average of a sequence of complex numeric values.
//
// # Returns
//
//	result complex128
//
// The average of the values in the sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrTypeIsNotNumber
//     When the source element type is not a complex numeric type.
//   - linq.ErrSourceContainsNoElements
//     When source contains no elements.
func (source enumerable[TSource]) AverageComplex[T constraints.Complex]() (result complex128, err error) {
	if _, ok := (any(*new(TSource))).(T); !ok {
		return result, ErrTypeIsNotNumber
	}

	count := 0
	var sum T

	for item := range source {
		sum += any(item).(T)
		count++
	}
	if count == 0 {
		return complex(math.NaN(), math.NaN()), ErrSourceContainsNoElements
	}
	return complex128(sum) / complex(float64(count), 0), nil
}

// Computes the average of the complex values produced by applying a selector function to the elements of a sequence.
//
// # Parameters
//
//	value Value[TSource, T]
//
// A function that transforms each source element into a complex numeric value.
//
// # Returns
//
//	result complex128
//
// The average of the selected values.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When source contains no elements.
func (source enumerable[TSource]) AverageComplexFunc[T constraints.Complex](value Value[TSource, T]) (result complex128, err error) {
	count := 0
	var sum T

	for item := range source {
		sum += value(item)
		count++
	}
	if count == 0 {
		return complex(math.NaN(), math.NaN()), ErrSourceContainsNoElements
	}
	return complex128(sum) / complex(float64(count), 0), nil
}

// Casts the elements of a sequence to the specified type.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] that contains each element of the source sequence cast to the specified type.
//
// # Panics
//
// When an element in the sequence cannot be cast to type TResult.
func (source enumerable[TSource]) Cast[TResult any]() (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			for item := range source {
				if !yield((any(item)).(TResult)) {
					return
				}
			}
		},
	}
}

// Splits the elements of a sequence into chunks of size at most size.
//
// # Parameters
//
//	size int
//
// The maximum size of each chunk.
//
// # Returns
//
//	result Enumerable[[]TSource]
//
// An Enumerable[[]TSource] containing the chunks of the input sequence, where each chunk is a slice of TSource values.
//
// # Remarks
//
// Each chunk except the last one has size elements. The last chunk contains the remaining elements and may be smaller.
//
// # Panics
//
// When size is below 1.
func (source enumerable[TSource]) Chunk[TResult []TSource](size int) (result Enumerable[TResult]) {
	if size < 1 {
		panic(ErrSizeIsBelowOne)
	}
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			chunk := make([]TSource, size)
			i := 0
			for item := range source {
				chunk[i] = item
				i++
				if i == size {
					if !yield(chunk) {
						return
					}
					chunk = make([]TSource, size)
					i = 0
				}
			}
			if i > 0 {
				if !yield(chunk[:i:i]) {
					return
				}
			}
		},
	}
}

// Concatenates the source sequence with another sequence.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// The sequence to append to the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] containing the source elements followed by the elements of other.
func (source enumerable[TSource]) Concat(other Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if !yield(item) {
					return
				}
			}
			for item := range other.enumerable {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Determines whether a sequence contains a specified element by using the default equality operator.
//
// # Parameters
//
//	value TSource
//
// The value to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains an element that has the specified value; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a matching element is found.
// Elements are compared with the specified value using Go's == operator.
func (source enumerable[TSource]) Contains[T comparable](value TSource) (result bool) {
	for item := range source {
		if (any(value)).(T) == (any(item)).(T) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains a specified element by using the Equal method.
//
// # Parameters
//
//	value TSource
//
// The value to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains an element that is equal to the specified value; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a matching element is found.
// Elements are compared with the specified value by calling the Equal method.
func (source enumerable[TSource]) ContainsEquatable[T iEquatable[T]](value TSource) (result bool) {
	for item := range source {
		if (any(value)).(T).Equal((any(item)).(T)) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains a specified element by using a specified equality function.
//
// # Parameters
//
//	value TSource
//
// The value to locate in the sequence.
//
//	equal Equal[TSource]
//
// A function to compare elements with the specified value.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains an element for which the equality function returns true; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a matching element is found.
// Elements are compared with the specified value by calling the equality function.
func (source enumerable[TSource]) ContainsFunc(value TSource, equal Equal[TSource]) (result bool) {
	for item := range source {
		if equal(item, value) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains any of the specified elements by using the equality operator.
//
// # Parameters
//
//	values ...TSource
//
// The list of values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains any element from the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as any matching element is found.
// Elements are compared with the specified values using Go's == operator.
func (source enumerable[TSource]) ContainsAny[T comparable](values ...TSource) (result bool) {
	for _, value := range values {
		if source.Contains[T](value) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains any of the specified elements by using the Equal method.
//
// # Parameters
//
//	values ...TSource
//
// The list of values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains any element that is equal to one of the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as any matching element is found.
// Elements are compared with the specified values by calling the Equal method.
func (source enumerable[TSource]) ContainsAnyEquatable[T iEquatable[T]](values ...TSource) (result bool) {
	for _, value := range values {
		if source.ContainsEquatable[T](value) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains any of the specified elements by using a specified equality function.
//
// # Parameters
//
//	equal Equal[TSource]
//
// A function to compare elements with the specified values.
//
//	values ...TSource
//
// The list of values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains any element for which the equality function returns true when compared with one of the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as any matching element is found.
// Elements are compared with the specified values by calling the equality function.
func (source enumerable[TSource]) ContainsAnyFunc(equal Equal[TSource], values ...TSource) (result bool) {
	for _, value := range values {
		if source.ContainsFunc(value, equal) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains all specified elements by using the equality operator.
//
// # Parameters
//
//	values ...TSource
//
// The values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains all elements from the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a specified value is not found.
// Elements are compared with the specified values using Go's == operator.
func (source enumerable[TSource]) ContainsAll[T comparable](values ...TSource) (result bool) {
	for _, value := range values {
		if !source.Contains[T](value) {
			return false
		}
	}
	return true
}

// Determines whether a sequence contains all specified elements by using the Equal method.
//
// # Parameters
//
//	values ...TSource
//
// The values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains all elements that are equal to the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a specified value is not found.
// Elements are compared with the specified values by calling the Equal method.
func (source enumerable[TSource]) ContainsAllEquatable[T iEquatable[T]](values ...TSource) (result bool) {
	for _, value := range values {
		if !source.ContainsEquatable[T](value) {
			return false
		}
	}
	return true
}

// Determines whether a sequence contains all specified elements by using a specified equality function.
//
// # Parameters
//
//	equal Equal[TSource]
//
// A function to compare elements with the specified values.
//
//	values ...TSource
//
// The values to locate in the sequence.
//
// # Returns
//
//	result bool
//
// True if the source sequence contains all elements for which the equality function returns true when compared with the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as a specified value is not found.
// Elements are compared with the specified values by calling the equality function.
func (source enumerable[TSource]) ContainsAllFunc(equal Equal[TSource], values ...TSource) (result bool) {
	for _, value := range values {
		if !source.ContainsFunc(value, equal) {
			return false
		}
	}
	return true
}

// Counts the elements in a sequence, optionally counting only the elements that satisfy a condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result int
//
// The number of elements in the sequence, or the number of elements that satisfy the predicate when one is supplied.
func (source enumerable[TSource]) Count(predicate ...Predicate[TSource]) (result int) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				result++
			}
		}
	} else {
		for range source {
			result++
		}
	}
	return result
}

// Returns distinct elements from a sequence by using the equality operator to compare values.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains distinct elements from the source sequence.
//
// # Remarks
//
// Elements are compared using Go's == operator.
// Only the first occurrence of each distinct value is returned.
func (source enumerable[TSource]) Distinct[T comparable]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			distinct := map[T]struct{}{}
			for item := range source {
				value := (any(item)).(T)
				if _, ok := distinct[value]; !ok {
					distinct[value] = struct{}{}
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// Returns distinct elements from a sequence by using the Equal method to compare values.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains distinct elements from the source sequence.
//
// # Remarks
//
// Elements are compared by calling the Equal method.
// Only the first occurrence of each distinct value is returned.
func (source enumerable[TSource]) DistinctEquatable[T iEquatable[T]]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			result := make([]T, 0)
			for item := range source {
				value := (any(item)).(T)
				if !containsEquatable(result, value) {
					result = append(result, value)
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

func containsEquatable[TSource iEquatable[TSource]](source []TSource, value TSource) (result bool) {
	for _, item := range source {
		if item.Equal(value) {
			return true
		}
	}
	return false
}

// Returns distinct elements from a sequence by using a specified equality function to compare values.
//
// # Parameters
//
//	equal Equal[TSource]
//
// A function to compare elements for equality.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains distinct elements from the source sequence.
//
// # Remarks
//
// Elements are compared by calling the equality function.
// Only the first occurrence of each distinct value is returned.
func (source enumerable[TSource]) DistinctFunc(equal Equal[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			result := make([]TSource, 0)
			for item := range source {
				if !containsFunc(result, item, equal) {
					result = append(result, item)
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

func containsFunc[TSource any](source []TSource, value TSource, equal Equal[TSource]) (result bool) {
	for _, item := range source {
		if equal(item, value) {
			return true
		}
	}
	return false
}

// Returns the element at a specified index in a sequence.
//
// # Parameters
//
//	index int
//
// The zero-based index of the element to retrieve.
//
// # Returns
//
//	result TSource
//
// The element at the specified position in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrIndexOutOfRange
//     When index is less than 0 or greater than or equal to the number of elements in source.
func (source enumerable[TSource]) ElementAt(index int) (result TSource, err error) {
	if index < 0 {
		return result, ErrIndexOutOfRange
	}
	for item := range source {
		if index == 0 {
			return item, nil
		}
		index--
	}
	return result, ErrIndexOutOfRange
}

// Returns the element at a specified index in a sequence or a default value if the index is out of range.
//
// # Parameters
//
//	index int
//
// The zero-based index of the element to retrieve, counted from the beginning of the sequence.
//
// # Returns
//
//	result TSource
//
// The zero value of TSource if index is outside the bounds of the source sequence; otherwise, the element at the specified position.
func (source enumerable[TSource]) ElementAtOrDefault(index int) (result TSource) {
	if index < 0 {
		return result
	}
	for item := range source {
		if index == 0 {
			return item
		}
		index--
	}
	return result
}

// Returns the element at a specified index in a sequence or a fallback value if the index is out of range.
//
// # Parameters
//
//	index int
//
// The zero-based index of the element to retrieve, counted from the beginning of the sequence.
//
//	fallback TSource
//
// The value to return if index is outside the bounds of the source sequence.
//
// # Returns
//
//	result TSource
//
// The fallback value if index is outside the bounds of the source sequence; otherwise, the element at the specified position.
//
// # Remarks
//
// No error is returned when index is negative or greater than or equal to the number of elements in the source sequence.
func (source enumerable[TSource]) ElementAtOrFallback(index int, fallback TSource) (result TSource) {
	if index < 0 {
		return fallback
	}
	for item := range source {
		if index == 0 {
			return item
		}
		index--
	}
	return fallback
}

// Returns a pointer to the element at a specified index in a sequence or nil if the index is out of range.
//
// # Parameters
//
//	index int
//
// The zero-based index of the element to retrieve, counted from the beginning of the sequence.
//
// # Returns
//
//	result *TSource
//
// A pointer to the element at the specified position, or nil if index is outside the bounds of the source sequence.
func (source enumerable[TSource]) ElementAtOrNil(index int) (result *TSource) {
	if index < 0 {
		return nil
	}
	for item := range source {
		if index == 0 {
			return &item
		}
		index--
	}
	return nil
}

// Produces the difference of two sequences by using a specified equality function to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose elements are excluded from the result when they match elements in the source sequence.
//
//	equal Equal[TSource]
//
// A function to compare elements for equality.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that do not match any element in other.
//
// # Remarks
//
// Elements are compared by calling the equality function.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) ExceptFunc(other Enumerable[TSource], equal Equal[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.DistinctFunc(equal).ToSlice()
		outer:
			for item := range source {
				for _, element := range elements {
					if equal(item, element) {
						continue outer
					}
				}
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Produces the difference of two sequences by using the equality operator to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose elements are excluded from the result when they match elements in the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that do not match any element in other.
//
// # Remarks
//
// Elements are compared using Go's == operator.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) Except[T comparable](other Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.Distinct[T]().Cast[T]().ToSlice()
		outer:
			for item := range source {
				value := (any(item)).(T)
				for _, element := range elements {
					if element == value {
						continue outer
					}
				}
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Produces the difference of two sequences by using the Equal method to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose elements are excluded from the result when they match elements in the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that do not match any element in other.
//
// # Remarks
//
// Elements are compared by calling the Equal method.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) ExceptEquatable[T iEquatable[T]](other Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.DistinctEquatable[T]().Cast[T]().ToSlice()
		outer:
			for item := range source {
				value := (any(item)).(T)
				for _, element := range elements {
					if value.Equal(element) {
						continue outer
					}
				}
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Returns the first element of a sequence, or the first element that satisfies a specified condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The first element in the sequence, or the first element that passes the test in the specified predicate function.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
//   - linq.ErrNoElementSatisfiesTheConditionInPredicate
//     When the source sequence contains elements but none of them passes the test in the specified predicate function.
func (source enumerable[TSource]) First(predicate ...Predicate[TSource]) (result TSource, err error) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		found := false
		for item := range source {
			found = true
			if Predicate(item) {
				return item, nil
			}
		}
		if !found {
			return *new(TSource), ErrSourceContainsNoElements
		}
		return *new(TSource), ErrNoElementSatisfiesTheConditionInPredicate
	}
	for item := range source {
		return item, nil
	}
	return *new(TSource), ErrSourceContainsNoElements
}

// Returns the first element of a sequence, or the first element that satisfies a condition, or the zero value if no such element is found.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The zero value of TSource if the source is empty or if no element passes the test specified by predicate; otherwise, the first matching element.
func (source enumerable[TSource]) FirstOrDefault(predicate ...Predicate[TSource]) (result TSource) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				return item
			}
		}
	} else {
		for item := range source {
			return item
		}
	}
	return *new(TSource)
}

// Returns the first element of a sequence, or the first element that satisfies a condition, or a fallback value if no such element is found.
//
// # Parameters
//
//	fallback TSource
//
// The value to return if no matching element is found.
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The fallback value if the source is empty or if no element passes the test specified by predicate; otherwise, the first matching element.
func (source enumerable[TSource]) FirstOrFallback(fallback TSource, predicate ...Predicate[TSource]) (result TSource) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				return item
			}
		}
	} else {
		for item := range source {
			return item
		}
	}
	return fallback
}

// Returns a pointer to the first element of a sequence, or a pointer to the first element that satisfies a condition, or nil if no such element is found.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result *TSource
//
// A pointer to the first matching element, or nil if the source is empty or if no element passes the test specified by predicate.
func (source enumerable[TSource]) FirstOrNil(predicate ...Predicate[TSource]) (result *TSource) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				return &item
			}
		}
	} else {
		for item := range source {
			return &item
		}
	}
	return nil
}

// Performs the specified action on each element of a sequence.
//
// # Parameters
//
//	action Action[TSource]
//
// An action to perform on each element.
func (source enumerable[TSource]) ForEach(action Action[TSource]) {
	for item := range source {
		action(item)
	}
}

// Groups the elements of a sequence according to a specified key selector function.
//
// # Parameters
//
//	key Key[TSource, TKey]
//
// A function that extracts the grouping key from each element.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] that contains one grouping for each distinct key.
// Each grouping contains the source elements associated with its key.
//
// # Remarks
//
// Elements within each grouping retain their order from the source sequence.
// The order of the groupings is not guaranteed.
func (source enumerable[TSource]) GroupBy[TKey comparable, TResult Grouping[TKey, TSource]](key Key[TSource, TKey]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			groups := make(map[TKey][]TSource)
			for item := range source {
				key := key(item)
				groups[key] = append(groups[key], item)
			}
			for key, group := range groups {
				if !yield(TResult{Key: key, enumerable: func(yield func(value TSource) bool) {
					for _, item := range group {
						if !yield(item) {
							return
						}
					}
				}}) {
					return
				}
			}
		},
	}
}

// Produces the intersection of two sequences by using the equality operator to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose matching elements are included in the result.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that match an element in other.
//
// # Remarks
//
// Elements are compared using Go's == operator.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) Intersect[TValue comparable](other Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.Distinct[TValue]().Cast[TValue]().ToSlice()
			for item := range source {
				value := (any(item)).(TValue)
				for _, element := range elements {
					if value == element {
						if !yield(item) {
							return
						}
					}
				}
			}
		},
	}
}

// Produces the intersection of two sequences by using the Equal method to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose matching elements are included in the result.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that match an element in other.
//
// # Remarks
//
// Elements are compared by calling the Equal method.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) IntersectEquatable[TValue iEquatable[TValue]](other Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.DistinctEquatable[TValue]().Cast[TValue]().ToSlice()
			for item := range source {
				value := (any(item)).(TValue)
				for _, element := range elements {
					if value.Equal(element) {
						if !yield(item) {
							return
						}
					}
				}
			}
		},
	}
}

// Produces the intersection of two sequences by using a specified equality function to compare values.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence whose matching elements are included in the result.
//
//	equal Equal[TSource]
//
// A function to compare elements for equality.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that match an element in other.
//
// # Remarks
//
// Elements are compared by calling the equality function.
// The source sequence order is preserved, including duplicate elements.
func (source enumerable[TSource]) IntersectFunc(other Enumerable[TSource], equal Equal[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			elements := other.DistinctFunc(equal).ToSlice()
			for item := range source {
				for _, element := range elements {
					if equal(item, element) {
						if !yield(item) {
							return
						}
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Key[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Key[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, TInner, TResult]
//
// A function to create a result element from each matching pair.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains elements from both sequences that share the same key.
func (source enumerable[TSource]) Join[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey Key[TSource, TKey], innerKey Key[TInner, TKey], join Join[TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := make(map[TKey][]TInner)
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				matches[innerKey] = append(matches[innerKey], innerItem)
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				if innerItems, ok := matches[outerKey]; ok {
					for _, innerItem := range innerItems {
						if !yield(join(outerItem, innerItem)) {
							return
						}
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys using an equatable comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, TInner, TResult]
//
// A function to create a result element from each matching pair.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains elements from both sequences whose keys are equal according to the equatable comparer.
func (source enumerable[TSource]) JoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], join Join[TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TInner]{}
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(innerKey)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, innerItem)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TInner]{Key: innerKey, Value: []TInner{innerItem}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(outerKey)
				})
				if index >= 0 {
					for _, innerItem := range matches[index].Value {
						if !yield(join(outerItem, innerItem)) {
							return
						}
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys using a custom equality comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	equal Equal[TKey]
//
// A function to compare keys for equality.
//
//	join Join[TSource, TInner, TResult]
//
// A function to create a result element from each matching pair.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains elements from both sequences whose keys are equal according to the supplied comparer.
func (source enumerable[TSource]) JoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], equal Equal[TKey], join Join[TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TInner]{}
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, innerKey)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, innerItem)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TInner]{Key: innerKey, Value: []TInner{innerItem}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, outerKey)
				})
				if index >= 0 {
					for _, innerItem := range matches[index].Value {
						if !yield(join(outerItem, innerItem)) {
							return
						}
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the current sequence.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Key[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Key[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, *TInner, TResult]
//
// A function to create a result element from an outer element and a matching inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the current sequence and matching elements from the second sequence, with nil for unmatched inner elements.
func (source enumerable[TSource]) LeftJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey Key[TSource, TKey], innerKey Key[TInner, TKey], join Join[TSource, *TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := map[TKey][]TInner{}
			for item := range inner.enumerable {
				key := innerKey(item)
				matches[key] = append(matches[key], item)
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				if innerItems, ok := matches[outerKey]; ok {
					for _, innerItem := range innerItems {
						if !yield(join(outerItem, &innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(outerItem, nil)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the second sequence.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Key[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Key[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[*TSource, TInner, TResult]
//
// A function to create a result element from a matching outer element and inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the second sequence and matching elements from the current sequence, with nil for unmatched outer elements.
func (source enumerable[TSource]) RightJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey Key[TSource, TKey], innerKey Key[TInner, TKey], join Join[*TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := map[TKey][]TSource{}
			for item := range source {
				key := outerKey(item)
				matches[key] = append(matches[key], item)
			}
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				if outerItems, ok := matches[innerKey]; ok {
					for _, outerItem := range outerItems {
						if !yield(join(&outerItem, innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(nil, innerItem)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the current sequence, using an equatable comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, *TInner, TResult]
//
// A function to create a result element from an outer element and a matching inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the current sequence and matching elements from the second sequence, with nil for unmatched inner elements.
func (source enumerable[TSource]) LeftJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], join Join[TSource, *TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TInner]{}
			for item := range inner.enumerable {
				key := innerKey(item)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(key)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, item)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TInner]{Key: key, Value: []TInner{item}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(outerKey)
				})
				if index >= 0 {
					for _, innerItem := range matches[index].Value {
						if !yield(join(outerItem, &innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(outerItem, nil)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the second sequence, using an equatable comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[*TSource, TInner, TResult]
//
// A function to create a result element from a matching outer element and inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the second sequence and matching elements from the current sequence, with nil for unmatched outer elements.
func (source enumerable[TSource]) RightJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], join Join[*TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TSource]{}
			for item := range source {
				key := outerKey(item)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TSource]) bool {
					return value.Key.Equal(key)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, item)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TSource]{Key: key, Value: []TSource{item}})
				}
			}
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TSource]) bool {
					return value.Key.Equal(innerKey)
				})
				if index >= 0 {
					for _, outerItem := range matches[index].Value {
						if !yield(join(&outerItem, innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(nil, innerItem)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the current sequence, using a custom equality comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	equal Equal[TKey]
//
// A function to compare keys for equality.
//
//	join Join[TSource, *TInner, TResult]
//
// A function to create a result element from an outer element and a matching inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the current sequence and matching elements from the second sequence, with nil for unmatched inner elements.
func (source enumerable[TSource]) LeftJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], equal Equal[TKey], join Join[TSource, *TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TInner]{}
			for item := range inner.enumerable {
				key := innerKey(item)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, key)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, item)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TInner]{Key: key, Value: []TInner{item}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, outerKey)
				})
				if index >= 0 {
					for _, innerItem := range matches[index].Value {
						if !yield(join(outerItem, &innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(outerItem, nil)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys, preserving all elements from the second sequence, using a custom equality comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to join with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	equal Equal[TKey]
//
// A function to compare keys for equality.
//
//	join Join[*TSource, TInner, TResult]
//
// A function to create a result element from a matching outer element and inner element, or nil when no match exists.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence that contains all elements from the second sequence and matching elements from the current sequence, with nil for unmatched outer elements.
func (source enumerable[TSource]) RightJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], equal Equal[TKey], join Join[*TSource, TInner, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := []KeyValuePair[TKey, []TSource]{}
			for item := range source {
				key := outerKey(item)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TSource]) bool {
					return equal(value.Key, key)
				})
				if index >= 0 {
					matches[index].Value = append(matches[index].Value, item)
				} else {
					matches = append(matches, KeyValuePair[TKey, []TSource]{Key: key, Value: []TSource{item}})
				}
			}
			for innerItem := range inner.enumerable {
				innerKey := innerKey(innerItem)
				index := slices.IndexFunc(matches, func(value KeyValuePair[TKey, []TSource]) bool {
					return equal(value.Key, innerKey)
				})
				if index >= 0 {
					for _, outerItem := range matches[index].Value {
						if !yield(join(&outerItem, innerItem)) {
							return
						}
					}
				} else {
					if !yield(join(nil, innerItem)) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys and groups matching elements from the second sequence for each element in the current sequence.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to group with the current sequence.
//
//	outerKey Key[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Key[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, Enumerable[TInner], TResult]
//
// A function to create a result element from an outer element and the matching inner elements.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence of grouped results where each element of the current sequence is paired with all matching elements from the second sequence.
func (source enumerable[TSource]) GroupJoin[TInner any, TKey comparable, TResult any](inner Enumerable[TInner], outerKey Key[TSource, TKey], innerKey Key[TInner, TKey], join Join[TSource, Enumerable[TInner], TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			matches := map[TKey][]TInner{}
			for item := range inner.enumerable {
				key := innerKey(item)
				matches[key] = append(matches[key], item)
			}
			for outerItem := range source {
				key := outerKey(outerItem)
				if !yield(join(outerItem, Enumerable[TInner]{
					enumerable: func(yield func(value TInner) bool) {
						for _, item := range matches[key] {
							if !yield(item) {
								return
							}
						}
					},
				})) {
					return
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys and groups matching elements from the second sequence for each element in the current sequence, using an equatable comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to group with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, Enumerable[TInner], TResult]
//
// A function to create a result element from an outer element and the matching inner elements.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence of grouped results where each element of the current sequence is paired with all matching elements from the second sequence.
func (source enumerable[TSource]) GroupJoinEquatable[TInner any, TKey iEquatable[TKey], TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], join Join[TSource, Enumerable[TInner], TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			innerByKey := []KeyValuePair[TKey, []TInner]{}
			for item := range inner.enumerable {
				key := innerKey(item)
				index := slices.IndexFunc(innerByKey, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(key)
				})
				if index >= 0 {
					innerByKey[index].Value = append(innerByKey[index].Value, item)
				} else {
					innerByKey = append(innerByKey, KeyValuePair[TKey, []TInner]{Key: key, Value: []TInner{item}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				if index := slices.IndexFunc(innerByKey, func(value KeyValuePair[TKey, []TInner]) bool {
					return value.Key.Equal(outerKey)
				}); index >= 0 {
					matches := innerByKey[index].Value
					if !yield(join(outerItem, Enumerable[TInner]{
						enumerable: func(yield func(value TInner) bool) {
							for _, item := range matches {
								if !yield(item) {
									return
								}
							}
						},
					})) {
						return
					}
				} else {
					if !yield(join(outerItem, Enumerable[TInner]{
						enumerable: func(yield func(value TInner) bool) {},
					})) {
						return
					}
				}
			}
		},
	}
}

// Correlates the elements of two sequences based on matching keys and groups matching elements from the second sequence for each element in the current sequence, using a custom equality comparer.
//
// # Parameters
//
//	inner Enumerable[TInner]
//
// The sequence to group with the current sequence.
//
//	outerKey Value[TSource, TKey]
//
// A function to extract the join key from each element of the current sequence.
//
//	innerKey Value[TInner, TKey]
//
// A function to extract the join key from each element of the second sequence.
//
//	join Join[TSource, Enumerable[TInner], TResult]
//
// A function to create a result element from an outer element and the matching inner elements.
//
//	equal Equal[TKey]
//
// A function to compare keys for equality.
//
// # Returns
//
//	result Enumerable[TResult]
//
// A sequence of grouped results where each element of the current sequence is paired with all matching elements from the second sequence.
func (source enumerable[TSource]) GroupJoinFunc[TInner any, TKey any, TResult any](inner Enumerable[TInner], outerKey Value[TSource, TKey], innerKey Value[TInner, TKey], join Join[TSource, Enumerable[TInner], TResult], equal Equal[TKey]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			innerByKey := []KeyValuePair[TKey, []TInner]{}
			for item := range inner.enumerable {
				key := innerKey(item)
				index := slices.IndexFunc(innerByKey, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, key)
				})
				if index >= 0 {
					innerByKey[index].Value = append(innerByKey[index].Value, item)
				} else {
					innerByKey = append(innerByKey, KeyValuePair[TKey, []TInner]{Key: key, Value: []TInner{item}})
				}
			}
			for outerItem := range source {
				outerKey := outerKey(outerItem)
				if index := slices.IndexFunc(innerByKey, func(value KeyValuePair[TKey, []TInner]) bool {
					return equal(value.Key, outerKey)
				}); index >= 0 {
					matches := innerByKey[index].Value
					if !yield(join(outerItem, Enumerable[TInner]{
						enumerable: func(yield func(value TInner) bool) {
							for _, item := range matches {
								if !yield(item) {
									return
								}
							}
						},
					})) {
						return
					}
				} else {
					if !yield(join(outerItem, Enumerable[TInner]{
						enumerable: func(yield func(value TInner) bool) {},
					})) {
						return
					}
				}
			}
		},
	}
}

// Returns the last element of a sequence, or the last element that satisfies a specified condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The last element in the sequence, or the last element that passes the test in the specified predicate function.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
//   - linq.ErrNoElementSatisfiesTheConditionInPredicate
//     When the source sequence contains elements but none of them passes the test in the specified predicate function.
func (source enumerable[TSource]) Last(predicate ...Predicate[TSource]) (result TSource, err error) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		anyItem := false
		Predicate := predicate[0]
		for item := range source {
			anyItem = true
			if Predicate(item) {
				found = true
				result = item
			}
		}
		if found {
			return result, nil
		}
		if !anyItem {
			return *new(TSource), ErrSourceContainsNoElements
		}
		return *new(TSource), ErrNoElementSatisfiesTheConditionInPredicate
	}
	for item := range source {
		found = true
		result = item
	}
	if found {
		return result, nil
	}
	return *new(TSource), ErrSourceContainsNoElements
}

// Returns the last element of a sequence, or the last element that satisfies a condition, or the zero value if no such element is found.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The zero value of TSource if the source is empty or if no element passes the test specified by predicate; otherwise, the last matching element.
func (source enumerable[TSource]) LastOrDefault(predicate ...Predicate[TSource]) (result TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				found = true
				result = item
			}
		}
		if found {
			return result
		}
		return *new(TSource)
	}
	for item := range source {
		found = true
		result = item
	}
	if found {
		return result
	}
	return *new(TSource)
}

// Returns the last element of a sequence, or the last element that satisfies a condition, or a fallback value if no such element is found.
//
// # Parameters
//
//	fallback TSource
//
// The fallback value to return if no matching element is found.
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The fallback value if the source is empty or if no element passes the test specified by predicate; otherwise, the last matching element.
func (source enumerable[TSource]) LastOrFallback(fallback TSource, predicate ...Predicate[TSource]) (result TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				found = true
				result = item
			}
		}
		if found {
			return result
		}
		return fallback
	}
	for item := range source {
		found = true
		result = item
	}
	if found {
		return result
	}
	return fallback
}

// Returns a pointer to the last element of a sequence, or a pointer to the last element that satisfies a condition, or nil if no such element is found.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result *TSource
//
// A pointer to the last matching element, or nil if the source is empty or if no element passes the test specified by predicate.
func (source enumerable[TSource]) LastOrNil(predicate ...Predicate[TSource]) (result *TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				found = true
				result = &item
			}
		}
		if found {
			return result
		}
		return nil
	}
	for item := range source {
		found = true
		result = &item
	}
	if found {
		return result
	}
	return nil
}

// Returns the maximum value in a sequence by using the comparison operators.
//
// # Returns
//
//	max TResult
//
// The maximum value in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) Max[TResult cmp.Ordered]() (max TResult, err error) {
	found := false
	for item := range source {
		value := (any(item)).(TResult)
		if !found {
			max = value
			found = true
			continue
		}
		if value > max {
			max = value
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the maximum element in a sequence by using the Compare method.
//
// # Returns
//
//	max TSource
//
// The maximum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MaxComparable[TResult iComparable[TSource]]() (max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			max = item
			found = true
			continue
		}
		if (any(max)).(iComparable[TSource]).Compare(item) < 0 {
			max = item
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the maximum element in a sequence by using a specified comparison function.
//
// # Parameters
//
//	compare Compare[TSource]
//
// A function to compare two elements.
//
// # Returns
//
//	max TSource
//
// The maximum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MaxFunc(compare Compare[TSource]) (max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			max = item
			found = true
			continue
		}
		if compare(item, max) > 0 {
			max = item
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the element associated with the maximum selected value in a sequence.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MaxBy[TResult cmp.Ordered](valueSelector func(TSource) TResult) (max TSource, err error) {
	found := false
	var valueMax TResult
	for item := range source {
		if !found {
			max = item
			valueMax = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if value > valueMax {
			max = item
			valueMax = value
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the element associated with the maximum selected value in a sequence by using the Compare method.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MaxByComparable[TResult iComparable[TResult]](valueSelector func(TSource) TResult) (max TSource, err error) {
	found := false
	var valueMax TResult
	for item := range source {
		if !found {
			max = item
			valueMax = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if value.Compare(valueMax) > 0 {
			max = item
			valueMax = value
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the element associated with the maximum selected value in a sequence by using a specified comparison function.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
//	compare Compare[TResult]
//
// A function to compare the selected values.
//
// # Returns
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MaxByFunc[TResult any](valueSelector func(TSource) TResult, compare Compare[TResult]) (max TSource, err error) {
	found := false
	var valueMax TResult
	for item := range source {
		if !found {
			max = item
			valueMax = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if compare(value, valueMax) > 0 {
			max = item
			valueMax = value
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

// Returns the minimum value in a sequence by using the comparison operators.
//
// # Returns
//
//	min TResult
//
// The minimum value in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) Min[TResult cmp.Ordered]() (min TResult, err error) {
	found := false
	for item := range source {
		value := (any(item)).(TResult)
		if !found {
			min = value
			found = true
			continue
		}
		if value < min {
			min = value
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the minimum element in a sequence by using the Compare method.
//
// # Returns
//
//	min TSource
//
// The minimum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinComparable[TResult iComparable[TSource]]() (min TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			found = true
			continue
		}
		if (any(min)).(iComparable[TSource]).Compare(item) > 0 {
			min = item
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the minimum element in a sequence by using a specified comparison function.
//
// # Parameters
//
//	compare Compare[TSource]
//
// A function to compare two elements.
//
// # Returns
//
//	min TSource
//
// The minimum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinFunc(compare Compare[TSource]) (min TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			found = true
			continue
		}
		if compare(item, min) < 0 {
			min = item
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the element associated with the minimum selected value in a sequence.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinBy[TResult cmp.Ordered](valueSelector func(TSource) TResult) (min TSource, err error) {
	found := false
	var valueMin TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if value < valueMin {
			min = item
			valueMin = value
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the element associated with the minimum selected value in a sequence by using the Compare method.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinByComparable[TResult iComparable[TResult]](valueSelector func(TSource) TResult) (min TSource, err error) {
	found := false
	var valueMin TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if value.Compare(valueMin) < 0 {
			min = item
			valueMin = value
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the element associated with the minimum selected value in a sequence by using a specified comparison function.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
//	compare Compare[TResult]
//
// A function to compare the selected values.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinByFunc[TResult any](valueSelector func(TSource) TResult, compare Compare[TResult]) (min TSource, err error) {
	found := false
	var valueMin TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			found = true
			continue
		}
		value := valueSelector(item)
		if compare(value, valueMin) < 0 {
			min = item
			valueMin = value
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

// Returns the minimum and maximum values in a sequence by using the comparison operators.
//
// # Returns
//
//	min TResult
//
// The minimum value in the source sequence.
//
//	max TResult
//
// The maximum value in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMax[TResult cmp.Ordered]() (min TResult, max TResult, err error) {
	found := false
	for item := range source {
		value := (any(item)).(TResult)
		if !found {
			min = value
			max = value
			found = true
			continue
		}
		if value < min {
			min = value
		} else if value > max {
			max = value
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// Returns the minimum and maximum elements in a sequence by using the Compare method.
//
// # Returns
//
//	min TSource
//
// The minimum element in the source sequence.
//
//	max TSource
//
// The maximum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMaxComparable[TResult iComparable[TSource]]() (min TSource, max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			max = item
			found = true
			continue
		}
		if (any(min)).(iComparable[TSource]).Compare(item) > 0 {
			min = item
		} else if (any(max)).(iComparable[TSource]).Compare(item) < 0 {
			max = item
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// Returns the minimum and maximum elements in a sequence by using a specified comparison function.
//
// # Parameters
//
//	compare Compare[TSource]
//
// A function to compare two elements.
//
// # Returns
//
//	min TSource
//
// The minimum element in the source sequence.
//
//	max TSource
//
// The maximum element in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMaxFunc(compare Compare[TSource]) (min TSource, max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			max = item
			found = true
			continue
		}
		if compare(item, min) < 0 {
			min = item
		} else if compare(item, max) > 0 {
			max = item
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// Returns the elements associated with the minimum and maximum selected values in a sequence.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMaxBy[TResult cmp.Ordered](valueSelector func(TSource) TResult) (min TSource, max TSource, err error) {
	found := false
	var valueMin TResult
	var valueMax TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			valueMax = valueMin
			found = true
			continue
		}
		value := valueSelector(item)
		if value < valueMin {
			min = item
			valueMin = value
		} else if value > valueMax {
			max = item
			valueMax = value
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// Returns the elements associated with the minimum and maximum selected values in a sequence by using the Compare method.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMaxByComparable[TResult iComparable[TResult]](valueSelector func(TSource) TResult) (min TSource, max TSource, err error) {
	found := false
	var valueMin TResult
	var valueMax TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			valueMax = valueMin
			found = true
			continue
		}
		value := valueSelector(item)
		if value.Compare(valueMin) < 0 {
			min = item
			valueMin = value
		} else if value.Compare(valueMax) > 0 {
			max = item
			valueMax = value
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// Returns the elements associated with the minimum and maximum selected values in a sequence by using a specified comparison function.
//
// # Parameters
//
//	valueSelector func(TSource) TResult
//
// A function that selects the value used for comparison from each element.
//
//	compare Compare[TResult]
//
// A function to compare the selected values.
//
// # Returns
//
//	min TSource
//
// The element whose selected value is the minimum in the source sequence.
//
//	max TSource
//
// The element whose selected value is the maximum in the source sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence contains no elements.
func (source enumerable[TSource]) MinMaxByFunc[TResult any](valueSelector func(TSource) TResult, compare Compare[TResult]) (min TSource, max TSource, err error) {
	found := false
	var valueMin TResult
	var valueMax TResult
	for item := range source {
		if !found {
			min = item
			valueMin = valueSelector(item)
			valueMax = valueMin
			found = true
			continue
		}
		value := valueSelector(item)
		if compare(value, valueMin) < 0 {
			min = item
			valueMin = value
		} else if compare(value, valueMax) > 0 {
			max = item
			valueMax = value
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

// OfType filters the elements of a sequence based on a specified type.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] that contains elements from the input sequence of type TResult.
func (source enumerable[TSource]) OfType[TResult any]() (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			for item := range source {
				if v, ok := any(item).(TResult); ok {
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// Sorts the elements of a sequence in ascending order by using the comparison operators.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in ascending order.
func (source enumerable[TSource]) Order[T cmp.Ordered]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			sort.Slice(source, func(i, j int) bool {
				return (any(source[i])).(T) < (any(source[j])).(T)
			})
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in ascending order by using the Compare method.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in ascending order.
func (source enumerable[TSource]) OrderComparable[T iComparable[T]]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			slices.SortFunc(source, func(x, y TSource) int {
				return (any(x)).(T).Compare((any(y)).(T))
			})
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in ascending order by using a specified comparison function.
//
// # Parameters
//
//	compare Compare[TSource]
//
// A function to compare two elements.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in ascending order.
func (source enumerable[TSource]) OrderFunc(compare Compare[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			slices.SortFunc(source, compare)
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in ascending order according to a specified key selector function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in ascending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderBy[TValue cmp.Ordered](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return cmp.Compare(value(x), value(y))
			}},
		},
	}
}

// Sorts the elements of a sequence in ascending order according to a specified key selector function by using the Compare method.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in ascending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderByComparable[TValue iComparable[TValue]](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return value(x).Compare(value(y))
			}},
		},
	}
}

// Sorts the elements of a sequence in ascending order according to a specified key selector function and comparison function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
//	compare Compare[TValue]
//
// A function to compare the selected values.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in ascending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderByFunc[TValue any](value Value[TSource, TValue], compare Compare[TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return compare(value(x), value(y))
			}},
		},
	}
}

// Sorts the elements of a sequence in descending order by using the comparison operators.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in descending order.
func (source enumerable[TSource]) OrderDescending[T cmp.Ordered]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			sort.Slice(source, func(i, j int) bool {
				return (any(source[i])).(T) > (any(source[j])).(T)
			})
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in descending order by using the Compare method.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in descending order.
func (source enumerable[TSource]) OrderDescendingComparable[T iComparable[T]]() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			slices.SortFunc(source, func(x, y TSource) int {
				return (any(y)).(T).Compare((any(x)).(T))
			})
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in descending order by using a specified comparison function.
//
// # Parameters
//
//	compare Compare[TSource]
//
// A function to compare two elements.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements are sorted in descending order.
func (source enumerable[TSource]) OrderDescendingFunc(compare Compare[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			source := source.ToSlice()
			slices.SortFunc(source, func(x, y TSource) int {
				return compare(y, x)
			})
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Sorts the elements of a sequence in descending order according to a specified key selector function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in descending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderDescendingBy[TValue cmp.Ordered](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return cmp.Compare(value(y), value(x))
			}},
		},
	}
}

// Sorts the elements of a sequence in descending order according to a specified key selector function by using the Compare method.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in descending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderDescendingByComparable[TValue iComparable[TValue]](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return value(y).Compare(value(x))
			}},
		},
	}
}

// Sorts the elements of a sequence in descending order according to a specified key selector function and comparison function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to sort each element.
//
//	compare Compare[TValue]
//
// A function to compare the selected values.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] whose elements are sorted in descending order by the selected values.
// The result can be further sorted using ThenBy methods.
func (source enumerable[TSource]) OrderDescendingByFunc[TValue any](value Value[TSource, TValue], compare Compare[TValue]) (result OrderedEnumerable[TSource]) {
	return OrderedEnumerable[TSource]{
		&orderedEnumerable[TSource]{
			enumerable: func(yield func(value TSource) bool) {
				source := source.ToSlice()
				slices.SortStableFunc(source, result.comparators.Compare)
				for _, item := range source {
					if !yield(item) {
						return
					}
				}
			},
			comparators: []Compare[TSource]{func(x, y TSource) int {
				return compare(value(y), value(x))
			}},
		},
	}
}

// Prepends values to the beginning of the sequence.
//
// # Parameters
//
//	other ...TSource
//
// The values to prepend to the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// A new Enumerable[TSource] that begins with other and is followed by the elements of the source sequence.
func (source enumerable[TSource]) Prepend(other ...TSource) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(TSource) bool) {
			for _, item := range other {
				if !yield(item) {
					return
				}
			}
			for item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Inverts the order of the elements in a sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] whose elements correspond to those of the input sequence in reverse order.
func (source enumerable[TSource]) Reverse() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for _, item := range slices.Backward(source.ToSlice()) {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Projects each element of a sequence into a new form by using a specified selector function.
//
// # Parameters
//
//	value Value[TSource, TResult]
//
// A function that transforms each source element into a result element.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] whose elements are the transformed elements of the source sequence.
//
// # Remarks
//
// Elements are transformed in source sequence order.
func (source enumerable[TSource]) Select[TResult any](value Value[TSource, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			for item := range source {
				if !yield(value(item)) {
					return
				}
			}
		},
	}
}

// Projects each element of a sequence into a sequence and flattens the resulting sequences into one sequence.
//
// # Parameters
//
//	value Value[TSource, []TResult]
//
// A function that transforms each source element into a slice of result elements.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] that contains the concatenated result elements produced for each source element.
//
// # Remarks
//
// Result elements are returned in source sequence order, with each inner slice processed in its original order.
// An empty inner slice contributes no elements to the result.
func (source enumerable[TSource]) SelectMany[TResult any](value Value[TSource, []TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			for item := range source {
				innerSource := value(item)
				for _, innerItem := range innerSource {
					if !yield(innerItem) {
						return
					}
				}
			}
		},
	}
}

// Determines whether two sequences are equal by comparing corresponding elements using the equality operator.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence to compare with the source sequence.
//
// # Returns
//
//	result bool
//
// True if the sequences have the same length and their corresponding elements are equal; otherwise, false.
func (source enumerable[TSource]) SequenceEqual[TValue comparable](other Enumerable[TSource]) (result bool) {
	next, stop := iter.Pull(iter.Seq[TSource](other.enumerable))
	defer stop()
	for item1 := range source {
		item2, ok := next()
		if !ok || (any(item1)).(TValue) != (any(item2)).(TValue) {
			return false
		}
	}
	_, ok := next()
	return !ok
}

// Determines whether two sequences are equal by comparing corresponding elements using the Equal method.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence to compare with the source sequence.
//
// # Returns
//
//	result bool
//
// True if the sequences have the same length and their corresponding elements are equal; otherwise, false.
func (source enumerable[TSource]) SequenceEqualEquatable[TValue iEquatable[TValue]](other Enumerable[TSource]) (result bool) {
	next, stop := iter.Pull(iter.Seq[TSource](other.enumerable))
	defer stop()
	for item1 := range source {
		item2, ok := next()
		if !ok {
			return false
		}
		if v, ok := (any(item1)).(TValue); ok {
			if !v.Equal((any(item2)).(TValue)) {
				return false
			}
		}
	}
	_, ok := next()
	return !ok
}

// Determines whether two sequences are equal by comparing corresponding elements using a specified equality function.
//
// # Parameters
//
//	other Enumerable[TSource]
//
// A sequence to compare with the source sequence.
//
//	equal Equal[TSource]
//
// A function to compare corresponding elements for equality.
//
// # Returns
//
//	result bool
//
// True if the sequences have the same length and the equality function returns true for every pair of corresponding elements; otherwise, false.
func (source enumerable[TSource]) SequenceEqualFunc(other Enumerable[TSource], equal Equal[TSource]) (result bool) {
	next, stop := iter.Pull(iter.Seq[TSource](other.enumerable))
	defer stop()
	for item1 := range source {
		item2, ok := next()
		if !ok || !equal(item1, item2) {
			return false
		}
	}
	_, ok := next()
	return !ok
}

// Shuffles the order of the elements of a sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// A new sequence with the elements of the source sequence in random order.
func (source enumerable[TSource]) Shuffle() (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			items := source.ToSlice()
			rand.Shuffle(len(items), func(i, j int) {
				items[i], items[j] = items[j], items[i]
			})
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Returns the only element of a sequence, or the only element that satisfies a specified condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single element of the source sequence, or the single element that satisfies the predicate.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrSourceContainsNoElements
//     When the source sequence is empty.
//   - linq.ErrSourceHasMoreThanOneElement
//     When the source sequence has more than one element and no predicate is supplied.
//   - linq.ErrNoElementSatisfiesTheConditionInPredicate
//     When the source sequence contains no element that satisfies the supplied predicate.
//   - linq.ErrMoreThanOneElementSatisfiesTheConditionInPredicate
//     When the source sequence contains more than one element that satisfies the supplied predicate.
func (source enumerable[TSource]) Single(predicate ...Predicate[TSource]) (result TSource, err error) {
	found := false
	anyItem := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			anyItem = true
			if Predicate(item) {
				if found {
					return *new(TSource), ErrMoreThanOneElementSatisfiesTheConditionInPredicate
				}
				result = item
				found = true
			}
		}
		if found {
			return result, nil
		}
		if !anyItem {
			return *new(TSource), ErrSourceContainsNoElements
		}
		return *new(TSource), ErrNoElementSatisfiesTheConditionInPredicate
	}
	for item := range source {
		if found {
			return *new(TSource), ErrSourceHasMoreThanOneElement
		}
		result = item
		found = true
	}
	if !found {
		return *new(TSource), ErrSourceContainsNoElements
	}
	return result, nil
}

// Returns the only element of a sequence, or the only element that satisfies a specified condition, or the zero value if none or more than one such element exists.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single matching element, or the zero value of TSource if no matching element exists or more than one matching element exists.
func (source enumerable[TSource]) SingleOrDefault(predicate ...Predicate[TSource]) (result TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				if found {
					return *new(TSource)
				}
				result = item
				found = true
			}
		}
		return result
	}
	for item := range source {
		if found {
			return *new(TSource)
		}
		result = item
		found = true
	}
	return result
}

// Returns the only element of a sequence, or the only element that satisfies a specified condition, or a fallback value if none or more than one such element exists.
//
// # Parameters
//
//	fallback TSource
//
// The fallback value to return if no matching element exists or more than one matching element exists.
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single matching element, or the fallback value if no matching element exists or more than one matching element exists.
func (source enumerable[TSource]) SingleOrFallback(fallback TSource, predicate ...Predicate[TSource]) (result TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				if found {
					return fallback
				}
				result = item
				found = true
			}
		}
	} else {
		for item := range source {
			if found {
				return fallback
			}
			result = item
			found = true
		}
	}
	if found {
		return result
	}
	return fallback
}

// Returns a pointer to the only element of a sequence, or a pointer to the only element that satisfies a specified condition, or nil if none or more than one such element exists.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result *TSource
//
// A pointer to the single matching element, or nil if no matching element exists or more than one matching element exists.
func (source enumerable[TSource]) SingleOrNil(predicate ...Predicate[TSource]) (result *TSource) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			if Predicate(item) {
				if found {
					return nil
				}
				result = &item
				found = true
			}
		}
		return result
	}
	for item := range source {
		if found {
			return nil
		}
		result = &item
		found = true
	}
	return result
}

// Bypasses a specified number of elements in a sequence and then returns the remaining elements.
//
// # Parameters
//
//	count int
//
// The number of elements to skip before returning the remaining elements.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements that occur after the specified number of elements in the source sequence.
//
// # Remarks
//
// If count is greater than or equal to the number of elements in the source sequence, this method returns an empty sequence. If count is less than or equal to zero, the source sequence is returned unchanged.
func (source enumerable[TSource]) Skip(count int) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if count > 0 {
					count--
					continue
				}
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Returns a sequence with the last count elements omitted.
//
// # Parameters
//
//	count int
//
// The number of elements to omit from the end of the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the source elements except for the last count elements.
func (source enumerable[TSource]) SkipLast(count int) (result Enumerable[TSource]) {
	if count <= 0 {
		return Enumerable[TSource]{
			enumerable: source,
		}
	}
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			buffer := make([]TSource, count)
			i := 0
			n := 0
			for item := range source {
				if n < count {
					buffer[i] = item
					i++
					if i == count {
						i = 0
					}
					n++
					continue
				}
				if !yield(buffer[i]) {
					return
				}
				buffer[i] = item
				i++
				if i == count {
					i = 0
				}
			}
		},
	}
}

// Bypasses elements in a sequence as long as a specified condition is true and then returns the remaining elements.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence starting at the first element that does not satisfy predicate.
func (source enumerable[TSource]) SkipWhile(predicate Predicate[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			skip := true
			for item := range source {
				if skip && predicate(item) {
					continue
				}
				skip = false
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Computes the sum of the elements in a sequence.
//
// The element type must be a numeric type or a string type.
//
// # Returns
//
//	result T
//
// The sum of the elements in the source sequence, or the zero value of T for an empty sequence.
//
// # Error
//
//	err error
//
// Errors that may be returned include:
//   - linq.ErrTypeIsNotNumberOrString
//     When the source element type is neither numeric nor a string type.
func (source enumerable[TSource]) Sum[T Number | ~string]() (result T, err error) {
	if _, ok := (any(*new(TSource))).(T); !ok {
		return result, ErrTypeIsNotNumberOrString
	}
	for item := range source {
		result += (any(item)).(T)
	}
	return result, nil
}

// Computes the sum of the values produced by applying a selector function to the elements of a sequence.
//
// # Parameters
//
//	value Value[TSource, T]
//
// A function that transforms each source element into a numeric or string value to include in the sum.
//
// # Returns
//
//	result T
//
// The sum of the selected values, or the zero value of T for an empty sequence.
func (source enumerable[TSource]) SumFunc[T Number | ~string](value Value[TSource, T]) (result T) {
	for item := range source {
		result += value(item)
	}
	return result
}

// Returns a specified number of contiguous elements from the start of a sequence.
//
// # Parameters
//
//	count int
//
// The number of elements to return.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains up to count elements from the beginning of the source sequence.
func (source enumerable[TSource]) Take(count int) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if count <= 0 {
					return
				}
				if !yield(item) {
					return
				}
				count--
			}
		},
	}
}

// Returns a sequence that contains the last count elements from the source sequence.
//
// # Parameters
//
//	count int
//
// The number of elements to take from the end of the source sequence.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains up to count elements from the end of the source sequence.
func (source enumerable[TSource]) TakeLast(count int) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			if count <= 0 {
				return
			}
			buffer := make([]TSource, count)
			i := 0
			n := 0
			ovf := false
			for item := range source {
				buffer[i] = item
				i++
				if i == count {
					i = 0
					ovf = true
				}
				if n < count {
					n++
				}
			}
			if !ovf {
				i = 0
			}
			for j := 0; j < n; j++ {
				if !yield(buffer[i]) {
					return
				}
				i++
				if i == count {
					i = 0
				}
			}
		},
	}
}

// Returns elements from a sequence as long as a specified condition is true, then stops at the first element that does not satisfy the condition.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the consecutive source elements that satisfy predicate.
func (source enumerable[TSource]) TakeWhile(predicate Predicate[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if !predicate(item) {
					return
				}
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Writes the elements of a sequence to a channel.
//
// # Parameters
//
//	result chan<- TSource
//
// The channel to which the source elements are written.
//
// # Remarks
//
// The channel is closed when the sequence is fully enumerated.
func (source enumerable[TSource]) ToChannel(result chan<- TSource) {
	defer close(result)
	for item := range source {
		result <- item
	}
}

// Creates a map[TKey]TValue from a sequence according to specified key and value selector functions.
//
// # Parameters
//
//	key Key[TSource, TKey]
//
// A function to extract a key from each element.
//
//	value Value[TSource, TValue]
//
// A function to produce a map value from each element.
//
// # Returns
//
//	result map[TKey]TValue
//
// A map[TKey]TValue that contains values selected from the source sequence.
//
// # Remarks
//
// If multiple source elements produce the same key, the value from the last such element replaces the previous value.
func (source enumerable[TSource]) ToMap[TKey comparable, TValue any](key Key[TSource, TKey], value Value[TSource, TValue]) (result map[TKey]TValue) {
	result = make(map[TKey]TValue)
	for item := range source {
		result[key(item)] = value(item)
	}
	return result
}

// Creates a slice of TSource from a sequence.
//
// # Returns
//
//	result []TSource
//
// A slice of TSource that contains the elements from the source sequence in source order.
// An empty source sequence produces an empty, non-nil slice.
func (source enumerable[TSource]) ToSlice() (result []TSource) {
	result = make([]TSource, 0)
	for item := range source {
		result = append(result, item)
	}
	return result
}

// Produces the union of two sequences by using the equality operator to compare values.
//
// # Parameters
//
//	sequence Enumerable[TSource]
//
// A sequence whose distinct elements are added to the result.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the distinct elements from both input sequences.
func (source enumerable[TSource]) Union[TValue comparable](sequence Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			items := source.Distinct[TValue]()
			for item := range items.enumerable {
				if !yield(item) {
					return
				}
			}
			for item := range sequence.Distinct[TValue]().enumerable {
				if !items.Contains[TValue](item) {
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// Produces the union of two sequences by using the Equal method to compare values.
//
// # Parameters
//
//	sequence Enumerable[TSource]
//
// A sequence whose distinct elements are added to the result.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the distinct elements from both input sequences.
func (source enumerable[TSource]) UnionEquatable[TValue iEquatable[TValue]](sequence Enumerable[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			items := source.DistinctEquatable[TValue]()
			for item := range items.enumerable {
				if !yield(item) {
					return
				}
			}
			for item := range sequence.DistinctEquatable[TValue]().enumerable {
				if !items.ContainsEquatable[TValue](item) {
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// Produces the union of two sequences by using a specified equality function to compare values.
//
// # Parameters
//
//	sequence Enumerable[TSource]
//
// A sequence whose distinct elements are added to the result.
//
//	equal Equal[TSource]
//
// A function to compare elements for equality.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the distinct elements from both input sequences.
func (source enumerable[TSource]) UnionFunc(sequence Enumerable[TSource], equal Equal[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			items := source.DistinctFunc(equal)
			for item := range items.enumerable {
				if !yield(item) {
					return
				}
			}
			for item := range sequence.DistinctFunc(equal).enumerable {
				if !items.ContainsFunc(item, equal) {
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// Filters a sequence of values based on a predicate.
//
// # Parameters
//
//	predicate Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Enumerable[TSource]
//
// An Enumerable[TSource] that contains the elements from the source sequence that satisfy predicate.
func (source enumerable[TSource]) Where(predicate Predicate[TSource]) (result Enumerable[TSource]) {
	return Enumerable[TSource]{
		enumerable: func(yield func(value TSource) bool) {
			for item := range source {
				if predicate(item) {
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// Produces a sequence by combining corresponding elements from two sequences.
//
// # Parameters
//
//	other Enumerable[TSecond]
//
// The second sequence to combine with the source sequence.
//
//	join Join[TFirst, TSecond, TResult]
//
// A function that combines corresponding elements from the two sequences.
//
// # Returns
//
//	result Enumerable[TResult]
//
// An Enumerable[TResult] produced by applying join to corresponding elements from the source and other sequence.
//
// # Remarks
//
// The result contains at most as many elements as the shorter input sequence.
// Elements are combined in sequence order.
func (source enumerable[TFirst]) Zip[TSecond any, TResult any](other Enumerable[TSecond], join Join[TFirst, TSecond, TResult]) (result Enumerable[TResult]) {
	return Enumerable[TResult]{
		enumerable: func(yield func(value TResult) bool) {
			next, stop := iter.Pull(iter.Seq[TSecond](other.enumerable))
			defer stop()
			for item1 := range source {
				item2, ok := next()
				if !ok {
					return
				}
				if !yield(join(item1, item2)) {
					return
				}
			}
		},
	}
}
