package bundle

import (
	"errors"
	"idunno/query"
	"strings"
)

type Bundle interface {
	Add(...query.Query)
	List() []query.Query
	Run(string) error
}

func New() Bundle {
	return &bundle{}
}

type bundle struct {
	queries []query.Query
}

func (b *bundle) Add(queries ...query.Query) {
	b.queries = append(b.queries, queries...)
}

func (b *bundle) List() []query.Query {
	return b.queries
}

func (b *bundle) Run(s string) error {
	//fmt.Println()
	//fmt.Println("running:", s)
	args := strings.Split(strings.TrimSpace(s), " ")
	indexExists := func(i int) bool { return len(args)-1 >= i }

	ctx := &context{}

	for _, q := range b.queries {
		ctx.args = make(map[string]string)
		//fmt.Println("query", j)

		for i, e := range q.Elements() {
			//fmt.Printf("%+v\n", e)

			if e.Type() == query.ElementTypeArgument {
				// element must exist and equal the key
				if !indexExists(i) || e.Key() != args[i] {
					break
				}
			} else if e.Type() == query.ElementTypeValue {
				if indexExists(i) {
					if len(q.Elements())-1 == i {
						// is last argument, append rest of args to key
						ctx.args[e.Key()] = strings.Join(args[i:], " ")
					} else {
						ctx.args[e.Key()] = args[i]
					}
				} else if !e.Optional() {
					// arg not exists & required
					break
				}

				// arg is optional, means must be last arg
				q.Run(ctx)
				return nil
			}
		}
	}
	return errors.New("no matching query found")
}

type context struct {
	args map[string]string
}

func (ctx *context) Array(key string) (value []string, exists bool) {
	panic("implement me")
}

func (ctx *context) RequireArray(key string) []string {
	panic("implement me")
}

func (ctx *context) Get(key string) (string, bool) {
	v, exists := ctx.args[key]
	return v, exists
}
func (ctx *context) Require(key string) string {
	v := ctx.args[key]
	return v
}
