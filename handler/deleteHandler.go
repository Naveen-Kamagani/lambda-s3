package handler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

// func (h *Handler) DeleteHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
// 	//objectKey := request.PathParameters["objectKey"]
// 	objectKey := request.PathParameters["objectKey"]
// 	if objectKey == "" {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       fmt.Sprintf("object key is passed as nil - %s", objectKey),
// 		}, nil
// 	}
// 	_, err := h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()),
// 		}, nil
// 	}
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       fmt.Sprintf("Deleted %s successfully", objectKey),
// 	}, nil
// }

func (h *Handler) DeleteHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
	pathObjectKey := request.PathParameters["objectKey"]
	objectKey, err := url.PathUnescape(pathObjectKey)
	if err != nil {
		return generateErrorResponse(bucketName, objectKey, "Invalid objectKey")
	}
	if objectKey == "" {
		return generateErrorResponse(bucketName, objectKey, "object key is passed as nil")
	}

	_, err = h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
	if err != nil {
		return generateErrorResponse(bucketName, objectKey, fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()))
	}

	return generateSuccessResponse(bucketName, objectKey, "", fmt.Sprintf("Deleted %s successfully", objectKey))
}
