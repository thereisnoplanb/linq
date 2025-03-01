package linq

import (
	"cmp"
	"iter"
	"reflect"
	"slices"

	"github.com/thereisnoplanb/generic"
)

type Iterator[TSource any] iter.Seq[TSource]

// Applies an accumulator function over a sequence. The specified seed value is used as the initial accumulator value.
//
// # Parameters
//
//	seed TSource
//
// The initial accumulator value.
//
//	accumulator generic.Accumulator[TSource,TSource]
//
// An accumulator function to be invoked on each element.
//
//	resultSelector func(TSource) TSource
//
// A function to transform the final accumulator value into the result value. [OPTIONAL]
//
// # Returns
//
//	result TSource - The final accumulator value.
func (source Iterator[TSource]) Aggregate(seed TSource, accumulator generic.Accumulator[TSource, TSource], resultSelector ...func(TSource) TSource) (result TSource) {
	result = seed
	for item := range source {
		result = accumulator(result, item)
	}
	if len(resultSelector) > 0 {
		result = resultSelector[0](result)
	}
	return result
}

// Applies an accumulator function over a sequence. The specified seed value is used as the initial accumulator value.
//
// # Parameters
//
//	seed TAccumulator
//
// The initial accumulator value.
//
//	accumulator generic.Accumulator[TSource, TAccumulator]
//
// An accumulator function to be invoked on each element.
//
//	resultSelector func(TAccumulator) TAccumulator
//
// A function to transform the final accumulator value into the result value. [OPTIONAL]
//
// # Returns
//
//	result TAccumulator
//
// The final accumulator value.
func Aggregate[TSource any, TAccumulator any](source Iterator[TSource], seed TAccumulator, accumulator generic.Accumulator[TSource, TAccumulator], resultSelector ...func(TAccumulator) TAccumulator) (result TAccumulator) {
	result = seed
	for item := range source {
		result = accumulator(result, item)
	}
	if len(resultSelector) > 0 {
		result = resultSelector[0](result)
	}
	return result
}

// Determines whether all elements of a sequence satisfy a condition.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result bool
//
// True if every element of the source sequence passes the test in the specified predicate, or if the sequence is empty; otherwise, false.
func (source Iterator[TSource]) All(predicate generic.Predicate[TSource]) (result bool) {
	for item := range source {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Determines whether any element of a sequence satisfies a condition.
//
// # Parameters
//
//	predicate generic.Predicate[TSource] - A function to test each element for a condition.
//
// # Returns
//
//	result bool
//
// True if the source sequence is not empty and at least one of its elements passes the test in the specified predicate; otherwise, false.
func (source Iterator[TSource]) Any(predicate ...generic.Predicate[TSource]) (result bool) {
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
//	elements TSource
//
// The values to append to source.
//
// # Returns
//
//	result Iterator[TSource]
//
// A new sequence that ends with elements.
func (source Iterator[TSource]) Append(elements ...TSource) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if !yield(item) {
				return
			}
		}
		for _, element := range elements {
			if !yield(element) {
				return
			}
		}
	}
}

// Computes the average of a sequence of numeric values.
//
// # Parameters
//
//	source Iterator[TSource]
//
// A sequence of values to calculate the average of.
//
// # Returns
//
//	result float64
//
// The average of the sequence of values.
//
// # Error
//
//	err error
//
// linq.ErrSourceContainsNoElements if source contains no elements.
func Average[TSource generic.Real](source Iterator[TSource]) (result float64, err error) {
	count := 0
	var sum TSource
	for item := range source {
		sum += item
		count++
	}
	if count == 0 {
		return result, ErrSourceContainsNoElements
	}
	return float64(sum) / float64(count), nil
}

// Casts the elements of an Iterator to the specified type.
//
// # Parameters
//
//	source Iterator[TSource]
//
// The Iterator that contains the elements to be cast to type TResult.
//
// # Returns
//
//	result Iterator[TResult]
//
// An Iterator[TResult] that contains each element of the source sequence cast to the specified type.
//
// # Panics
//
// When an element in the sequence cannot be cast to type TResult.
func Cast[TSource any, TResult any](source Iterator[TSource]) (result Iterator[TResult]) {
	return func(yield func(value TResult) bool) {
		for item := range source {
			if !yield((any(item)).(TResult)) {
				return
			}
		}
	}
}

// func (source Iterator[TSource]) Chunk(size int) (result Iterator[Iterator[TSource]]) {
// 	if size < 1 {
// 		panic(ErrSizeIsBelowOne)
// 	}
// 	return func(yield func(value Iterator[TSource]) bool) {
// 		slice := make([]TSource, 0)
// 		for item := range source {
// 			slice = append(slice, item)
// 		}
// 		iterators := make([]Iterator[TSource], 0)
// 		for i := 0; i < len(slice); i += size {
// 			end := min(size, len(slice[i:]))
// 			iterators = append(iterators, FromSlice(slice[i:i+end:i+end]))
// 		}
// 		for _, iterator := range iterators {
// 			if !yield(iterator) {
// 				return
// 			}
// 		}
// 	}
// }

// Concatenates two sequences.
//
// # Parameters
//
//	sequence Iterator[TSource]
//
// The sequence to concatenate to the original sequence.
//
// # Returns
//
//	result ITerator[TSource]
//
// An Iterator[TSource] that contains the concatenated elements of the original sequence and input sequence.
func (source Iterator[TSource]) Concat(sequence Iterator[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if !yield(item) {
				return
			}
		}
		for item := range sequence {
			if !yield(item) {
				return
			}
		}
	}
}

