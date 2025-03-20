package chat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/llmclient"
)

type Service struct {
	store  ChatStore
	client llmclient.LlmClient
}

func NewChatService(store ChatStore, client llmclient.LlmClient) *Service {
	return &Service{
		store:  store,
		client: client,
	}
}

// should return a response to the router
func (chat *Service) Chat(request ChatRequest) (ChatResponse, error) {
	chatHistory, err := chat.store.GetChatHistory(uuid.New())
	if err != nil {
		return ChatResponse{}, nil
	}
	tools := GetTools()

	if len(request.Attachment) > 0 {
		parts := []llmclient.Part{
			{
				Text: request.Message,
			},
		}
		for _, attachment := range request.Attachment {
			parts = append(parts, llmclient.Part{
				InlineData: &llmclient.ImageData{
					MimeType: attachment.MimeType,
					Data:     string(attachment.Content),
				},
			})
		}
		chatHistory = append(chatHistory, llmclient.Message{
			Role:  llmclient.USER,
			Parts: parts,
		})
	} else {
		chatHistory = append(chatHistory, llmclient.Message{
			Role: llmclient.USER,
			Parts: []llmclient.Part{
				{
					Text: request.Message,
				},
			},
		})
	}
	response, err := chat.client.CallGemini(chatHistory, tools)
	if err != nil {
		return ChatResponse{}, err
	}
	chatHistory = append(chatHistory, llmclient.Message{
		Role:  llmclient.MODEL,
		Parts: response.Candidates[0].Content.Parts,
	})

	if err != nil {
		return ChatResponse{}, err
	}

	var chatResponse ChatResponse

	for response.Candidates[0].Content.Parts[0].FunctionCall != nil {
		call_result, err := chat.client.HandleFunctionCall(response)
		if err != nil {
			log.Println(err)
			return ChatResponse{}, nil
		}

		chatHistory = append(chatHistory, call_result.ModelResponse)

		if call_result.Products != nil {
			chatResponse.Products = append(chatResponse.Products, call_result.Products...)
		}
		if call_result.Location != "" {
			chatResponse.Location = call_result.Location
		}

		response, err = chat.client.CallGemini(chatHistory, tools)

		chatHistory = append(chatHistory, llmclient.Message{
			Role:  llmclient.MODEL,
			Parts: response.Candidates[0].Content.Parts,
		})
		if err != nil {
			return ChatResponse{}, err
		}
		str, err := json.Marshal(response)
		if err != nil {
			log.Print(err)
			return ChatResponse{}, err
		}
		fmt.Println(string(str))
	}
	Messsages = chatHistory
	MessageHistory = append(MessageHistory, ChatResponse{
		Content: request.Message,
		Attachment:   request.Attachment,
		Role:         "user",
	})

	chatResponse.Content = response.Candidates[0].Content.Parts[0].Text
	chatResponse.Role = "model"
	MessageHistory = append(MessageHistory, chatResponse)

	return chatResponse, nil

}

func GetTools() []llmclient.Tool {
	return []llmclient.Tool{
		{
			FunctionDeclarations: []llmclient.FunctionDeclaration{
				{
					Name:        "QueryProducts",
					Description: "a function to get products that matches the query passed",
					Parameters: llmclient.Parameters{
						Type: "object",
						Properties: map[string]llmclient.Property{
							"query": {
								Type: "string",
							},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "CompanyInfo",
					Description: "a function to ask questions about a company's identity and general info",
					Parameters: llmclient.Parameters{
						Type: "object",
						Properties: map[string]llmclient.Property{
							"query": {
								Type: "string",
							},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "TrackOrder",
					Description: "a function to get the location of a user's order ",
					Parameters: llmclient.Parameters{
						Type: "object",
						Properties: map[string]llmclient.Property{
							"orderID": {
								Type: "string",
							},
						},
						Required: []string{"orderID"},
					},
				},
			},
		},
	}
}
