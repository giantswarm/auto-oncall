package opsgenie

type Alert struct {
	Message string
	Team    string
}

type Condition struct {
	Order         int    `json:"order"`
	Not           bool   `json:"not"`
	ExpectedValue string `json:"expectedValue"`
}

type Criteria struct {
	Type       string      `json:"type"`
	Conditions []Condition `json:"conditions"`
}

// Internal representation.
type RoutingRule struct {
	Name       string
	User       string
	Cluster    string
	Conditions []Rule
	Type       string
}

type Rule struct {
	Not   bool
	Value string
}

type TeamRoutingRule struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IsDefault bool     `json:"isDefault"`
	Criteria  Criteria `json:"criteria"`
}

type TeamRoutingRuleData struct {
	Rule *TeamRoutingRule `json:"data"`
}

type TeamRoutingRules struct {
	Rules []TeamRoutingRule `json:"data"`
}
