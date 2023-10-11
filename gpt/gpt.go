package gpt

import (
	"akil_telegram_bot/bootstrap"
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var env *bootstrap.Env
var availableFunctions = map[string]func(string) string{"GetOrganization": GetOrganizations, "GetOpportunity": GetOpportunity}

func SetEnv(e *bootstrap.Env) {
	env = e
}

func GetResponse(messages []openai.ChatCompletionMessage) string {
	client := openai.NewClient(env.OpenAIAPIKey)
	systemMessage, err := os.ReadFile(env.SystemMessagePath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	systemMessageChat := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: string(systemMessage)}
	tempMessages := []openai.ChatCompletionMessage{systemMessageChat}

	tempMessages = append(tempMessages, messages...)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  tempMessages,
			Functions: Functions,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "Try again later"
	}

	responseMessage := resp.Choices[0].Message
	if responseMessage.FunctionCall != nil {
		fmt.Println("____________Function Called______________")
		functionName := resp.Choices[0].Message.FunctionCall.Name
		functionToCall := availableFunctions[functionName]
		functionArguments := responseMessage.FunctionCall.Arguments

		functionResponse := functionToCall(functionArguments)

		tempMessages = append(tempMessages, responseMessage)
		tempMessages = append(tempMessages, openai.ChatCompletionMessage{
			Name:    functionName,
			Role:    openai.ChatMessageRoleFunction,
			Content: fmt.Sprint(functionResponse),
		})

		secondResponse, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: tempMessages,
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return "Try again later"
		}

		print(secondResponse.Choices[0].Message.Content)
		return secondResponse.Choices[0].Message.Content
	}
	return resp.Choices[0].Message.Content
}
