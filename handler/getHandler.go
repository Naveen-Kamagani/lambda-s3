package handler

import (
	"aws-lambda-s3/repositories"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Response struct {
	Objects               []string `json:"objects"`
	NextContinuationToken string   `json:"nextContinuationToken"`
}

type Handler struct {
	S3Client *repositories.S3Client
}

func (h *Handler) GetHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
	//objectKey := request.PathParameters["objectKey"]
	objectKey := request.QueryStringParameters["objectKey"]
	if objectKey == "" {
		// Fetch all objects using pagination
		var allObjects []string
		continuationToken := ""

		// Call ListObjectsV2 with the current continuation token
		result, err := h.S3Client.ListObjectsV2(context.TODO(), bucketName)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}

		// Append the current set of objects to the list
		for _, object := range result.Contents {
			allObjects = append(allObjects, aws.ToString(object.Key))
		}

		// // Check if there's a continuation token for more results
		if result.IsTruncated != nil && *result.IsTruncated {
			continuationToken = string(*result.NextContinuationToken)
		} else {
			continuationToken = ""
		}
		// Return the complete list of objects
		resp := Response{
			Objects:               allObjects,
			NextContinuationToken: continuationToken, // No continuation token since all objects are fetched
		}
		body, _ := json.Marshal(resp)
		// response = events.APIGatewayProxyResponse{
		// 	StatusCode: http.StatusOK,
		// 	Body:       string(body),
		// }
		fmt.Printf("Objects in ListObjectsV2 - %s", string(body))
		//body, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	} else {
		// Get a single object
		getObjectOutput, err := h.S3Client.GetObject(context.TODO(), bucketName, objectKey)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}
		defer getObjectOutput.Body.Close()
		body, _ := io.ReadAll(getObjectOutput.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	}
}
