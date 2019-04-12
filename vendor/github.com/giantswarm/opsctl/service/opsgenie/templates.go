package opsgenie

import (
	"text/template"
)

var fnplus = template.FuncMap{
	"plus1": func(x int) int {
		return x + 1
	},
}

// escalationTmpl describes object, used for creating new escalations
// in OpsGenie (https://docs.opsgenie.com/docs/escalation-api).
const escalationTmpl = `{
    "name" : "{{.Name}}",
    "rules" : [
        {
            "delay": {
                "timeAmount" : 1
            },
            "recipient":{
                "type" : "user",
                "username": "{{.User}}"
            },
            "notifyType" : "default",
            "condition": "if-not-acked"
        }
    ],
    "ownerTeam" : {
        "name" : "ops_team"
    },
    "repeat": {
      "waitInterval": 5,
      "count": 20,
      "resetRecipientStates": false,
      "closeAlertAfterAll": false
    }
}`

// routingRuleTmpl describes object, used for creating new routing rules
// in OpsGenie (https://docs.opsgenie.com/docs/team-routing-rule-api).
const routingRuleTmpl = `{
    "name": "{{.Name}}",
    "order": 0,
    "criteria": {
        "type": "{{.Type}}",
        "conditions": [
            {{ if ne .Cluster "" -}}
            {
                "field": "message",
                "operation": "contains",
                "expectedValue": "{{ .Cluster }}"
            },
            {{ end -}}
            {{$n := len .Conditions }}{{ range $i, $e := .Conditions -}}
            {
                "field": "description",
                "not": {{$e.Not}},
                "operation": "contains",
                "expectedValue": "{{$e.Value}}"
            }{{if ne (plus1 $i) $n}},{{end}}
            {{ end -}}
        ]
    },
    "notify": {
        "name":"{{.Name}}",
        "type":"escalation"
    }
}`

// alertTmpl describes object, used for creating new alerts
// in OpsGenie (https://docs.opsgenie.com/docs/alert-api).
const alertTmpl = `{
    "message": "{{.Message}}",
    "responders": [
        {
            "type": "team",
            "name": "{{.Team}}"
        },
	{
	    "type": "escalation",
	    "name": "{{.Team}}_panic_escalation"
	}
    ]
}`
