package linq

// Represents a pair of key and value.
//
//	TKey
//
// The type of the key.
//
//	TValue
//
// The type of the value.
type KeyValuePair[TKey any, TValue any] struct {
	// Key is the key component of the pair.
	Key TKey
	// Value is the value component of the pair.
	Value TValue
}
