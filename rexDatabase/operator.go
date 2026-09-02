package rexDatabase

type ConditionOperator string

const (
	ConditionEq ConditionOperator = "eq"
	ConditionNe ConditionOperator = "ne"

	ConditionGt  ConditionOperator = "gt"
	ConditionGte ConditionOperator = "gte"
	ConditionLt  ConditionOperator = "lt"
	ConditionLte ConditionOperator = "lte"

	ConditionIn    ConditionOperator = "in"
	ConditionNotIn ConditionOperator = "notIn"

	ConditionContains   ConditionOperator = "contains"
	ConditionStartsWith ConditionOperator = "startsWith"
	ConditionEndsWith   ConditionOperator = "endsWith"

	ConditionLike  ConditionOperator = "like"
	ConditionILike ConditionOperator = "ilike"

	ConditionBetween ConditionOperator = "between"

	ConditionIsNull    ConditionOperator = "isNull"
	ConditionIsNotNull ConditionOperator = "isNotNull"

	ConditionRegex ConditionOperator = "regex"
)

type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "and"
	LogicalOr  LogicalOperator = "or"
	LogicalNot LogicalOperator = "not"
)
