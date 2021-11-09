package q2sql

const ErrUndefinedCondition = Error("condition is not defined")

// Condition returns wrapped SQL representation which describes certain condition
// and a list of the condition arguments
type Condition func(field string, args ...interface{}) Sqlizer

// ConditionFactory creates specific conditions by the name
type ConditionFactory interface {
	// CreateCondition creates condition by the name
	CreateCondition(name string) (Condition, error)
}

// ConditionMap maps a string key to a Condition.
// It is used to provide a simple way of provisioning ConditionFactory.
type ConditionMap map[string]Condition

// CreateCondition implements ConditionFactory
func (m ConditionMap) CreateCondition(name string) (Condition, error) {
	c, ok := m[name]
	if !ok {
		return nil, ErrUndefinedCondition
	}
	return c, nil
}
