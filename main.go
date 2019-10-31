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

	q, err := query.New("say {text}", func(ctx query.Context) {
		fmt.Println("say", ctx.Require("text"))
	})
	if err != nil {
		log.Fatal(err)
	}
	q2, err := query.New("say {text?}", func(ctx query.Context) {
		fmt.Println("Usage: say <text>")
	})
	if err != nil {
		log.Fatal(err)
	}

	b := bundle.New()
	b.Add(q, q2)

	err = b.Run("say Hello world!")
	if err != nil {
		log.Fatal(err)
	}
	err = b.Run("say")
	if err != nil {
		log.Fatal(err)
	}
}
