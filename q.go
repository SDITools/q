package q

// Tiny abstraction library for consuming/polling AWS SQS messages

import (
	"errors"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"
)

const defaultMaxMessages = 10
const defaultVisibilityTimeout = 1
const defaultWaitTimeSeconds = 20

var ErrNotFound = errors.New("No messages found")

type Queue struct {
	endpoint string
	params   QueueParams
	sqs      *sqs.SQS
	poll     bool
}

type QueueParams struct {
	MaxMessages       int64
	VisibilityTimeout int64
	WaitTimeSeconds   int64

	receiveMessageInput *sqs.ReceiveMessageInput
}

type Message struct {
	ID            *string
	ReceiptHandle *string
	Body          string
}

func (qq *Queue) Receive() (msgs []*Message, err error) {
	ms, err := qq.sqs.ReceiveMessage(qq.params.receiveMessageInput)

	for _, m := range ms.Messages {
		msgs = append(msgs, &Message{
			ID:            m.MessageID,
			ReceiptHandle: m.ReceiptHandle,
			Body:          *m.Body,
		})
	}

	if len(msgs) == 0 && err == nil {
		return msgs, ErrNotFound
	}
	return msgs, err
}

func (qq *Queue) Delete(msgs []*Message) error {
	input := &sqs.DeleteMessageBatchInput{
		QueueURL: aws.String(qq.endpoint),
		Entries:  make([]*sqs.DeleteMessageBatchRequestEntry, 0),
	}
	for _, msg := range msgs {
		input.Entries = append(input.Entries, &sqs.DeleteMessageBatchRequestEntry{
			ID:            msg.ID,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}
	_, err := qq.sqs.DeleteMessageBatch(input)
	return err
}

func (qq *Queue) Shift() ([]*Message, error) {
	msgs, err := qq.Receive()
	if err != nil {
		return msgs, err
	}

	delErr := qq.Delete(msgs)
	return msgs, delErr
}

func (qq *Queue) Poll(msgCh chan<- string, errCh chan<- error) {
	qq.poll = true

	for {
		if !qq.poll {
			break
		}

		msgs, err := qq.Shift()

		if len(msgs) > 0 {
			for _, msg := range msgs {
				msgCh <- msg.Body
			}
		}

		if err != nil && err != ErrNotFound {
			errCh <- err
		}
	}
}

func (qq *Queue) StopPoll() {
	qq.poll = false
}

func (qp *QueueParams) defaults(endpoint string) {
	if qp.MaxMessages == 0 {
		qp.MaxMessages = defaultMaxMessages
	}
	if qp.VisibilityTimeout == 0 {
		qp.VisibilityTimeout = defaultVisibilityTimeout
	}
	if qp.WaitTimeSeconds == 0 {
		qp.WaitTimeSeconds = defaultWaitTimeSeconds
	}

	qp.receiveMessageInput = &sqs.ReceiveMessageInput{
		QueueURL:            aws.String(endpoint),
		MaxNumberOfMessages: aws.Long(qp.MaxMessages),
		VisibilityTimeout:   aws.Long(qp.VisibilityTimeout),
		WaitTimeSeconds:     aws.Long(qp.WaitTimeSeconds),
	}
}

func New(endpoint string, region string, params QueueParams) *Queue {
	params.defaults(endpoint)
	return &Queue{
		sqs:      sqs.New(&aws.Config{Region: region}),
		params:   params,
		endpoint: endpoint,
	}
}
