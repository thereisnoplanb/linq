package linq

import "github.com/thereisnoplanb/generic"

// Returns the input typed as Iterator[TSource].
//
// # Parameters
//
//	source []TSource
//
// The sequence of TSource.
//
// # Returns
//
//	result Iterator[TSource]
//
// The input sequence typed as Iterator[TSource].
func FromSlice[TSlice ~[]TSource, TSource any](source TSlice) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		for _, value := range source {
			if !yield(value) {
				return
			}
		}
	}
}

// Returns the input typed as Iterator[TSource].
//
// # Parameters
//
//	source Iterator[TSource]
//
// The sequence of TSource.
//
// # Returns
//
//	result Iterator[TSource]
//
// The input sequence typed as Iterator[TSource].
func FromIterator[TSource any](source Iterator[TSource]) Iterator[TSource] {
	return source
}

// Returns the input typed as Iterator[generic.KeyValuePair[TKey, TValue]].
//
// # Parameters
//
//	source map[TKey]TValue
//
// The sequence of generic.KeyValuePair[TKey, TValue].
//
// # Returns
//
//	result Iterator[generic.KeyValuePair[TKey, TValue]]
//
// The input sequence typed as Iterator[generic.KeyValuePair[TKey, TValue]].
func FromMap[TMap ~map[TKey]TValue, TKey comparable, TValue any](source TMap) Iterator[generic.KeyValuePair[TKey, TValue]] {
	return func(yield func(value generic.KeyValuePair[TKey, TValue]) bool) {
		for key, value := range source {
			if !yield(generic.KeyValuePair[TKey, TValue]{
				Key:   key,
				Value: value,
			}) {
				return
			}
		}
	}
}

// Returns the input typed as Iterator[rune].
//
// # Parameters
//
//	source string
//
// The sequence of runes.
//
// # Returns
//
//	result Iterator[rune]
//
// The input sequence typed as Iterator[rune].
func FromString(source string) Iterator[rune] {
	return func(yield func(value rune) bool) {
		for _, value := range source {
			if !yield(value) {
				return
			}
		}
	}
}

// Generates a sequence that contains one repeated value.
//
// # Parameters
//
//	element TSource
//
// The value to be repeated.
//
//	count int
//
// The number of times to repeat the value in the generated sequence.
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains a repeated value.
func Repeat[TSource any](element TSource, count int) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		for count > 0 {
			if !yield(element) {
				return
			}
			count--
		}
	}
}

// Generates a sequence of integral numbers within a specified range.
//
// # Parameters
//
//	start int
//
// The value of the first integer in the sequence.
//
//	count int
//
// The number of sequential integers to generate.
//
// # Returns
//
//	result Iterator[int]
//
// An Iterator[int] that contains a range of sequential integral numbers.
func Range(start int, count int) Iterator[int] {
	return func(yield func(value int) bool) {
		for count > 0 {
			if !yield(start) {
				return
			}
			start++
			count--
		}
	}
}
