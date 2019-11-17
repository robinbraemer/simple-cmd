package query

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Query interface {
	Elements() []Element
	Element(key string) Element
	Run(Context) error
}

type Element interface {
	Key() string
	Type() ElementType
	Optional() bool
	IsArray() bool
	ArraySize() *int
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
	fn       reflect.Value // the function
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

func (q *query) Run(ctx Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error running query: %v", r)
		}
	}()

	fnType := q.fn.Type()
	numIn := fnType.NumIn()
	args := make([]reflect.Value, 0, numIn)

	for i := 0; i < numIn; i++ {
		switch fnType.In(i) {
		case reflect.TypeOf((*Context)(nil)).Elem():
			args = append(args, reflect.ValueOf(ctx))
		default:
			return errors.New("unsupported query function")
		}
	}

	q.fn.Call(args)
	return nil
}

type Context interface {
	Get(key string) (value string, exists bool)
	Require(key string) string

	Array(key string) (value []string, exists bool)
	RequireArray(key string) []string
}

func New(rawQuery string, fn interface{}) (Query, error) {
	if fn == nil || reflect.TypeOf(fn).Kind() != reflect.Func {
		return nil, errors.New("fn must be a function")
	}

	elements, err := parse(rawQuery)
	if err != nil {
		return nil, err
	}
	return &query{
		elements: elements,
		fn:       reflect.ValueOf(fn),
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
	/*
		"{a}"
		"{a?}"
		"{[]a}"
		"{[]a?}"
		"{[3]a}"
		"{[3]a?}"
	*/
	if rawElement[0] == '{' {
		// element is a value
		var optional bool
		var array bool
		var arraySize *int

		if rawElement[len(rawElement)-1] != '}' {
			return nil, errors.New("missing closing bracket")
		}

		leftCursor := 1
		if rawElement[1] == '[' {
			// value is an array
			array = true
			// find closing bracket
			arrayClosingAt := 2
			for {
				if rawElement[arrayClosingAt] == ']' {
					break
				}
				arrayClosingAt++
				if arrayClosingAt == len(rawElement)-1 {
					// end of element
					return nil, errors.New("missing closing array bracket")
				}
			}
			leftCursor = arrayClosingAt + 1
			if arrayClosingAt != 2 {
				// array is not infinite as it doesn't close immediately -> "[]"
				i, err := strconv.Atoi(rawElement[2:arrayClosingAt])
				if err != nil {
					return nil, errors.New("invalid array size")
				}
				if i < 1 {
					// array size to small
					return nil, errors.New("array size must not be < 1")
				}
				arraySize = &i
			}
		}

		last := len(rawElement) - 1
		if rawElement[last-1] == '?' {
			optional = true
			last--
		}
		key := rawElement[leftCursor:last]
		if len(key) == 0 {
			return nil, errors.New("missing key name")
		}
		return &element{
			key:         key,
			elementType: ElementTypeValue,
			array:       array,
			arraySize:   arraySize,
			optional:    optional,
		}, nil
	} else if rawElement[len(rawElement)-1] == '}' {
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
	arraySize *int
	// Whether the element is optional.
	// Must be the last element in the query.
	optional bool
}

func (e *element) IsArray() bool {
	return e.array
}

func (e *element) ArraySize() *int {
	return e.arraySize
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
