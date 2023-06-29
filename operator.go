package mongoent

var (
	Eq    = "$eq"
	Ne    = "$ne"
	Gt    = "$gt"
	Lt    = "$lt"
	Gte   = "$gte"
	Lte   = "$lte"
	Regex = "$regex"
	In    = "$in"
	NotIn = "$nin"
)

var ComparisonOperators = map[string][]string{
	"uint":   {Eq, Ne, Gt, Lt, Gte, Lte},
	"uint64": {Eq, Ne, Gt, Lt, Gte, Lte},
	"int":    {Eq, Ne, Gt, Lt, Gte, Lte},
	"int64":  {Eq, Ne, Gt, Lt, Gte, Lte},
	"string": {Eq, Ne, Regex},
	"bool":   {Eq, Ne},
}

var ComparisonInOperators = map[string][]string{
	"uint":   {In, NotIn},
	"uint64": {In, NotIn},
	"int":    {In, NotIn},
	"int64":  {In, NotIn},
	"string": {In, NotIn},
	"bool":   {In, NotIn},
}
