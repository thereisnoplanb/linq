package linq

import "errors"

var ErrSourceContainsNoElements = errors.New("the source contains no elements")
var ErrNoElementSatisfiesTheConditionInPredicate = errors.New("no element satisfies the condition in predicate")
var ErrMoreThanOneElementSatisfiesTheConditionInPredicate = errors.New("more than one element satisfies the condition in predicate")
var ErrSourceHasMoreThanOneElement = errors.New("the source has more than one element")
var ErrSizeIsBelowOne = errors.New("size is below 1")
var ErrIndexOutOfRange = errors.New("index out of range")
