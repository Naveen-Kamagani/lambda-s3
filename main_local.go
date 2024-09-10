package main

import (
	"aws-lambda-s3/handler"
	"aws-lambda-s3/repositories"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	env    = os.Getenv("ENV")
	region = os.Getenv("AWS_REGION")
	hndlr  MyHandler
)

// type Handler struct {
// 	S3Client *repositories.S3Client
// }

type MyHandler struct {
	Handler handler.Handler
}

type Msg struct {
	Message string
}

type Response struct {
	Objects               []string `json:"objects"`
	NextContinuationToken string   `json:"nextContinuationToken"`
}

func initializeAWSCfg(ctx context.Context) (cfg aws.Config, error error) {
	if strings.ToLower(env) == "dev" {
		return config.LoadDefaultConfig(
			ctx,
			config.WithRegion(region),
			config.WithClientLogMode(aws.LogRequest),
			config.WithClientLogMode(aws.LogRequestWithBody),
			config.WithClientLogMode(aws.LogResponseWithBody),
			config.WithClientLogMode(aws.LogResponse),
		)
	}
	return config.LoadDefaultConfig(ctx, config.WithRegion(region))
}

func main() {
	// Initialize the AWS configuration
	awsCfg, awsCfgErr := initializeAWSCfg(context.TODO())
	if awsCfgErr != nil {
		panic("unable to initialize AWS Config")
	}

	// Initialize S3 Client
	s3Client := repositories.S3Client{
		Client: repositories.InitializeS3Client(awsCfg),
	}
	hndlr = MyHandler{
		Handler: handler.Handler{
			S3Client: &s3Client,
		},
	}

	// Create a mock API Gateway request
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		//HTTPMethod: "DELETE",
		// PathParameters: map[string]string{
		// 	"source":                   "create_folder",
		// 	"action":                   "SELECT",
		// 	"baseEncodedDocumentTitle": "UFJUWV9JRCNDVVNUX0lECg==",
		// },
		PathParameters: map[string]string{
			"source":                   "create_folder",
			"action":                   "INSERT",
			"baseEncodedDocumentTitle": "",
		},
		// PathParameters: map[string]string{
		// 	"source":                   "create_folder",
		// 	"action":                   "INSERT",
		// 	"baseEncodedDocumentTitle": "UFJUWV9JRCNDVVNUX0lECg==",
		// },
	}

	// Set environment variables locally
	os.Setenv("BUCKET_NAME", "bcdr-lambda-function")
	os.Setenv("ENV", "dev")
	os.Setenv("AWS_REGION", "us-east-1")

	// Call the handler function
	response, err := hndlr.handleRequest(request)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Response: ", response)
	}
}

func (h *MyHandler) handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	// if bucketName == "" {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Body:       "BUCKET_NAME environment variable is not set",
	// 	}, nil
	// }

	// if h.Handler.S3Client == nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Body:       "S3Client is not initialized",
	// 	}, nil
	// }

	switch request.HTTPMethod {
	case "GET":
		return h.Handler.GetHandler(request, bucketName)
	case "DELETE":
		return h.Handler.DeleteHandler(request, bucketName)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Unsupported HTTP method",
		}, nil
	}
}