// Determines whether a sequence contains a specified element by using a specified generic.Equality[TSource].
//
// # Parameters
//
//	value TSource
//
// The value to locate in the sequence.
//
//	comparer generic.Equality[TSource]
//
// An equality comparer to compare values. [OPTIONAL]
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
// If the comparer parameter is omitted or nil, the default equality comparator is used to compare elements to the specified value.
// Before doing this, it is checked whether the type TSource implements the generic.IEquatable interface.
// If so, the Equals() method from that interface is used to compare elements to the specified value.
func (source Iterator[TSource]) Contains(value TSource, comparer ...generic.Equality[TSource]) (result bool) {
	if len(comparer) > 0 && comparer[0] != nil {
		Equal := comparer[0]
		for item := range source {
			if Equal(item, value) {
				return true
			}
		}
	} else if v, ok := (any(value)).(generic.IEquatable[TSource]); ok {
		for item := range source {
			if v.Equal(item) {
				return true
			}
		}
	} else if v := reflect.ValueOf(value); v.Comparable() {
		for item := range source {
			if v.Equal(reflect.ValueOf(item)) {
				return true
			}
		}
	} else {
		for item := range source {
			if reflect.DeepEqual(item, value) {
				return true
			}
		}
	}
	return false
}

// Determines whether a sequence contains any of the specified elements by using a specified generic.Equality[TSource].
//
// # Parameters
//
//	values []TSource
//
// The list of values to locate in the sequence.
//
//	comparer generic.Equality[TSource]
//
// An equality comparer to compare values. [OPTIONAL]
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
// If the comparer parameter is omitted or nil, the default equality comparator is used to compare elements to the specified values.
// Before doing this, it is checked whether the type TSource implements the generic.IEquatable interface.
// If so, the Equals() method from that interface is used to compare elements to the specified values.
func (source Iterator[TSource]) ContainsAny(values []TSource, comparer ...generic.Equality[TSource]) (result bool) {
	for _, value := range values {
		if source.Contains(value, comparer...) {
			return true
		}
	}
	return false
}

// Determines whether a sequence contains all specified elements by using a specified generic.Equality[TSource].
//
// # Parameters
//
//	values TSource
//
// The values to locate in the sequence.
//
//	comparer generic.Equality[TSource]
//
// An equality comparer to compare values. [OPTIONAL]
//
// # Returns
//
//	result bool
//
// True if the source sequence contains all elements from the specified values; otherwise, false.
//
// # Remarks
//
// Iteration is terminated as soon as any matching element is not found.
// If the comparer parameter is omitted or nil, the default equality comparator is used to compare elements to the specified values.
// Before doing this, it is checked whether the type TSource implements the generic.IEquatable interface.
// If so, the Equals() method from that interface is used to compare elements to the specified values.
func (source Iterator[TSource]) ContainsAll(values []TSource, comparer ...generic.Equality[TSource]) (result bool) {
	for _, value := range values {
		if !source.Contains(value, comparer...) {
			return false
		}
	}
	return true
}

