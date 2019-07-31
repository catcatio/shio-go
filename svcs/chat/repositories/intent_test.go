package repositories

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"encoding/json"
	"fmt"
	"github.com/catcatio/shio-go/nub"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"testing"
)

func TestX(t *testing.T) {

	credential := ``

	ctx := context.Background()
	client, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsJSON([]byte(credential)))
	fmt.Println(err)
	fmt.Println(client)

	projectID := ""
	sessionID := nub.NewID()
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	message := "buy  meiji"

	request := &dialogflowpb.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         message,
					LanguageCode: "us",
				},
			},
		},
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone: "Asia/Bangkok",
			Payload: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"source": {
						Kind: &structpb.Value_StringValue{
							StringValue: "LINE",
						},
					},
					"platform": {
						Kind: &structpb.Value_StringValue{
							StringValue: "LINE",
						},
					},
				},
			},
		},
	}

	response, err := client.DetectIntent(ctx, request)

	parameters := make(entities.Parameters)
	for _, p := range response.QueryResult.Intent.Parameters {
		parameters[p.Name] = p.Value
	}

	intent := &entities.Intent{
		Name:       response.QueryResult.Action,
		Parameters: parameters,
	}

	fmt.Println(intent)
	fmt.Println(err)

	requestJsonByte, _ := json.MarshalIndent(request, "", "  ")
	fmt.Println(string(requestJsonByte))

	responseJsonByte, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(responseJsonByte))

	fmt.Println(response.QueryResult.FulfillmentText)
	fmt.Println(response.QueryResult.QueryText)
	fmt.Println(response.QueryResult.Action)
	fmt.Println(response.QueryResult.Intent)
	fmt.Println(response.QueryResult.Intent.Action)
	fmt.Println(response.QueryResult.Intent.Parameters)
	fmt.Println(response.QueryResult.Intent.FollowupIntentInfo)
}
