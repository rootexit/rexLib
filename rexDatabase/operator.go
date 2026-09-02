package rexDatabase

type Operator string

const (
	OpEq Operator = "eq"
	OpNe Operator = "ne"

	OpGt  Operator = "gt"
	OpGte Operator = "gte"
	OpLt  Operator = "lt"
	OpLte Operator = "lte"

	OpIn    Operator = "in"
	OpNotIn Operator = "notIn"

	OpContains   Operator = "contains"
	OpStartsWith Operator = "startsWith"
	OpEndsWith   Operator = "endsWith"

	OpLike  Operator = "like"
	OpILike Operator = "ilike"

	OpBetween Operator = "between"

	OpIsNull    Operator = "isNull"
	OpIsNotNull Operator = "isNotNull"

	OpRegex Operator = "regex"
)
