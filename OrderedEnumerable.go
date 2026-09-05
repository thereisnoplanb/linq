package linq

import "cmp"

// OrderedEnumerable represents a sequence whose elements have an ordering criterion.
//
// It supports applying additional ordering criteria through ThenBy and ThenDescendingBy methods.
type OrderedEnumerable[TSource any] struct {
	*orderedEnumerable[TSource]
}

type orderedEnumerable[TSource any] struct {
	enumerable[TSource]
	comparators comparators[TSource]
}

type comparators[TSource any] []Compare[TSource]

// Compares two elements using the sequence of comparators.
//
// # Parameters
//
//	x TSource
//
// The first element to compare.
//
//	y TSource
//
// The second element to compare.
//
// # Returns
//
//	result int
//
// A negative value if x is less than y, zero if x equals y, and a positive value if x is greater than y.
func (this comparators[TSource]) Compare(x, y TSource) int {
	for _, compare := range this {
		if result := compare(x, y); result != 0 {
			return result
		}
	}
	return 0
}

// Performs a subsequent ascending ordering of the elements in an ordered sequence according to a specified key selector function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenBy[TValue cmp.Ordered](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return cmp.Compare(value(x), value(y))
	})
	return OrderedEnumerable[TSource]{source}
}

// Performs a subsequent ascending ordering of the elements in an ordered sequence according to a specified key selector function by using the Compare method.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenByComparable[TValue iComparable[TValue]](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return value(x).Compare(value(y))
	})
	return OrderedEnumerable[TSource]{source}
}

// Performs a subsequent ascending ordering of the elements in an ordered sequence according to a specified key selector and comparison function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
//	compare Compare[TValue]
//
// A function to compare the selected values.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenByFunc[TValue any](value Value[TSource, TValue], compare Compare[TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return compare(value(x), value(y))
	})
	return OrderedEnumerable[TSource]{source}
}

// Performs a subsequent descending ordering of the elements in an ordered sequence according to a specified key selector function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values in descending order.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenDescendingBy[TValue cmp.Ordered](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return cmp.Compare(value(y), value(x))
	})
	return OrderedEnumerable[TSource]{source}
}

// Performs a subsequent descending ordering of the elements in an ordered sequence according to a specified key selector function by using the Compare method.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values in descending order.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenDescendingByComparable[TValue iComparable[TValue]](value Value[TSource, TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return value(y).Compare(value(x))
	})
	return OrderedEnumerable[TSource]{source}
}

// Performs a subsequent descending ordering of the elements in an ordered sequence according to a specified key selector and comparison function.
//
// # Parameters
//
//	value Value[TSource, TValue]
//
// A function that selects the value used to order each element.
//
//	compare Compare[TValue]
//
// A function to compare the selected values.
//
// # Returns
//
//	result OrderedEnumerable[TSource]
//
// An OrderedEnumerable[TSource] ordered by the selected values in descending order.
//
// # Remarks
//
// This ordering is applied after the existing orderings. The result can be further ordered using ThenBy or ThenDescendingBy methods.
func (source *orderedEnumerable[TSource]) ThenDescendingByFunc[TValue any](value Value[TSource, TValue], compare Compare[TValue]) (result OrderedEnumerable[TSource]) {
	source.comparators = append(source.comparators, func(x, y TSource) int {
		return compare(value(y), value(x))
	})
	return OrderedEnumerable[TSource]{source}
}
