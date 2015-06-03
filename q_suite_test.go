package q_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sditools/q"
)

const REGION = "us-east-1"
const REDIS_ADDR = "localhost:6379"

var svc *sqs.SQS
var queueURL string
var queue *Queue

var _ = BeforeSuite(func() {

	svc = sqs.New(&aws.Config{
		Credentials: aws.DefaultChainCredentials,
		Region:      *aws.String(REGION),
	})

	timestamp := time.Now().Local().Format("20060102150405")
	resp, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("test-queue-" + timestamp),
	})

	if err != nil {
		fmt.Println(err)
	}

	queueURL = *resp.QueueURL

	queue = New(queueURL, REGION, QueueParams{})
})

var _ = AfterSuite(func() {
	_, err := svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueURL: aws.String(queueURL),
	})

	if err != nil {
		fmt.Println(err)
	}
})

func sendTestMessage(message string, queueURL string) {
	svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueURL:    aws.String(queueURL),
	})
}

func TestQ(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Q Suite")
}
