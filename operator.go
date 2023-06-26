package go_mongo

var (
	Eq    = "$eq"
	Ne    = "$ne"
	Gt    = "$gt"
	Lt    = "$lt"
	Gte   = "$gte"
	Lte   = "$lte"
	Regex = "$regex"
)
var ComparisonOperators = map[string][]string{
	"uint":   []string{Eq, Ne, Gt, Lt, Gte, Lte},
	"uint64": []string{Eq, Ne, Gt, Lt, Gte, Lte},
	"int":    []string{Eq, Ne, Gt, Lt, Gte, Lte},
	"int64":  []string{Eq, Ne, Gt, Lt, Gte, Lte},
	"string": []string{Eq, Ne, Regex},
	"bool":   []string{Eq, Ne},
}
