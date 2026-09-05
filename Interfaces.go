package linq

// Defines a generalized comparison method that is implemented to create a type-specific comparison method for ordering or sorting its instances.
//
// Type Parameters
//
//	T
//
// The type of object to compare.
type iComparable[T any] interface {

	// Compares the current instance with another object of the same type and returns an integer that indicates whether the current instance precedes,
	// follows, or occurs in the same position in the sort order as the other object.
	//
	// # Parameters
	//
	//	other T
	//
	// An object to compare with current object.
	//
	// # Returns
	//
	//	int
	//
	// A signed integer that indicates the relative values of current and another, as shown in the following table.
	//
	//	+-----------------------+-------------------------------+
	//	| Result                | Meaning                       |
	//	+-----------------------+-------------------------------+
	//	| Less than 0           | current is less than other.   |
	//	| 0                     | current equals other.         |
	//	| Greater than 0        | current is greater than other.|
	//	+-----------------------+-------------------------------+
	Compare(other T) int
}

// Defines a generalized method that is implemented to create a type-specific method for determining equality of instances.
//
// Type Parameters
//
//	T
//
// The type of objects to compare.
type iEquatable[T any] interface {

	// Indicates whether the current object is equal to another object of the same type.
	//
	// # Parameters
	//
	//	other T
	//
	// An object to compare with current object.
	//
	// # Returns
	//
	//	bool
	//
	// True if the current object is equal to the other parameter; otherwise, false.
	Equal(other T) bool
}
