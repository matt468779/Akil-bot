package gpt

import (
	"fmt"
	"io"
	"net/http"
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
	res, err := http.Get(env.BackendURL + "opportunities/all")
	if err != nil {
		return "Something went wrong, Please try again later"
	}
	result := fmt.Sprint(io.ReadAll(res.Body))
	println(result)
	return result
}
