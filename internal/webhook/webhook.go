package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	api.UnimplementedWebhookServiceServer
}

func NewService(options ...func(*Service)) *Service {
	service := Service{}
	for _, option := range options {
		option(&service)
	}
	return &service
}

func (s *Service) Webhook(ctx context.Context, input *api.WebhookRequest) (*api.WebhookResponse, error) {
	if input.Headers == nil {
		*input.Headers = structpb.Struct{}
	}
	if input.Body == nil {
		*input.Body = structpb.Struct{}
	}
	// Convert header and body to io.Reader
	bodyBytes := new(bytes.Buffer)
	json.NewEncoder(bodyBytes).Encode(input.Body)
	request, err := http.NewRequest(input.Method, input.Uri, bodyBytes)
	if err != nil {
		return nil, err
	}
	// Set headers
	for key, value := range input.Headers.Fields {
		request.Header.Set(key, value.GetStringValue())
	}
	// Send the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// Convert response body to structpb.Value
	responseBody := new(map[string]interface{})
	json.NewDecoder(response.Body).Decode(responseBody)
	responseStruct, err := conversion.MapToProtobufStruct(*responseBody)
	if err != nil {
		return nil, err
	}
	// Convert response headers to structpb.Value
	responseHeaders := new(map[string]interface{})
	for key, value := range response.Header {
		(*responseHeaders)[key] = value
	}
	responseHeadersStruct, err := conversion.MapToProtobufStruct(*responseHeaders)
	if err != nil {
		return nil, err
	}
	return &api.WebhookResponse{
		Status:  uint32(response.StatusCode),
		Headers: responseHeadersStruct,
		Body:    responseStruct,
	}, nil
}
