package linq

// Represents a function that accumulates an object into an accumulator.
//
// # Parameters
//
//	accumulator TAccumulator
//
// The current accumulator value.
//
//	object T
//
// The object to accumulate.
//
// # Returns
//
//	TAccumulator
//
// The accumulated value produced from the accumulator and the object.
type Accumulate[TAccumulator any, T any] = func(accumulator TAccumulator, object T) TAccumulator

// Represents the method that has a single parameter and does not return a value.
//
// # Parameters
//
//	object T
//
// The parameter of the method that this delegate encapsulates.
type Action[T any] = func(object T)

// Represents the method that compares two objects of the same type.
//
// # Parameters
//
//	x T
//
// The first object to compare.
//
//	y T
//
// The second object to compare.
//
// # Returns
//
//	int
//
// A signed integer that indicates the relative values of x and y, as shown in the following table.
//
//	+-----------------------+----------------------+
//	| Result                | Meaning              |
//	+-----------------------+----------------------+
//	| Less than 0           | x is less than y.    |
//	| 0                     | x equals y.          |
//	| Greater than 0        | x is greater than y. |
//	+-----------------------+----------------------+
type Compare[T any] = func(x, y T) int

// Represents the method that reports two objects of the same type are equal.
//
// # Parameters
//
//	x T
//
// The first object to compare for equality.
//
//	y T
//
// The second object to compare for equality.
//
// # Returns
//
//	bool
//
// True if the specified objects are equal; otherwise, false.
type Equal[T any] = func(x, y T) bool

// Represents a function that creates a result element from two input elements.
//
// # Parameters
//
//	outer TOuter
//
// The first input element.
//
//	inner TInner
//
// The second input element.
//
// # Returns
//
//	TResult
//
// The result produced from the input elements.
type Join[TOuter any, TInner any, TResult any] = func(outer TOuter, inner TInner) TResult

// Represents a function that extracts a key from the specified object.
//
// # Parameters
//
//	object T
//
// The object to extract a key from.
//
// # Returns
//
//	TKey
//
// The key extracted from the object.
type Key[T any, TKey comparable] = func(object T) TKey

// Represents a function that defines a set of criteria and determines whether the specified object meets those criteria.
//
// # Parameters
//
//	object T
//
// The object to compare against the criteria defined within the method represented by this Predicate.
//
// # Returns
//
//	bool
//
// True if object meets the criteria defined within the method represented by this Predicate; otherwise, false.
type Predicate[T any] = func(object T) bool

// Represents a function that produces a value from the specified object.
//
// # Parameters
//
//	object T
//
// The object to produce a value from.
//
// # Returns
//
//	TValue
//
// The value produced from the object.
type Value[T any, TValue any] = func(object T) TValue