// Returns the number of elements in a sequence or returns a number that represents how many elements in the specified sequence satisfy a condition in predicate if passed.
//
// # Parameters
//
//	predicate genreic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result int
//
// The number of elements in the input sequence or a number that represents how many elements in the sequence satisfy the condition in the predicate function if passed.
func (source Iterator[TSource]) Count(predicate ...generic.Predicate[TSource]) (result int) {
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

// Returns distinct elements from a sequence by using a specified generic.Equality[TSource] to compare values.
//
// # Parameters
//
//	comparer generic.Equality[TSource]
//
// An generic.Equality[TSource] to compare values. [OPTIONAL]
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains distinct elements from the source sequence.
//
// # Remarks
//
// If the comparer parameter is omitted or nil, the default equality comparator is used to compare elements to the specified value.
// Before doing this, it is checked whether the type TSource implements the generic.IEquatable interface.
// If so, the Equals() method from that interface is used to compare elements to the specified value.
func (source Iterator[TSource]) Distinct(comparer ...generic.Equality[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		result := make([]TSource, 0)
		for item := range source {
			if !FromSlice(result).Contains(item, comparer...) {
				result = append(result, item)
			}
		}
		for _, item := range result {
			if !yield(item) {
				return
			}
		}
	}
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
// ErrIndexOutOfRange when index is less than 0 or greater than or equal to the number of elements in source.
func (source Iterator[TSource]) ElementAt(index int) (result TSource, err error) {
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
// The index of the element to retrieve, which is either from the beginning or the end of the sequence.
//
// # Returns
//
//	reslut TSource
//
// Default value if index is outside the bounds of the source sequence; otherwise, the element at the specified position in the source sequence.
func (source Iterator[TSource]) ElementAtOrDefault(index int) (result TSource) {
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
// The index of the element to retrieve, which is either from the beginning or the end of the sequence.
//
// # Returns
//
//	reslut TSource
//
// Fallback value if index is outside the bounds of the source sequence; otherwise, the element at the specified position in the source sequence.
func (source Iterator[TSource]) ElementAtOrFallback(index int, fallback TSource) (result TSource) {
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

// Produces the set difference of two sequences.
//
// # Parameters
//
//	sequence Iterator[TSource]
//
// An Iterator[TSource] whose distinct elements that also occur in the source sequence will cause those elements to be removed from the returned sequence.
//
//	comparer generic.Equality[TSource]
//
// An Equality function to compare values. [OPTIONAL]
//
// # Returns
//
//	result Iterator[TSource]
//
// A sequence that contains the set difference of the elements of two sequences.
//
// # Remarks
//
// The set difference of two sets is defined as the members of the source set that don't appear in the sequence set.
// This method returns those elements in source that don't appear in sequence.
// It doesn't return those elements in sequence that don't appear in source.
// Only unique elements are returned.
//
// # Example
//
//	source := FromSlice([]int{1, 2, 3, 1, 2, 3})
//	sequence := FroSlice([]int{1, 2, 1, 1})
//	result := source.Except(sequence).ToSlice()
//	/*This code produces the following output result = []int{3}*/
func (source Iterator[TSource]) Except(sequence Iterator[TSource], comparer ...generic.Equality[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		if len(comparer) > 0 {
			isEqual := comparer[0]
			for item := range source.Distinct(isEqual) {
				for other := range sequence.Distinct(isEqual) {
					if isEqual(item, other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		} else if _, ok := (any(*new(TSource))).(generic.IEquatable[TSource]); ok {
			for item := range source.Distinct() {
				c := any(item).(generic.IEquatable[TSource])
				for other := range sequence.Distinct() {
					if c.Equal(other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		} else {
			for item := range source.Distinct() {
				for other := range sequence.Distinct() {
					if reflect.DeepEqual(item, other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		}
	}
}

// Returns the first element of a sequence or returns the first element in a sequence that satisfies a specified condition in predicate if passed.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The first element in the specified sequence or the first element in the sequence that passes the test in the specified predicate function if passed.
//
// # Error
//
//	err error
//
// ErrSourceContainsNoElements - When sequence contains no elelements,
//
// ErrNoElementSatisfiesTheConditionInPredicate - When sequence contains elements but none of them passes the test in the specified predicate function if passed.
func (source Iterator[TSource]) First(predicate ...generic.Predicate[TSource]) (result TSource, err error) {
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		any := false
		for item := range source {
			any = true
			if Predicate(item) {
				return item, nil
			}
		}
		if !any {
			return *new(TSource), ErrSourceContainsNoElements
		}
		return *new(TSource), ErrNoElementSatisfiesTheConditionInPredicate
	}
	for item := range source {
		return item, nil
	}
	return *new(TSource), ErrSourceContainsNoElements
}

// Returns the first element of a sequence, or a default value if the sequence contains no elements or returns the first element of the sequence that satisfies a condition if passed or a default value if no such element is found.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// Default value if source is empty or if no element passes the test specified by predicate; otherwise, the first element in source that passes the test specified by predicate.
func (source Iterator[TSource]) FirstOrDefault(predicate ...generic.Predicate[TSource]) (result TSource) {
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

// Returns the first element of a sequence, or a fallback value if the sequence contains no elements or returns the first element of the sequence that satisfies a condition if passed or a fallback value if no such element is found.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// Fallback value if source is empty or if no element passes the test specified by predicate; otherwise, the first element in source that passes the test specified by predicate.
func (source Iterator[TSource]) FirstOrFallback(fallback TSource, predicate ...generic.Predicate[TSource]) (result TSource) {
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

func GroupBy[TSource any, TKey comparable](source Iterator[TSource], keySelector generic.KeySelector[TSource, TKey]) (result Iterator[generic.KeyValuePair[TKey, Iterator[TSource]]]) {
	return func(yield func(object generic.KeyValuePair[TKey, Iterator[TSource]]) bool) {
		groups := make(map[TKey][]TSource)
		for item := range source {
			key := keySelector(item)
			groups[key] = append(groups[key], item)
		}
		for key, value := range groups {
			if !yield(generic.KeyValuePair[TKey, Iterator[TSource]]{
				Key:   key,
				Value: FromSlice(value),
			}) {
				return
			}
		}
	}
}

// Produces the set intersection of two sequences.
//
// # Parameters
//
//	sequence Iterator[TSource]
//
// An Iterator[TSource] whose distinct elements that also appear in the source sequence will be returned.
//
//	comparer generic.Equality[TSource]
//
// An Equality function to compare values. [OPTIONAL]
//
// # Returns
//
//	result Iterator[TSource]
//
// A sequence that contains the elements that form the set intersection of two sequences.
//
// # Remarks
//
// The set intersection of two sets is defined as the members of the source that also appear in sequence, but no other elements.
// This method returns those elements in source that also appear in sequence.
// It doesn't return those elements in sequence that don't appear in source.
// Only unique elements are returned.
//
// # Example
//
//	source := FromSlice([]int{1, 2, 3, 1, 2, 3})
//	sequence := FroSlice([]int{1, 2, 1, 1})
//	result := source.Except(sequence).ToSlice()
//	/*This code produces the following output result = []int{1, 2}*/
func (source Iterator[TSource]) Intersect(sequence Iterator[TSource], comparer ...generic.Equality[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		if len(comparer) > 0 {
			isEqual := comparer[0]
			for item := range source.Distinct(isEqual) {
				for other := range sequence.Distinct(isEqual) {
					if isEqual(item, other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		} else if _, ok := (any(*new(TSource))).(generic.IEquatable[TSource]); ok {
			for item := range source.Distinct() {
				c := any(item).(generic.IEquatable[TSource])
				for other := range sequence.Distinct() {
					if c.Equal(other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		} else {
			for item := range source.Distinct() {
				for other := range sequence.Distinct() {
					if reflect.DeepEqual(item, other) {
						if !yield(item) {
							return
						}
					}
				}
			}
		}
	}
}

func joinComparer[TOuter any, TInner any, TKey any, TResult any](outer Iterator[TOuter], inner Iterator[TInner], outerKeySelector generic.ValueSelector[TOuter, TKey], innerKeySelector generic.ValueSelector[TInner, TKey], resultSelector func(outer TOuter, inner TInner) TResult, isEqual generic.Equality[TKey]) (result Iterator[TResult]) {
	return func(yield func(value TResult) bool) {
		for outerItem := range outer {
			outerKey := outerKeySelector(outerItem)
			for innerItem := range inner {
				innerKey := innerKeySelector(innerItem)
				if isEqual(outerKey, innerKey) {
					if !yield(resultSelector(outerItem, innerItem)) {
						return
					}
				}
			}
		}
	}
}

func joinEquatable[TOuter any, TInner any, TKey any, TResult any](outer Iterator[TOuter], inner Iterator[TInner], outerKeySelector generic.ValueSelector[TOuter, TKey], innerKeySelector generic.ValueSelector[TInner, TKey], resultSelector func(outer TOuter, inner TInner) TResult) (result Iterator[TResult]) {
	return func(yield func(value TResult) bool) {
		for outerItem := range outer {
			outerKey := any(outerKeySelector(outerItem)).(generic.IEquatable[TKey])
			for innerItem := range inner {
				innerKey := innerKeySelector(innerItem)
				if outerKey.Equal(innerKey) {
					if !yield(resultSelector(outerItem, innerItem)) {
						return
					}
				}
			}
		}
	}
}

func joinComparable[TOuter any, TInner any, TKey any, TResult any](outer Iterator[TOuter], inner Iterator[TInner], outerKeySelector generic.ValueSelector[TOuter, TKey], innerKeySelector generic.ValueSelector[TInner, TKey], resultSelector func(outer TOuter, inner TInner) TResult) (result Iterator[TResult]) {
	return func(yield func(value TResult) bool) {
		for outerItem := range outer {
			outerKey := outerKeySelector(outerItem)
			for innerItem := range inner {
				innerKey := innerKeySelector(innerItem)
				if reflect.DeepEqual(outerKey, innerKey) {
					if !yield(resultSelector(outerItem, innerItem)) {
						return
					}
				}
			}
		}
	}
}

func Join[TOuter any, TInner any, TKey any, TResult any](outer Iterator[TOuter], inner Iterator[TInner], outerKeySelector generic.ValueSelector[TOuter, TKey], innerKeySelector generic.ValueSelector[TInner, TKey], resultSelector func(outer TOuter, inner TInner) TResult, comparer ...generic.Equality[TKey]) (result Iterator[TResult]) {
	if len(comparer) > 0 && comparer[0] != nil {
		return joinComparer(outer, inner, outerKeySelector, innerKeySelector, resultSelector, comparer[0])
	}
	if _, ok := any(*new(TKey)).(generic.IEquatable[TKey]); ok {
		return joinEquatable(outer, inner, outerKeySelector, innerKeySelector, resultSelector)
	}
	return joinComparable(outer, inner, outerKeySelector, innerKeySelector, resultSelector)
}

// Returns the last element of a sequence or returns the last element in a sequence that satisfies a specified condition in predicate if passed.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The last element in the specified sequence or the last element in the sequence that passes the test in the specified predicate function if passed.
//
// # Error
//
//	err error
//
// linq.ErrSourceContainsNoElements - When sequence contains no elelements.
//
// linq.ErrNoElementSatisfiesTheConditionInPredicate - When sequence contains elements but none of them passes the test in the specified predicate function if passed.
func (source Iterator[TSource]) Last(predicate ...generic.Predicate[TSource]) (result TSource, err error) {
	found := false
	if len(predicate) > 0 && predicate[0] != nil {
		any := false
		Predicate := predicate[0]
		for item := range source {
			any = true
			if Predicate(item) {
				found = true
				result = item
			}
		}
		if found {
			return result, nil
		}
		if !any {
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

// Returns the last element of a sequence, or a default value if the sequence contains no elements or returns the last element of the sequence that satisfies a condition in passed predicate or a default value if no such element is found.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// Default value if source is empty or if no element passes the test specified by predicate; otherwise, the last element in source that passes the test specified by predicate.
func (source Iterator[TSource]) LastOrDefault(predicate ...generic.Predicate[TSource]) (result TSource) {
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

// Returns the last element of a sequence, or a fallback value if the sequence contains no elements or returns the last element of the sequence that satisfies a condition if passed or a fallback value if no such element is found.
//
// # Parameters
//
//	fallback TSource
//
// The fallback value to return if the sequence is empty or constains none element that satisfies predicate (if passed).
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// Fallback value if source is empty or if no element passes the test specified by predicate; otherwise, the last element in source or the last element that passes the test specified by predicate (if passed).
func (source Iterator[TSource]) LastOrFallback(fallback TSource, predicate ...generic.Predicate[TSource]) (result TSource) {
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

func minmax[TSource any, T generic.Real | generic.String](source Iterator[TSource]) (resultMin, resultMax TSource, err error) {
	found := false
	var min, max T
	for item := range source {
		value := any(item).(T)
		if !found {
			min = value
			max = value
			found = true
			continue
		}
		if value < min {
			min = value
			continue
		}
		if value > max {
			max = value
		}
	}
	if !found {
		return resultMin, resultMax, ErrSourceContainsNoElements
	}
	return any(min).(TSource), any(max).(TSource), nil
}

func (source Iterator[TSource]) Max(compare ...generic.Comparison[TSource]) (max TSource, err error) {
	found := false
	if len(compare) > 0 {
		for item := range source {
			if !found {
				max = item
				found = true
				continue
			}
			if compare[0](max, item) > 0 {
				max = item
			}
		}
		if !found {
			return max, ErrSourceContainsNoElements
		}
		return max, nil
	}
	if v, ok := (any(max)).(generic.IComparable[TSource]); ok {
		for item := range source {
			if !found {
				max = item
				found = true
				continue
			}
			if v.Compare(item) > 0 {
				max = item
			}
		}
		if !found {
			return max, ErrSourceContainsNoElements
		}
		return max, nil
	}
	switch (any(*new(TSource))).(type) {
	case int:
		_, max, err = minmax[TSource, int](source)
	case int8:
		_, max, err = minmax[TSource, int8](source)
	case int16:
		_, max, err = minmax[TSource, int16](source)
	case int32:
		_, max, err = minmax[TSource, int32](source)
	case int64:
		_, max, err = minmax[TSource, int64](source)
	case uint:
		_, max, err = minmax[TSource, uint](source)
	case uint8:
		_, max, err = minmax[TSource, uint8](source)
	case uint16:
		_, max, err = minmax[TSource, uint16](source)
	case uint32:
		_, max, err = minmax[TSource, uint32](source)
	case uint64:
		_, max, err = minmax[TSource, uint64](source)
	case uintptr:
		_, max, err = minmax[TSource, uintptr](source)
	case float32:
		_, max, err = minmax[TSource, float32](source)
	case float64:
		_, max, err = minmax[TSource, float64](source)
	case string:
		_, max, err = minmax[TSource, string](source)
	default:
		panic("unsupported type for Max")
	}
	return max, err
}

func Max[TSource generic.Comparable](source Iterator[TSource]) (max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			max = item
			found = true
			continue
		}
		if item > max {
			max = item
		}
	}
	if !found {
		return max, ErrSourceContainsNoElements
	}
	return max, nil
}

func (source Iterator[TSource]) Min(compare ...generic.Comparison[TSource]) (min TSource, err error) {
	found := false
	if len(compare) > 0 {
		for item := range source {
			if !found {
				min = item
				found = true
				continue
			}
			if compare[0](min, item) < 0 {
				min = item
			}
		}
		if !found {
			return min, ErrSourceContainsNoElements
		}
		return min, nil
	}
	if v, ok := (any(min)).(generic.IComparable[TSource]); ok {
		for item := range source {
			if !found {
				min = item
				found = true
				continue
			}
			if v.Compare(item) < 0 {
				min = item
			}
		}
		if !found {
			return min, ErrSourceContainsNoElements
		}
		return min, nil
	}
	switch (any(*new(TSource))).(type) {
	case int:
		min, _, err = minmax[TSource, int](source)
	case int8:
		min, _, err = minmax[TSource, int8](source)
	case int16:
		min, _, err = minmax[TSource, int16](source)
	case int32:
		min, _, err = minmax[TSource, int32](source)
	case int64:
		min, _, err = minmax[TSource, int64](source)
	case uint:
		min, _, err = minmax[TSource, uint](source)
	case uint8:
		min, _, err = minmax[TSource, uint8](source)
	case uint16:
		min, _, err = minmax[TSource, uint16](source)
	case uint32:
		min, _, err = minmax[TSource, uint32](source)
	case uint64:
		min, _, err = minmax[TSource, uint64](source)
	case uintptr:
		min, _, err = minmax[TSource, uintptr](source)
	case float32:
		min, _, err = minmax[TSource, float32](source)
	case float64:
		min, _, err = minmax[TSource, float64](source)
	case string:
		min, _, err = minmax[TSource, string](source)
	default:
		panic("unsupported type for Max")
	}
	return min, err
}

func Min[TSource generic.Comparable](source Iterator[TSource]) (min TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			found = true
			continue
		}
		if item < min {
			min = item
		}
	}
	if !found {
		return min, ErrSourceContainsNoElements
	}
	return min, nil
}

func (source Iterator[TSource]) MinMax(compare ...generic.Comparison[TSource]) (min TSource, max TSource, err error) {
	found := false
	if len(compare) > 0 {
		for item := range source {
			if !found {
				min = item
				max = item
				found = true
				continue
			}
			if compare[0](min, item) < 0 {
				min = item
				continue
			}
			if compare[0](max, item) > 0 {
				max = item
			}
		}
		if !found {
			return min, max, ErrSourceContainsNoElements
		}
		return min, max, nil
	}
	vmin, ok := (any(min)).(generic.IComparable[TSource])
	vmax, ok2 := (any(min)).(generic.IComparable[TSource])
	if ok && ok2 {
		for item := range source {
			if !found {
				min = item
				max = item
				found = true
				continue
			}
			if vmin.Compare(item) < 0 {
				min = item
			}
			if vmax.Compare(item) > 0 {
				max = item
			}
		}
		if !found {
			return min, max, ErrSourceContainsNoElements
		}
		return min, max, nil
	}
	switch (any(*new(TSource))).(type) {
	case int:
		min, max, err = minmax[TSource, int](source)
	case int8:
		min, max, err = minmax[TSource, int8](source)
	case int16:
		min, max, err = minmax[TSource, int16](source)
	case int32:
		min, max, err = minmax[TSource, int32](source)
	case int64:
		min, max, err = minmax[TSource, int64](source)
	case uint:
		min, max, err = minmax[TSource, uint](source)
	case uint8:
		min, max, err = minmax[TSource, uint8](source)
	case uint16:
		min, max, err = minmax[TSource, uint16](source)
	case uint32:
		min, max, err = minmax[TSource, uint32](source)
	case uint64:
		min, max, err = minmax[TSource, uint64](source)
	case uintptr:
		min, max, err = minmax[TSource, uintptr](source)
	case float32:
		min, max, err = minmax[TSource, float32](source)
	case float64:
		min, max, err = minmax[TSource, float64](source)
	case string:
		min, max, err = minmax[TSource, string](source)
	default:
		panic("unsupported type for Max")
	}
	return min, max, err
}

func MinMax[TSource generic.Comparable](source Iterator[TSource]) (min, max TSource, err error) {
	found := false
	for item := range source {
		if !found {
			min = item
			max = item
			found = true
			continue
		}
		if item < min {
			min = item
			continue
		}
		if item > max {
			max = item
		}
	}
	if !found {
		return min, max, ErrSourceContainsNoElements
	}
	return min, max, nil
}

func sort[TSource any, T generic.Comparable](source []TSource) {
	slices.SortFunc(source, func(x, y TSource) int {
		return cmp.Compare(any(x).(T), any(y).(T))
	})
}

func (source Iterator[TSource]) Order(compare ...generic.Comparison[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		result1 := make([]TSource, 0)
		for item := range source {
			result1 = append(result1, item)
		}
		if len(compare) > 0 {
			slices.SortFunc(result1, compare[0])
		} else if _, ok := (any(*new(TSource))).(generic.IComparable[TSource]); ok {
			slices.SortFunc(result1, func(first, second TSource) int {
				return (any(first)).(generic.IComparable[TSource]).Compare(second)
			})
		} else {
			switch (any(*new(TSource))).(type) {
			case int:
				sort[TSource, int](result1)
			case int8:
				sort[TSource, int8](result1)
			case int16:
				sort[TSource, int16](result1)
			case int32:
				sort[TSource, int32](result1)
			case int64:
				sort[TSource, int64](result1)
			case uint:
				sort[TSource, uint](result1)
			case uint8:
				sort[TSource, uint8](result1)
			case uint16:
				sort[TSource, uint16](result1)
			case uint32:
				sort[TSource, uint32](result1)
			case uint64:
				sort[TSource, uint64](result1)
			case uintptr:
				sort[TSource, uintptr](result1)
			case float32:
				sort[TSource, float32](result1)
			case float64:
				sort[TSource, float64](result1)
			case string:
				sort[TSource, string](result1)
			default:
				panic("unsupported type for Order")
			}
		}
		for _, item := range result1 {
			if !yield(item) {
				return
			}
		}
	}
}

func Order[TSource generic.Comparable](source Iterator[TSource], compare ...generic.Comparison[TSource]) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		result := make([]TSource, 0)
		for item := range source {
			result = append(result, item)
		}
		if len(compare) > 0 {
			slices.SortFunc(result, compare[0])
		} else {
			slices.Sort(result)
		}
		for _, item := range result {
			if !yield(item) {
				return
			}
		}
	}
}

func OrderBy[TSource any, TValue generic.Comparable](source Iterator[TSource], valueSelector generic.ValueSelector[TSource, TValue], compare ...generic.Comparison[TValue]) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		result := make([]generic.ValuePair[TSource, TValue], 0)
		for item := range source {
			result = append(result, generic.ValuePair[TSource, TValue]{
				Item1: item,
				Item2: valueSelector(item),
			})
		}
		if len(compare) > 0 {
			Compare := compare[0]
			slices.SortFunc(result, func(x, y generic.ValuePair[TSource, TValue]) int {
				return Compare(x.Item2, y.Item2)
			})
		} else {
			slices.SortFunc(result, func(x, y generic.ValuePair[TSource, TValue]) int {
				return cmp.Compare(x.Item2, y.Item2)
			})
		}
		for _, item := range result {
			if !yield(item.Item1) {
				return
			}
		}
	}
}

func sortDescending[TSource any, T generic.Real | generic.String](source []TSource) {
	slices.SortFunc(source, func(x, y TSource) int {
		return cmp.Compare(any(y).(T), any(x).(T))
	})
}

func (source Iterator[TSource]) OrderDescending(compare ...generic.Comparison[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		result1 := make([]TSource, 0)
		for item := range source {
			result1 = append(result1, item)
		}
		if len(compare) > 0 {
			Compare := compare[0]
			slices.SortFunc(result1, func(x, y TSource) int {
				return Compare(y, x)
			})
		} else if _, ok := (any(*new(TSource))).(generic.IComparable[TSource]); ok {
			slices.SortFunc(result1, func(first, second TSource) int {
				return (any(second)).(generic.IComparable[TSource]).Compare(first)
			})
		} else {
			switch (any(*new(TSource))).(type) {
			case int:
				sortDescending[TSource, int](result1)
			case int8:
				sortDescending[TSource, int8](result1)
			case int16:
				sortDescending[TSource, int16](result1)
			case int32:
				sortDescending[TSource, int32](result1)
			case int64:
				sortDescending[TSource, int64](result1)
			case uint:
				sortDescending[TSource, uint](result1)
			case uint8:
				sortDescending[TSource, uint8](result1)
			case uint16:
				sortDescending[TSource, uint16](result1)
			case uint32:
				sortDescending[TSource, uint32](result1)
			case uint64:
				sortDescending[TSource, uint64](result1)
			case uintptr:
				sortDescending[TSource, uintptr](result1)
			case float32:
				sortDescending[TSource, float32](result1)
			case float64:
				sortDescending[TSource, float64](result1)
			case string:
				sortDescending[TSource, string](result1)
			default:
				panic("unsupported type for Order")
			}
		}
		for _, item := range result1 {
			if !yield(item) {
				return
			}
		}
	}
}

func OrderDescending[TSource generic.Comparable](source Iterator[TSource], compare ...generic.Comparison[TSource]) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		result := make([]TSource, 0)
		for item := range source {
			result = append(result, item)
		}
		if len(compare) > 0 {
			Compare := compare[0]
			slices.SortFunc(result, func(x, y TSource) int {
				return Compare(y, x)
			})
		} else {
			slices.SortFunc(result, func(x, y TSource) int {
				return cmp.Compare(y, x)
			})
		}
		for _, item := range result {
			if !yield(item) {
				return
			}
		}
	}
}

func OrderByDescending[TSource any, TValue generic.Comparable](source Iterator[TSource], valueSelector generic.ValueSelector[TSource, TValue], compare ...generic.Comparison[TValue]) Iterator[TSource] {
	return func(yield func(value TSource) bool) {
		result := make([]generic.ValuePair[TSource, TValue], 0)
		for item := range source {
			result = append(result, generic.ValuePair[TSource, TValue]{
				Item1: item,
				Item2: valueSelector(item),
			})
		}
		if len(compare) > 0 {
			Compare := compare[0]
			slices.SortFunc(result, func(x, y generic.ValuePair[TSource, TValue]) int {
				return Compare(y.Item2, x.Item2)
			})
		} else {
			slices.SortFunc(result, func(x, y generic.ValuePair[TSource, TValue]) int {
				return cmp.Compare(y.Item2, x.Item2)
			})
		}
		for _, item := range result {
			if !yield(item.Item1) {
				return
			}
		}
	}
}

// Prepends values to the beggining of the sequence.
//
// # Parameters
//
//	elements TSource
//
// The values to prepend to source.
//
// # Returns
//
//	result Iterator[TSource]
//
// A new sequence that begins with elements.
func (source Iterator[TSource]) Prepend(elements ...TSource) (result Iterator[TSource]) {
	return func(yield func(TSource) bool) {
		for _, element := range elements {
			if !yield(element) {
				return
			}
		}
		for item := range source {
			if !yield(item) {
				return
			}
		}
	}
}

// Inverts the order of the elements in a sequence.
//
// # Returns
//
//	result Iterator[TSource]
//
// A sequence whose elements correspond to those of the input sequence in reverse order.
func (source Iterator[TSource]) Reverse() (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		reverse := make([]TSource, 0)
		for item := range source {
			reverse = append(reverse, item)
		}
		for i := len(reverse); i >= 0; i-- {
			if !yield(reverse[i]) {
				return
			}
		}
	}
}

func Select[TSource any, TResult any](source Iterator[TSource], valueSelector generic.ValueSelector[TSource, TResult]) Iterator[TResult] {
	return func(yield func(value TResult) bool) {
		for item := range source {
			if !yield(valueSelector(item)) {
				return
			}
		}
	}
}

func SelectMany[TSource any, TResult any](source Iterator[TSource], valueSelector generic.ValueSelector[TSource, []TResult]) Iterator[TResult] {
	return func(yield func(value TResult) bool) {
		for item := range source {
			innerSource := valueSelector(item)
			for _, innerItem := range innerSource {
				if !yield(innerItem) {
					return
				}
			}
		}
	}
}

// Determines whether two sequences are equal by comparing their elements by using a specified Equality[TSource].
//
// # Parameters
//
//	sequence Iterator[TSource]
//
// An Iterator[TSource] to compare to the source sequence.
//
//	comparer generic.Equality[TSource]
//
// An Equality[TSource] to use to compare elements. [OPTIONAL]
//
// # Returns
//
//	result bool
//
// True if the two sequences are equal of length and their corresponding elements compare equal according to comparer; otherwise, false.
func (source Iterator[TSource]) SequenceEqual(sequence Iterator[TSource], comparer ...generic.Equality[TSource]) (result bool) {
	next, stop := iter.Pull(iter.Seq[TSource](sequence))
	defer stop()
	if len(comparer) > 0 && comparer[0] != nil {
		Compare := comparer[0]
		for item1 := range source {
			item2, ok := next()
			if !ok || !Compare(item1, item2) {
				return false
			}
		}
	} else {
		for item1 := range source {
			item2, ok := next()
			if !ok {
				return false
			}
			if v, ok := (any(item1)).(generic.IEquatable[TSource]); ok {
				if !v.Equal(item2) {
					return false
				}
				continue
			}
			if !reflect.DeepEqual(item1, item2) {
				return false
			}
		}
	}
	_, ok := next()
	return !ok
}

// Returns the only element of a sequence or the only element of a sequence that satisfies a specified condition in predicate (if passed), and returns an error if none or more than one such element exists.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test an element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single element of the input sequence.
//
// # Error
//
//	err error
//
// ErrSourceContainsNoElements - When the input sequence is empty.
//
// ErrSourceHasMoreThanOneElement - When the input sequence has more then one element, predicate is not passed.
//
// ErrNoElementSatisfiesTheConditionInPredicate - When the input sequence contains no element that satifies a conditoin in passed predidate.
//
// ErrMoreThanOneElementSatisfiesTheConditionInPredicate - When the input sequence contains more than one element that satisfies a condition in passed predicate.
func (source Iterator[TSource]) Single(predicate ...generic.Predicate[TSource]) (result TSource, err error) {
	found := false
	any := false
	if len(predicate) > 0 && predicate[0] != nil {
		Predicate := predicate[0]
		for item := range source {
			any = true
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
		if !any {
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

// Returns the only element of a sequence or the only element of a sequence that satisfies a specified condition in predicate (if passed), or returns a default value if none or more than one such element exists.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test an element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single element of the input sequence, or a default value if no such element is found.
func (source Iterator[TSource]) SingleOrDefault(predicate ...generic.Predicate[TSource]) (result TSource) {
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

// Returns the only element of a sequence or the only element of a sequence that satisfies a specified condition in predicate (if passed), or returns a fallback value if none or more than one such element exists.
//
// # Parameters
//
//	fallback TSource
//
// The fallback value to return if the sequence is empty or constains more than one element that satisfies predicate (if passed).
//
//	predicate generic.Predicate[TSource]
//
// A function to test an element for a condition. [OPTIONAL]
//
// # Returns
//
//	result TSource
//
// The single element of the input sequence, or a fallback value if no such element is found.
func (source Iterator[TSource]) SingleOrFallback(fallback TSource, predicate ...generic.Predicate[TSource]) (result TSource) {
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
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains the elements that occur after the specified index in the input sequence.
//
// # Remarks
//
// If count is greater then collection length, this method returns an empty iterable collection.
func (source Iterator[TSource]) Skip(count int) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if count > 0 {
				count--
				continue
			}
			if !yield(item) {
				return
			}
		}
	}
}

// Returns a new iterable collection that contains the elements from source with the last <count> elements of the source collection omitted.
//
// # Parameters
//
//	count int
//
// The number of elements to omit from the end of the collection.
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains the elements from source minus count elements from the end of the collection.
//
// # Remarks
//
// If count is greater then collection length, this method returns an empty iterable collection.
func (source Iterator[TSource]) SkipLast(count int) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		count = source.Count() - count
		for item := range source {
			if count <= 0 {
				return
			}
			if !yield(item) {
				return
			}
			count--
		}
	}
}

// Bypasses elements in a sequence as long as a specified condition is true and then returns the remaining elements.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains the elements from the input sequence starting at the first element in the linear series that does not pass the test specified by predicate.
func (source Iterator[TSource]) SkipWhile(predicate generic.Predicate[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		skip := true
		for item := range source {
			if skip && !predicate(item) {
				continue
			}
			skip = false
			if !yield(item) {
				return
			}
		}
	}
}

func Sum[TValue generic.Number | generic.String](source Iterator[TValue]) (result TValue) {
	for item := range source {
		result += item
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
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains the specified number of elements from the start of the input sequence.
//
// # Remarks
//
// If count is not a positive number, this method returns an empty iterable collection.
func (source Iterator[TSource]) Take(count int) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if count <= 0 {
				return
			}
			if !yield(item) {
				return
			}
			count--
		}
	}
}

// Returns a new iterable collection that contains the last count elements from source.
//
// # Parameters
//
//	count int
//
// The number of elements to take from the end of the collection.
//
// # Returns
//
//	result Iterator[TSource]
//
// A new iterable collection that contains the last count elements from source.
//
// # Remarks
//
// If count is not a positive number, this method returns an empty iterable collection.
func (source Iterator[TSource]) TakeLast(count int) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		count = source.Count() - count
		for item := range source {
			if count > 0 {
				count--
				continue
			}
			if !yield(item) {
				return
			}
		}
	}
}

// Returns elements from a sequence as long as a specified condition is true, and then skips the remaining elements.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[Tsource] that contains the elements from the input sequence that occur before the element at which the test no longer passes.
func (source Iterator[TSource]) TakeWhile(predicate generic.Predicate[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if predicate(item) {
				return
			}
			if !yield(item) {
				return
			}
		}
	}
}

// Creates a slice of TSource from an Iterator[TSource]
//
// # Returns
//
//	result []TSource
//
// A slice of TSource that contains elements from the input sequence.
func (source Iterator[TSource]) ToSlice() (result []TSource) {
	result = make([]TSource, 0)
	for item := range source {
		result = append(result, item)
	}
	return result
}

// Creates a map[TKey]TValue from an Iterator[TSource] according to specified key selector and value selector functions.
//
// # Parameters
//
//	source Iterator[TSource]
//
// An Iterator[TSource] to create a map[TKey]TValue from.
//
//	keySelector generic.KeySelector[TSource, TKey]
//
// A function to extract a key from each element.
//
//	valueSelector generic.ValueSelector[TSource,TValue]
//
// A transform function to produce a result element value from each element.
//
// # Returns
//
//	result map[TKey],TValue
//
// A map[TKey]TValue that contains values of type TValue selected from the input sequence.
func ToMap[TSource any, TKey comparable, TValue any](source Iterator[TSource], keySelector generic.KeySelector[TSource, TKey], valueSelector generic.ValueSelector[TSource, TValue]) (result map[TKey]TValue) {
	result = make(map[TKey]TValue)
	for item := range source {
		result[keySelector(item)] = valueSelector(item)
	}
	return result
}

// Produces the set union of two sequences by using a specified Equality[TSource] function.
//
// # Parameters
//
//	sequence Iterator[TSource]
//
// An Iterator[TSource] whose distinct elements form the second set for the union.
//
//	comparer generic.Equality[TSource]
//
// The Equality[TSource] function to compare values.
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains the elements from both input sequences, excluding duplicates.
func (source Iterator[TSource]) Union(sequence Iterator[TSource], comparer ...generic.Equality[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source.Distinct(comparer...) {
			if !yield(item) {
				return
			}
		}
		for item := range sequence.Distinct(comparer...) {
			if !source.Contains(item, comparer...) {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Filters a sequence of values based on a predicate.
//
// # Parameters
//
//	predicate generic.Predicate[TSource]
//
// A function to test each element for a condition.
//
// # Returns
//
//	result Iterator[TSource]
//
// An Iterator[TSource] that contains elements from the input sequence that satisfy the condition in predicate.
func (source Iterator[TSource]) Where(predicate generic.Predicate[TSource]) (result Iterator[TSource]) {
	return func(yield func(value TSource) bool) {
		for item := range source {
			if predicate(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Produces a sequence of tuples with elements from the two specified sequences.
//
// # Parameters
//
//	first Iterator[TFirst]
//
// The first sequence to merge.
//
//	second Iterator[TSecond]
//
// The second sequence to merge.
//
// # Returns
//
//	resutl Iterator[generic.ValuePair[TFirst, TSecond]]
//
// A sequence of pairs with elements taken from the first and second sequences, in that order.
func Zip[TFirst any, TSecond any](source Iterator[TFirst], sequence Iterator[TSecond]) (result Iterator[generic.ValuePair[TFirst, TSecond]]) {
	return func(yield func(value generic.ValuePair[TFirst, TSecond]) bool) {
		next, stop := iter.Pull(iter.Seq[TSecond](sequence))
		defer stop()
		for item1 := range source {
			item2, ok := next()
			if !ok {
				return
			}
			if !yield(generic.ValuePair[TFirst, TSecond]{
				Item1: item1,
				Item2: item2,
			}) {
				return
			}
		}
	}
}
