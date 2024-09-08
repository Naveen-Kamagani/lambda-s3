package main

// import (
// 	"aws-lambda-s3/repositories"
// 	"context"
// 	"net/http"
// 	"os"
// 	"strings"

// 	ddlambda "github.com/DataDog/datadog-lambda-go"
// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// )

// var (
// 	env        = os.Getenv("ENV")
// 	region     = os.Getenv("AWS_REGION")
// 	bucketName = os.Getenv("BUCKET_NAME")
// 	handler    Handler
// )

// type Handler struct {
// 	S3Client *repositories.S3Client
// }

// type Msg struct {
// 	Message string
// }

// type Response struct {
// 	Objects               []string `json:"objects"`
// 	NextContinuationToken string   `json:"nextContinuationToken"`
// }

// func initializeAWSCfg(ctx context.Context) (cfg aws.Config, error error) {
// 	if strings.ToLower(env) == "dev" {
// 		return config.LoadDefaultConfig(
// 			ctx,
// 			config.WithRegion(region),
// 			config.WithClientLogMode(aws.LogRequest),
// 			config.WithClientLogMode(aws.LogRequestWithBody),
// 			config.WithClientLogMode(aws.LogResponseWithBody),
// 			config.WithClientLogMode(aws.LogResponse),
// 		)
// 	}
// 	return config.LoadDefaultConfig(ctx, config.WithRegion(region))
// }

// func init() {
// 	awsCfg, awsCfgErr := initializeAWSCfg(context.TODO())
// 	if awsCfgErr != nil {
// 		panic("unable to initialize AWS Config")
// 	}

// 	s3Client := repositories.S3Client{
// 		Client: repositories.InitializeS3Client(awsCfg),
// 	}
// 	handler = Handler{
// 		S3Client: &s3Client,
// 	}
// }

// func (h *Handler) handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	var response events.APIGatewayProxyResponse
// 	bucketName = os.Getenv("BUCKET_NAME")
// 	if bucketName == "" {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       "BUCKET_NAME environment variable is not set",
// 		}, nil
// 	}

// 	if h.S3Client == nil {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       "S3Client is not initialized",
// 		}, nil
// 	}

// 	switch request.HTTPMethod {
// 	case "GET":
// 		// objectKey := request.PathParameters["objectKey"]
// 		// if objectKey == "" {
// 		// 	// Fetch all objects using pagination
// 		// 	var allObjects []string
// 		// 	continuationToken := ""

// 		// 	// Call ListObjectsV2 with the current continuation token
// 		// 	result, err := h.S3Client.ListObjectsV2(context.TODO(), bucketName)
// 		// 	if err != nil {
// 		// 		return events.APIGatewayProxyResponse{
// 		// 			StatusCode: http.StatusInternalServerError,
// 		// 			Body:       err.Error(),
// 		// 		}, nil
// 		// 	}

// 		// 	// Append the current set of objects to the list
// 		// 	for _, object := range result.Contents {
// 		// 		allObjects = append(allObjects, aws.ToString(object.Key))
// 		// 	}

// 		// 	// // Check if there's a continuation token for more results
// 		// 	if result.IsTruncated != nil && *result.IsTruncated {
// 		// 		continuationToken = string(*result.NextContinuationToken)
// 		// 	} else {
// 		// 		continuationToken = ""
// 		// 	}
// 		// 	// Return the complete list of objects
// 		// 	resp := Response{
// 		// 		Objects:               allObjects,
// 		// 		NextContinuationToken: continuationToken, // No continuation token since all objects are fetched
// 		// 	}
// 		// 	body, _ := json.Marshal(resp)
// 		// 	// response = events.APIGatewayProxyResponse{
// 		// 	// 	StatusCode: http.StatusOK,
// 		// 	// 	Body:       string(body),
// 		// 	// }
// 		// 	fmt.Printf("Objects in ListObjectsV2 - %s", string(body))
// 		// 	//body, _ := json.Marshal(resp)
// 		// 	response = events.APIGatewayProxyResponse{
// 		// 		StatusCode: http.StatusOK,
// 		// 		Body:       string(body),
// 		// 	}
// 		// } else {
// 		// 	// Get a single object
// 		// 	getObjectOutput, err := h.S3Client.GetObject(context.TODO(), bucketName, objectKey)
// 		// 	if err != nil {
// 		// 		return events.APIGatewayProxyResponse{
// 		// 			StatusCode: http.StatusInternalServerError,
// 		// 			Body:       err.Error(),
// 		// 		}, nil
// 		// 	}
// 		// 	defer getObjectOutput.Body.Close()
// 		// 	body, _ := io.ReadAll(getObjectOutput.Body)
// 		// 	response = events.APIGatewayProxyResponse{
// 		// 		StatusCode: http.StatusOK,
// 		// 		Body:       string(body),
// 		// 	}
// 		// }
// 		return handler.GetHandler(request, bucketName)
// 	case "DELETE":
// 		// objectKey := request.PathParameters["objectKey"]
// 		// if objectKey == "" {
// 		// 	return events.APIGatewayProxyResponse{
// 		// 		StatusCode: http.StatusInternalServerError,
// 		// 		Body:       fmt.Sprintf("Object key is passed as nil -  %s", objectKey),
// 		// 	}, nil
// 		// }
// 		// _, err := h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
// 		// if err != nil {
// 		// 	return events.APIGatewayProxyResponse{
// 		// 		StatusCode: http.StatusInternalServerError,
// 		// 		Body:       fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()),
// 		// 	}, nil
// 		// }
// 		// response = events.APIGatewayProxyResponse{
// 		// 	StatusCode: http.StatusOK,
// 		// 	Body:       fmt.Sprintf("Deleted %s successfully", objectKey),
// 		// }
// 		return handler.DeleteHandler(request, bucketName)
// 	}

// 	return response, nil
// }

// func main() {
// 	lambda.Start(ddlambda.WrapFunction(handler.handleRequest, nil))
// }
