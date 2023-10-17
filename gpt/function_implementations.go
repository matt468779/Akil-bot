package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Filter interface {
}

func GetOrganizations(arguments string) string {

	res, err := http.Get(env.BackendURL + "organizations/all")
	if err != nil {
		return "Something went wrong, Please try again later"
	}
	return fmt.Sprint(io.ReadAll(res.Body))
}

func GetOpportunity(arguments string) string {
	url := env.BackendURL + "opportunities/search?"

	opportunityFilter := OpportunityFilter{}
	json.Unmarshal([]byte(arguments), &opportunityFilter)

	queryString, err := query.Values(opportunityFilter)
	if err != nil {
		return "Something went wrong, Please try again later"
	}
	println(arguments)
	println(url + queryString.Encode())
	res, err := http.Get(url + queryString.Encode())
	if err != nil {
		return "Something went wrong, Please try again later"
	}
	result := fmt.Sprint(io.ReadAll(res.Body))
	println(result)
	return result
}
