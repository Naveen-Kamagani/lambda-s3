package handler

import (
	"aws-lambda-s3/repositories"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// func TestDeleteHandler(t *testing.T) {
// 	tests := []struct {
// 		name               string
// 		request            events.APIGatewayProxyRequest
// 		expectedStatusCode int
// 		expectedBody       string
// 	}{
// 		{
// 			name: "valid",
// 			request: events.APIGatewayProxyRequest{
// 				QueryStringParameters: map[string]string{"objectKey": "object1.txt"},
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedBody:       "Deleted object1.txt successfully",
// 		},
// 		{
// 			name: "invalid",
// 			request: events.APIGatewayProxyRequest{
// 				QueryStringParameters: map[string]string{"objectKey": ""},
// 			},
// 			expectedStatusCode: http.StatusInternalServerError,
// 			expectedBody:       "object key is passed as nil - ",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &Handler{
// 				S3Client: &repositories.S3Client{
// 					Client: &repositories.S3Mock{
// 						MockTest: tt.name,
// 					},
// 				},
// 			}
// 			response, err := h.DeleteHandler(tt.request, "test-bucket")
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
// 			assert.Equal(t, tt.expectedBody, response.Body)
// 		})
// 	}
// }

func generateDeleteExpectedBody(bucketName, key, errorMessage, successMessage string) (string, error) {
	output := OutputResponse{
		BucketName:   bucketName,
		Key:          key,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now(),
		Body:         successMessage,
	}
	body, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

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
			expectedBody: func() string {
				body, _ := generateDeleteExpectedBody("test-bucket", "object1.txt", "", "Deleted object1.txt successfully")
				return body
			}(),
		},
		{
			name: "invalid",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": ""},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: func() string {
				body, _ := generateDeleteExpectedBody("test-bucket", "", "object key is passed as nil", "")
				return body
			}(),
		},
		{
			name: "error",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "object2.txt"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: func() string {
				body, _ := generateDeleteExpectedBody("test-bucket", "object2.txt", "Error deleting object2.txt: error deleting object", "")
				return body
			}(),
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

			// Compare the body while ignoring the timestamp
			var expectedOutput, actualOutput OutputResponse
			_ = json.Unmarshal([]byte(tt.expectedBody), &expectedOutput)
			_ = json.Unmarshal([]byte(response.Body), &actualOutput)

			expectedOutput.Timestamp = actualOutput.Timestamp // Ignore timestamp for comparison
			assert.Equal(t, expectedOutput, actualOutput)
		})
	}
}
