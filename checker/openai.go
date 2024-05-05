package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message represents a single message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestBody defines the structure of the API request body.
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ResponseChoice represents a choice object in the API response.
type ResponseChoice struct {
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// APIResponse defines the structure of the API response.
type APIResponse struct {
	Choices []ResponseChoice `json:"choices"`
	ID      string           `json:"id"`
	Model   string           `json:"model"`
	Object  string           `json:"object"`
	Usage   map[string]int   `json:"usage"`
}

// Client holds the configuration for an API client.
type Client struct {
	APIKey string
	URL    string
    Model string
}

// NewClient initializes a new instance of Client.
func NewClient(apiKey, url, model string) *Client {
	return &Client{
		APIKey: apiKey,
		URL:    url,
        Model: model,
	}
}

func (c *Client) call(messages []Message) ([]ResponseChoice, error) {
    requestBody := RequestBody{
        Model:    c.Model,
        Messages: messages,
    }

    response, err := c.sendRequest(requestBody)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }

    return response.Choices, nil
}


// sendRequest sends a request to the OpenAI API and returns the parsed response.
func (c *Client) sendRequest(requestBody RequestBody) (*APIResponse, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to API: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response data: %v", err)
	}

	var apiResponse APIResponse
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	return &apiResponse, nil
}
