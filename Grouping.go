package linq

// Grouping represents a group of source elements associated with a key.
//
// A Grouping can be iterated to access the elements that belong to the group.
type Grouping[TKey comparable, TSource any] struct {
	// Key is the value shared by the elements in the group.
	Key TKey
	enumerable[TSource]
}