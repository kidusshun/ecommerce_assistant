package llmclient

import (
	"errors"
	"fmt"
	"log"
)


type llmclient struct {
	tools Tools
}

func NewLlmClient(tools Tools) *llmclient {
	return &llmclient{
		tools: tools,
	}
}

func (client *llmclient) CallGemini(messageHistory []Message, tools []Tool) (*GeminiResponseBody, error) {

	geminiRequest := GeminiRequestBody{
		SystemInstruction: map[string]interface{}{
			"parts": map[string]string{
				"text": SystemInstruction,
			},
		},
		Contents: messageHistory,
		Tools:    tools,
		ToolConfig: FunctionCallingConfig{
			FunctionCallingConfig: Mode{
				Mode: AUTO,
			},
		},
	}
	res, err := GeminiClient(geminiRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}


func (client *llmclient) GenerateEmail(messageHistory []Message) (*GeminiResponseBody, error) {
	log.Println("history for email",messageHistory)
	geminiRequest := GeminiRequestBody{
		SystemInstruction: map[string]interface{}{
			"parts": map[string]string{
				"text": EmailSystemInstructions,
			},
		},
		Contents: messageHistory,
		ToolConfig: FunctionCallingConfig{
			FunctionCallingConfig: Mode{
				Mode: NONE,
			},
		},
		GenerationConfig: GenerationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"subject": map[string]string{
						"type": "string",
					},
					"body": map[string]string{
						"type": "string",
					},
				},
			},
			
		},
	}
	res, err := GeminiClient(geminiRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *llmclient) GenerateCouponEmail(couponCode string, discountPercentage float64, expirationDate string) (*GeminiResponseBody, error) {
	geminiRequest := GeminiRequestBody{
		SystemInstruction: map[string]interface{}{
			"parts": map[string]string{
				"text": CouponSystemInstructions,
			},
		},
		Contents: []Message{
			{
				Role: USER,
				Parts: []Part{
					{Text: fmt.Sprintf("the coupon has code %s with discount percentage of %s. It will expire on %s", couponCode, discountPercentage, expirationDate)},
				},
			},
		},
		ToolConfig: FunctionCallingConfig{
			FunctionCallingConfig: Mode{
				Mode: NONE,
			},
		},
		GenerationConfig: GenerationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"subject": map[string]string{
						"type": "string",
					},
					"body": map[string]string{
						"type": "string",
					},
				},
			},
			
		},
	}

	res, err := GeminiClient(geminiRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *llmclient) HandleFunctionCall(geminiresponse *GeminiResponseBody) (*ToolCallResponse, error) {

	log.Println(geminiresponse.Candidates[0].Content.Parts[0].Text)
	log.Println(geminiresponse.Candidates[0].Content.Parts[0].FunctionCall.Name)
	log.Println(geminiresponse.Candidates[0].Content.Parts[0].FunctionCall.Args)

	nameFunctionMap := map[string]interface{}{
		"QueryProducts": client.tools.QueryProducts,
		"CompanyInfo":   client.tools.CompanyInfo,
		"TrackOrder": 	 client.tools.TrackOrder,
	}

	functionCalls := geminiresponse.Candidates[0].Content.Parts

	for _, call := range functionCalls {
		if call.FunctionCall.Name != "" {
			functionName := call.FunctionCall.Name
			args := call.FunctionCall.Args
			if function, exists := nameFunctionMap[functionName]; exists {
				switch functionName {
				case "QueryProducts":
					query, ok := args["query"].(string)
					if !ok {
						return &ToolCallResponse{}, errors.New("Invalid argument for QueryProducts")
					}
					result, err := function.(func(string) (*ToolCallResponse, error))(query)
					if err != nil {
						return &ToolCallResponse{}, err
					}
					return result, nil
				case "CompanyInfo":
					query, ok := args["query"].(string)
					if !ok {
						return &ToolCallResponse{}, errors.New("Invalid argument for CompanyInfo")
					}
					result,err := function.(func(string) (*ToolCallResponse, error))(query)
					if err != nil {
						return &ToolCallResponse{},err
					}
					return result, nil
				case "TrackOrder":
					orderID, ok := args["orderID"].(string)
					if !ok {
						return &ToolCallResponse{}, errors.New("Invalid argument for TrackOrder")
					}
					result, err := function.(func(string) (*ToolCallResponse, error))(orderID)
					if err != nil {
						return &ToolCallResponse{}, err
					}
					return result, nil
				default:
					return &ToolCallResponse{},errors.New("Function not found")
				}
				} else {
				return &ToolCallResponse{},errors.New("Function not found")
			}
		}
	}
	return &ToolCallResponse{},errors.New("Function not found")
}
