package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) DeleteHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
	//objectKey := request.PathParameters["objectKey"]
	objectKey := request.QueryStringParameters["objectKey"]
	if objectKey == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("object key is passed as nil - %s", objectKey),
		}, nil
	}
	_, err := h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("Deleted %s successfully", objectKey),
	}, nil
}
