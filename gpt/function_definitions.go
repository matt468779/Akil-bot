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
				"keyword": {
					Type:        jsonschema.String,
					Description: "The keyword could be the name of the organization, the opportunity type or category etc. this value is going to be searched in every field and try to pick one",
				},
				"opportunityType": {
					Type: jsonschema.String,
					Enum: []string{"virtual", "in-person"},
				},
				"categories": {
					Type: jsonschema.String,
					Enum: []string{"Technology", "Education", "Health", "Economic Development", "Food and Nutrition", "Disaster Response and Recovery", "Child and Maternal Health"},
				},
				"organizationType": {
					Type: jsonschema.String,
					Enum: []string{"local", "international"},
				},
				"location": {
					Type:        jsonschema.String,
					Description: "Location of the opportunity",
				},
			},
			Required: []string{"keyword"},
		},
	},
}

type OrganizationFilter struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type OpportunityFilter struct {
	Keyword          string   `json:"keyword" url:"keyword"`
	OpportunityType  string   `json:"opportunityType" url:"opportunityType"`
	Categories       []string `json:"categories" url:"categories"`
	OrganizationType string   `json:"organizationType" url:"organizationType"`
	Location         []string `json:"location" url:"location"`
}
