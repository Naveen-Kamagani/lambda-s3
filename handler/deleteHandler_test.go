package handler

import (
	"aws-lambda-s3/repositories"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name               string
		request            events.APIGatewayProxyRequest
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "valid",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "object1.txt"},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Deleted object1.txt successfully",
		},
		{
			name: "invalid",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": ""},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "object key is passed as nil - ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				S3Client: &repositories.S3Client{
					Client: &repositories.S3Mock{
						MockTest: tt.name,
					},
				},
			}
			response, err := h.DeleteHandler(tt.request, "test-bucket")
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tt.expectedBody, response.Body)
		})
	}
}
