package handler

import (
	"aws-lambda-s3/repositories"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Response struct {
	Objects               []string `json:"objects"`
	NextContinuationToken string   `json:"nextContinuationToken"`
}

type OutputResponse struct {
	BucketName   string    `json:"bucketName"`
	Key          string    `json:"key"`
	ErrorMessage string    `json:"errorMessage"`
	Timestamp    time.Time `json:"timestamp"`
	Body         string    `json:"body"`
}

type Handler struct {
	S3Client *repositories.S3Client
}

func (h *Handler) GetHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
	source, err := url.PathUnescape(request.PathParameters["source"])
	if err != nil {
		return generateErrorResponse(bucketName, "", "Invalid source")
	}
	action, err := url.PathUnescape(request.PathParameters["action"])
	if err != nil {
		return generateErrorResponse(bucketName, "", "Invalid action")
	}
	baseEncodedDocumentTitle, err := url.PathUnescape(request.PathParameters["baseEncodedDocumentTitle"])
	if err != nil {
		return generateErrorResponse(bucketName, "", "Invalid baseEncodedDocumentTitle")
	}

	var bucketPrefix, objectKey string
	if source == "" && action == "" && baseEncodedDocumentTitle == "" {
		bucketPrefix = ""
	}
	if source != "" && action == "" && baseEncodedDocumentTitle == "" {
		bucketPrefix = fmt.Sprintf("%s/", source)
	}
	if source != "" && action != "" && baseEncodedDocumentTitle == "" {
		bucketPrefix = fmt.Sprintf("%s/%s/", source, action)
	}
	if source != "" && action != "" && baseEncodedDocumentTitle != "" {
		decodedKey, err := base64.URLEncoding.DecodeString(baseEncodedDocumentTitle)
		if err != nil {
			return generateErrorResponse(bucketName, objectKey, "Invalid objectKey")
		}
		objectKey = strings.TrimSpace(fmt.Sprintf("%s/%s/%s", source, action, string(decodedKey)))
	}

	//fmt.Println("Object Key = ", objectKey)
	if objectKey == "" {
		// Fetch all objects using pagination
		var allObjects []string
		continuationToken := ""

		result, err := h.S3Client.ListObjectsV2(context.TODO(), bucketName, bucketPrefix)
		if err != nil {
			return generateErrorResponse(bucketName, objectKey, err.Error())
		}

		for _, object := range result.Contents {
			allObjects = append(allObjects, aws.ToString(object.Key))
		}

		if result.IsTruncated != nil && *result.IsTruncated {
			continuationToken = string(*result.NextContinuationToken)
		}

		resp := Response{
			Objects:               allObjects,
			NextContinuationToken: continuationToken,
		}
		body, _ := json.Marshal(resp)
		return generateSuccessResponse(bucketName, objectKey, "", string(body))
	} else {

		getObjectOutput, err := h.S3Client.GetObject(context.TODO(), bucketName, objectKey)
		if err != nil {
			return generateErrorResponse(bucketName, objectKey, err.Error())
		}
		defer getObjectOutput.Body.Close()

		body, _ := io.ReadAll(getObjectOutput.Body)
		return generateSuccessResponse(bucketName, objectKey, "", string(body))
	}
}

// Helper function to generate a success response
func generateSuccessResponse(bucketName, key, errorMessage string, successOutput string) (events.APIGatewayProxyResponse, error) {
	output := OutputResponse{
		BucketName:   bucketName,
		Key:          key,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now(),
		Body:         successOutput,
	}
	body, _ := json.Marshal(output)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

// Helper function to generate an error response
func generateErrorResponse(bucketName, key, errorMessage string) (events.APIGatewayProxyResponse, error) {
	output := OutputResponse{
		BucketName:   bucketName,
		Key:          key,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now(),
	}
	body, _ := json.Marshal(output)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(body),
	}, nil
}
