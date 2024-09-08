package main

// import (
// 	"aws-lambda-s3/handler"
// 	"aws-lambda-s3/repositories"
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// )

// var (
// 	env    = os.Getenv("ENV")
// 	region = os.Getenv("AWS_REGION")
// 	//bucketName = os.Getenv("BUCKET_NAME")
// 	hndlr MyHandler
// )

// type MyHandler struct {
// 	Handler handler.Handler
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
// 	hndlr = MyHandler{
// 		Handler: handler.Handler{
// 			S3Client: &s3Client,
// 		},
// 	}
// }

// type Response struct {
// 	Objects               []string `json:"objects"`
// 	NextContinuationToken string   `json:"nextContinuationToken"`
// }

// func main() {

// 	//return response, nil
// 	//handler.handleRequest()
// 	hndlr.handleRequest()

// }

// func (h *MyHandler) handleRequest() {
// 	bucketName := os.Getenv("BUCKET_NAME")
// 	objectKey := os.Getenv("OBJECT_KEY")
// 	methodCall := os.Getenv("METHOD_CALL")
// 	//var request events.APIGatewayProxyRequest
// 	// switch methodCall {
// 	// case "GET":
// 	// 	//objectKey := request.PathParameters["objectKey"]
// 	// 	if objectKey == "" {
// 	// 		// Fetch all objects using pagination
// 	// 		var allObjects []string
// 	// 		continuationToken := ""
// 	// 		// Call ListObjectsV2 with the current continuation token
// 	// 		result, err := h.S3Client.ListObjectsV2(context.TODO(), bucketName)
// 	// 		if err != nil {
// 	// 			// return events.APIGatewayProxyResponse{
// 	// 			// 	StatusCode: http.StatusInternalServerError,
// 	// 			// 	Body:       err.Error(),
// 	// 			// }, nil
// 	// 			fmt.Printf("error in listobjectsv2 - %v", err.Error())
// 	// 		}

// 	// 		// Append the current set of objects to the list
// 	// 		for _, object := range result.Contents {
// 	// 			allObjects = append(allObjects, aws.ToString(object.Key))
// 	// 		}

// 	// 		fmt.Printf("\ncontinuation token - %v", *result.IsTruncated)

// 	// 		// // Check if there's a continuation token for more results
// 	// 		if *result.IsTruncated {
// 	// 			continuationToken = string(*result.NextContinuationToken)
// 	// 		} else {
// 	// 			continuationToken = ""
// 	// 		}
// 	// 		// Return the complete list of objects
// 	// 		resp := Response{
// 	// 			Objects:               allObjects,
// 	// 			NextContinuationToken: continuationToken, // No continuation token since all objects are fetched
// 	// 		}
// 	// 		body, _ := json.Marshal(resp)
// 	// 		// response = events.APIGatewayProxyResponse{
// 	// 		// 	StatusCode: http.StatusOK,
// 	// 		// 	Body:       string(body),
// 	// 		// }
// 	// 		fmt.Printf("Objects in ListObjectsV2 - %s", string(body))
// 	// 	} else {
// 	// 		// Get a single object
// 	// 		getObjectOutput, err := h.S3Client.GetObject(context.TODO(), bucketName, objectKey)
// 	// 		if err != nil {
// 	// 			// return events.APIGatewayProxyResponse{
// 	// 			// 	StatusCode: http.StatusInternalServerError,
// 	// 			// 	Body:       err.Error(),
// 	// 			// }, nil
// 	// 			fmt.Printf("error in GetObject - %v", err)
// 	// 		}
// 	// 		defer getObjectOutput.Body.Close()
// 	// 		body, _ := io.ReadAll(getObjectOutput.Body)
// 	// 		// response = events.APIGatewayProxyResponse{
// 	// 		// 	StatusCode: http.StatusOK,
// 	// 		// 	Body:       string(body),
// 	// 		// }
// 	// 		fmt.Printf("Successfully tested GetObject - %s and %s", objectKey, string(body))
// 	// 	}

// 	// case "DELETE":
// 	// 	//objectKey := request.PathParameters["objectKey"]
// 	// 	_, err := h.S3Client.DeleteObject(context.TODO(), bucketName, objectKey)
// 	// 	if err != nil {
// 	// 		// return events.APIGatewayProxyResponse{
// 	// 		// 	StatusCode: http.StatusInternalServerError,
// 	// 		// 	Body:       fmt.Sprintf("Error deleting %s: %v", objectKey, err.Error()),
// 	// 		// }, nil
// 	// 		fmt.Printf("Error in DeleteObject - %v", err)
// 	// 	}
// 	// 	// response = events.APIGatewayProxyResponse{
// 	// 	// 	StatusCode: http.StatusOK,
// 	// 	// 	Body:       fmt.Sprintf("Deleted %s successfully", objectKey),
// 	// 	// }
// 	// 	fmt.Printf("Successfully tested in DeleteObject - %s", objectKey)
// 	// }

// 	switch methodCall {
// 	case "GET":
// 		fmt.Println(h.Handler.GetHandler(objectKey, bucketName))
// 	case "DELETE":
// 		fmt.Println(h.Handler.DeleteHandler(objectKey, bucketName))
// 	default:
// 		fmt.Printf("Unsupported HTTP method")
// 		// return events.APIGatewayProxyResponse{
// 		// 	StatusCode: http.StatusMethodNotAllowed,
// 		// 	Body:       "Unsupported HTTP method",
// 		// }, nil
// 	}
// }
