package common

import (
	"os/exec"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/lambda"
	"encoding/json"
)

type Action struct {
	typ string
	action string
}

func (a *Action) Perform(instanceId, terminationTime string) {
	switch a.typ {
	case "shell":
		performShell(a.action, instanceId, terminationTime)
	case "sns":
		performSns(a.action, instanceId, terminationTime)
	case "lambda":
		performLambda(a.action, instanceId, terminationTime)
	default:
		panic("unexpected action type")
	}
}

func performShell(cmd, instanceId, timeArg string) {
	// NOTE args are unused because env vars have already been set
	_ = exec.Command("bash", "-c", cmd).Wait()
}

func performSns(topic, instanceId, timeArg string) {
	sess := session.Must(session.NewSession())
	client := sns.New(sess)

	payloadBytes := jsonPayload(instanceId, timeArg)
	payloadStr := string(payloadBytes)

	client.Publish(&sns.PublishInput{
		TopicArn: &topic,
		Message: &payloadStr,
	})
}

func performLambda(function, instanceId, timeArg string) {
	sess := session.Must(session.NewSession())
	client := lambda.New(sess)

	client.Invoke(&lambda.InvokeInput{
		FunctionName: &function,
		Payload: jsonPayload(instanceId, timeArg),
	})
}

func jsonPayload(instanceId, timeArg string) []byte {
	payload := map[string]string {
		"InstanceId": instanceId,
		"TerminationTime": timeArg,
	}

	encoded, _ := json.Marshal(payload)
	return encoded
}
