package common

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"os"
	"time"
)

func Run() {
	terminationTime := WaitUntilTermination()
	instanceId := GetInstanceId()
	TriggerActions(instanceId, terminationTime)
	time.Sleep(3 * time.Minute) // 3 minutes ought to be enough for anyone, right?
}

func WaitUntilTermination() string {
	sess := session.Must(session.NewSession())
	meta := ec2metadata.New(sess)

	for {
		time.Sleep(3 * time.Second)
		resp, err := meta.GetDynamicData("spot/termination-time")
		if err == nil {
			return resp
		}
	}
}

func GetInstanceId() string {
	sess := session.Must(session.NewSession())
	meta := ec2metadata.New(sess)

	identity, _ := meta.GetInstanceIdentityDocument()
	return identity.InstanceID
}

func TriggerActions(instanceId, terminationTime string) {
	// set env vars before executing any of the shell commands
	os.Setenv("TERMINATION_TIME", terminationTime)
	os.Setenv("INSTANCE_ID", instanceId)

	for _, action := range DefaultActionsList().Actions {
		go func() { action.Perform(instanceId, terminationTime) }()
		time.Sleep(1 * time.Second) // just to avoid a possible explosion of events at once
	}
}
