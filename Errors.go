package linq

import "errors"

// ErrSourceContainsNoElements indicates that the source sequence contains no elements.
var ErrSourceContainsNoElements = errors.New("the source contains no elements")

// ErrNoElementSatisfiesTheConditionInPredicate indicates that no source element satisfies the predicate.
var ErrNoElementSatisfiesTheConditionInPredicate = errors.New("no element satisfies the condition in predicate")

// ErrMoreThanOneElementSatisfiesTheConditionInPredicate indicates that more than one source element satisfies the predicate.
var ErrMoreThanOneElementSatisfiesTheConditionInPredicate = errors.New("more than one element satisfies the condition in predicate")

// ErrSourceHasMoreThanOneElement indicates that the source sequence contains more than one element.
var ErrSourceHasMoreThanOneElement = errors.New("the source has more than one element")

// ErrSizeIsBelowOne indicates that a requested size is less than one.
var ErrSizeIsBelowOne = errors.New("size is below 1")

// ErrIndexOutOfRange indicates that an index is outside the bounds of a sequence.
var ErrIndexOutOfRange = errors.New("index out of range")

// ErrTypeIsNotNumberOrString indicates that a type is neither numeric nor a string type.
var ErrTypeIsNotNumberOrString = errors.New("type is not number or string")

// ErrTypeIsNotNumber indicates that a type is not numeric.
var ErrTypeIsNotNumber = errors.New("type is not number")
