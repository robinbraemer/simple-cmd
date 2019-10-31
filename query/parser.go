package query

import (
	"errors"
	"strings"
)

type Query interface {
	Elements() []Element
	Element(key string) Element
	Run(Context)
}

type Element interface {
	Key() string
	Type() ElementType
	Optional() bool
}

type ElementType string

const (
	// A required argument that equals the element key.
	ElementTypeArgument ElementType = "ARGUMENT"
	// A dynamic argument, may be a string or an array.
	ElementTypeValue ElementType = "VALUE"
)

// .......................

type query struct {
	elements []Element
	run      func(ctx Context)
}

func (q *query) Elements() []Element {
	return q.elements
}
func (q *query) Element(key string) Element {
	for _, e := range q.elements {
		if e.Key() == key {
			return e
		}
	}
	return nil
}

func (q *query) Run(ctx Context) {
	q.run(ctx)
}

type Context interface {
	Get(key string) (value string, exists bool)
	Require(key string) string

	Array(key string) (value []string, exists bool)
	RequireArray(key string) []string
}

func New(rawQuery string, fn func(Context)) (Query, error) {
	elements, err := parse(rawQuery)
	if err != nil {
		return nil, err
	}
	return &query{
		elements: elements,
		run:      fn,
	}, nil
}

func parse(rawQuery string) ([]Element, error) {
	rawElements := strings.Split(rawQuery, " ")
	elements := make([]Element, 0, len(rawElements))
	keys := make(map[string]struct{})
	for _, e := range rawElements {
		element, err := parseElement(e)
		if err != nil {
			return nil, err
		}
		if _, ok := keys[element.Key()]; ok {
			return nil, errors.New("key names must be unique")
		}
		keys[element.Key()] = struct{}{}
		elements = append(elements, element)
	}
	return elements, nil
}

func parseElement(rawElement string) (Element, error) {
	var optional bool
	if rawElement[0] == '{' {
		if rawElement[len(rawElement)-1] != '}' {
			return nil, errors.New("missing closing bracket")
		}
		last := len(rawElement) - 1
		if rawElement[last-1:last] == "?" {
			optional = true
			last--
		}
		key := rawElement[1:last]
		if len(key) == 0 {
			return nil, errors.New("missing key name")
		}
		return &element{
			key:         key,
			elementType: ElementTypeValue,
			optional:    optional,
		}, nil
	}
	if rawElement[len(rawElement)-1] == '}' {
		return nil, errors.New("missing opening brackets")
	}
	return &element{
		key:         rawElement,
		elementType: ElementTypeArgument,
		optional:    false,
	}, nil
}

type element struct {
	// The key name of the element.
	key string
	// The element type.
	elementType ElementType
	// Whether the element is an array.
	array bool
	// The size of the array.
	// Must be >= 2 or nil as infinitely as the last element in the query
	arraySize *uint
	// Whether the element is optional.
	// Must be the last element in the query.
	optional bool
}

func (e *element) Key() string {
	return e.key
}

func (e *element) Type() ElementType {
	return e.elementType
}

func (e *element) Optional() bool {
	return e.optional
}
