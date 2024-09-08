package handler

import (
	"aws-lambda-s3/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// Helper function to generate the expected body for success responses
func generateExpectedBody(bucketName, key, errorMessage string, objectCount int, nextContinuationToken string) (string, error) {
	var objectList []string
	var resp Response
	output := OutputResponse{
		BucketName:   bucketName,
		Key:          key,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now(),
	}
	if errorMessage != "" {
		output.Body = ""
		body, err := json.Marshal(output)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	if objectCount == 1 {
		objectContent := "{\"SomethingCool\":\"cool beans\"}"
		output.Body = objectContent
		body, err := json.Marshal(output)
		if err != nil {
			return "", err
		}
		return string(body), nil
	} else {
		for i := 1; i <= objectCount; i++ {
			objectList = append(objectList, fmt.Sprintf("object%d.txt", i))
		}
		resp = Response{
			Objects:               objectList,
			NextContinuationToken: nextContinuationToken,
		}
	}
	respBody, _ := json.Marshal(resp)

	output.Body = string(respBody)
	body, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name               string
		request            events.APIGatewayProxyRequest
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "less-than-50",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": ""},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				body, _ := generateExpectedBody("test-bucket", "", "", 30, "")
				return body
			}(),
		},
		{
			name: "valid",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "object1.txt"},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				body, _ := generateExpectedBody("test-bucket", "object1.txt", "", 1, "")
				return body
			}(),
		},
		{
			name: "invalid key",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "non-existent.txt"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: func() string {
				body, _ := generateExpectedBody("test-bucket", "non-existent.txt", "invalid key or object key does not exist in the bucket", 0, "")
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
			response, err := h.GetHandler(tt.request, "test-bucket")
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
