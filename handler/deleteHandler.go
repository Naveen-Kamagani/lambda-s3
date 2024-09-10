package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) DeleteHandler(request events.APIGatewayProxyRequest, bucketName string) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
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
	var objectKey string

	if source != "" && action != "" && baseEncodedDocumentTitle != "" {
		decodedKey, err := base64.URLEncoding.DecodeString(baseEncodedDocumentTitle)
		if err != nil {
			return generateErrorResponse(bucketName, objectKey, "Invalid objectKey")
		}
		objectKey = strings.TrimSpace(fmt.Sprintf("%s/%s/%s", source, action, string(decodedKey)))
	}

	if objectKey == "" {
		return generateErrorResponse(bucketName, objectKey, "object key is passed as nil")
	}

	_, err = h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
	if err != nil {
		return generateErrorResponse(bucketName, objectKey, fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()))
	}

	return generateSuccessResponse(bucketName, objectKey, "", fmt.Sprintf("Deleted %s successfully", objectKey), headers)
}
