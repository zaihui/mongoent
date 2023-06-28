package mongoent

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
	"uint":   {Eq, Ne, Gt, Lt, Gte, Lte},
	"uint64": {Eq, Ne, Gt, Lt, Gte, Lte},
	"int":    {Eq, Ne, Gt, Lt, Gte, Lte},
	"int64":  {Eq, Ne, Gt, Lt, Gte, Lte},
	"string": {Eq, Ne, Regex},
	"bool":   {Eq, Ne},
}
