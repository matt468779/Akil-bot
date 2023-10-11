package gpt

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var Functions = []openai.FunctionDefinition{
	{
		Name:        "GetOrganization",
		Description: "This function will get a list of organization",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"location": {
					Type:        jsonschema.String,
					Description: "The city and state, e.g. San Francisco, CA",
				},
				"unit": {
					Type: jsonschema.String,
					Enum: []string{"celcius", "fahrenheit"},
				},
			},
			Required: []string{"location"},
		},
	},
	{
		Name:        "GetOpportunity",
		Description: "This function will get a list of opportunities",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"name": {
					Type:        jsonschema.String,
					Description: "The name of the organization",
				},
				"address": {
					Type:        jsonschema.String,
					Description: "Location of the opportunity",
				},
			},
			Required: []string{"location"},
		},
	},
}

type OrganizationFilter struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type OpportunityFilter struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
