package main

import (
	"fmt"
	"idunno/bundle"
	"idunno/query"
	"log"
)

func main() {
	/*
		"say {a}"
		"say {a?}"
		"say {[]a}"
		"say {[]a?}"
		"say {[3]a}"
		"say {[3]a} test"
		"say {[3]a} test {[]b}"
		"say {[3]a} test {[]b?}"
		"say {[3]a} test {[]b?}"
		"say {[3]a} test {b?}"
	*/

	b1 := mustBundle(
		mustQuery{
			rawQuery: "say {text}",
			fn: func(ctx query.Context) {
				fmt.Println("say", ctx.Require("text"))
			},
		},
		mustQuery{
			rawQuery: "say {text?}",
			fn: func(ctx query.Context) {
				fmt.Println("Usage: say <text>")
			},
		},
		mustQuery{
			rawQuery: "say {text?}",
			fn: func(ctx query.Context) {
				fmt.Println("Usage: say <text>")
			},
		},
		mustQuery{
			rawQuery: "hi lol",
			fn: func(ctx query.Context) {
				fmt.Println("hi lol")
			},
		},
	)

	err := b1.Run("say Hello world!")
	check(err)
	err = b1.Run("say")
	check(err)
	err = b1.Run("hi lol")
	check(err)

	fmt.Println()
	fmt.Println()
	fmt.Println()
	/*
		"say {a}"
		"say {a?}"
		"say {[]a}"
		"say {[]a?}"
		"say {[3]a}"
		"say {[3]a} test"
		"say {[3]a} test {[]b}"
		"say {[3]a} test {[]b?}"
		"say {[3]a} test {[]b?}"
		"say {[3]a} test {b?}"
	*/
	b2 := mustBundle(
		mustQuery{
			rawQuery: "say {text?}",
			fn: func(s string) {
				fmt.Println("REFLECTED")
			},
		},
	)

	err = b2.Run("say")
	check(err)

}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustBundle(queries ...mustQuery) bundle.Bundle {
	b := bundle.New()
	for _, q := range queries {
		que, err := query.New(q.rawQuery, q.fn)
		must(err)
		b.Add(que)
	}
	return b
}

type mustQuery struct {
	rawQuery string
	fn       interface{}
}
