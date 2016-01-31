package pubsub

import (
	"encoding/base64"
	"log"
	"reflect"
	"testing"

	"github.com/xhroot/gizmo/config"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/protobuf/proto"
)

func TestSQSSubscriber(t *testing.T) {
	test1 := "hey hey hey!"
	test2 := "ho ho ho!"
	test3 := "yessir!"
	test4 := "nope!"
	sqstest := &TestSQSAPI{
		Messages: [][]*sqs.Message{
			[]*sqs.Message{
				&sqs.Message{
					Body:          &test1,
					ReceiptHandle: &test1,
				},
				&sqs.Message{
					Body:          &test2,
					ReceiptHandle: &test2,
				},
			},
			[]*sqs.Message{
				&sqs.Message{
					Body:          &test3,
					ReceiptHandle: &test3,
				},
				&sqs.Message{
					Body:          &test4,
					ReceiptHandle: &test4,
				},
			},
		},
	}

	cfg := &config.SQS{}
	defaultSQSConfig(cfg)
	sub := &SQSSubscriber{
		sqs:        sqstest,
		cfg:        cfg,
		toDelete:   make(chan *deleteRequest),
		deleteDone: make(chan bool),
		stop:       make(chan chan error, 1),
	}

	queue := sub.Start()
	verifySQSSub(t, queue, sqstest, test1, 0)
	verifySQSSub(t, queue, sqstest, test2, 1)
	verifySQSSub(t, queue, sqstest, test3, 2)
	verifySQSSub(t, queue, sqstest, test4, 3)
	sub.Stop()
}

func verifySQSSub(t *testing.T, queue <-chan SubscriberMessage, testsqs *TestSQSAPI, want string, index int) {
	gotRaw := <-queue
	got := string(gotRaw.Message())
	if got != want {
		t.Errorf("SQSSubscriber expected:\n%#v\ngot:\n%#v", want, got)
	}
	gotRaw.Done()

	if len(testsqs.Deleted) != (index + 1) {
		t.Errorf("SQSSubscriber expected %d deleted message, got: %d", index+1, len(testsqs.Deleted))
	}

	if *testsqs.Deleted[index].ReceiptHandle != want {
		t.Errorf("SQSSubscriber expected receipt handle of \"%s\" , got: \"%s\"",
			want,
			*testsqs.Deleted[index].ReceiptHandle)
	}
}

func TestSQSSubscriberProto(t *testing.T) {
	test1 := &TestProto{"hey hey hey!"}
	test2 := &TestProto{"ho ho ho!"}
	test3 := &TestProto{"yessir!"}
	test4 := &TestProto{"nope!"}
	sqstest := &TestSQSAPI{
		Messages: [][]*sqs.Message{
			[]*sqs.Message{
				&sqs.Message{
					Body:          makeB64String(test1),
					ReceiptHandle: &test1.Value,
				},
				&sqs.Message{
					Body:          makeB64String(test2),
					ReceiptHandle: &test2.Value,
				},
			},
			[]*sqs.Message{
				&sqs.Message{
					Body:          makeB64String(test3),
					ReceiptHandle: &test3.Value,
				},
				&sqs.Message{
					Body:          makeB64String(test4),
					ReceiptHandle: &test4.Value,
				},
			},
		},
	}
	cfg := &config.SQS{ConsumeProtobuf: true}
	defaultSQSConfig(cfg)
	sub := &SQSSubscriber{
		sqs:        sqstest,
		cfg:        cfg,
		toDelete:   make(chan *deleteRequest),
		deleteDone: make(chan bool),
		stop:       make(chan chan error, 1),
	}

	queue := sub.Start()

	verifySQSSubProto(t, queue, sqstest, test1, 0)
	verifySQSSubProto(t, queue, sqstest, test2, 1)
	verifySQSSubProto(t, queue, sqstest, test3, 2)
	verifySQSSubProto(t, queue, sqstest, test4, 3)

	sub.Stop()
}

func verifySQSSubProto(t *testing.T, queue <-chan SubscriberMessage, testsqs *TestSQSAPI, want *TestProto, index int) {
	gotRaw := <-queue
	got := makeProto(gotRaw.Message())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SQSSubscriber expected:\n%#v\ngot:\n%#v", want, got)
	}
	gotRaw.Done()

	if len(testsqs.Deleted) != (index + 1) {
		t.Errorf("SQSSubscriber expected %d deleted message, got: %d", index+1, len(testsqs.Deleted))
	}

	if *testsqs.Deleted[index].ReceiptHandle != want.Value {
		t.Errorf("SQSSubscriber expected receipt handle of \"%s\" , got: \"%s\"",
			want.Value,
			*testsqs.Deleted[index].ReceiptHandle)
	}
}

func makeB64String(p proto.Message) *string {
	b, _ := proto.Marshal(p)
	s := base64.StdEncoding.EncodeToString(b)
	return &s
}

func makeProto(b []byte) *TestProto {
	t := &TestProto{}
	err := proto.Unmarshal(b, t)
	if err != nil {
		log.Printf("unable to unmarshal protobuf: %s", err)
	}
	return t
}

type TestSQSAPI struct {
	Offset   int
	Messages [][]*sqs.Message
	Deleted  []*sqs.DeleteMessageBatchRequestEntry
}

func (s *TestSQSAPI) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if s.Offset >= len(s.Messages) {
		return &sqs.ReceiveMessageOutput{}, nil
	}
	out := s.Messages[s.Offset]
	s.Offset++
	return &sqs.ReceiveMessageOutput{Messages: out}, nil
}

func (s *TestSQSAPI) DeleteMessageBatch(i *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	s.Deleted = append(s.Deleted, i.Entries...)
	return nil, errNotImpl
}

///////////
// ALL METHODS BELOW HERE ARE EMPTY AND JUST SATISFYING THE SQSAPI interface
///////////

func (s *TestSQSAPI) DeleteMessage(d *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, errNotImpl
}

func (s *TestSQSAPI) DeleteMessageBatchRequest(i *sqs.DeleteMessageBatchInput) (*request.Request, *sqs.DeleteMessageBatchOutput) {
	return nil, nil
}

func (s *TestSQSAPI) AddPermissionRequest(*sqs.AddPermissionInput) (*request.Request, *sqs.AddPermissionOutput) {
	return nil, nil
}
func (s *TestSQSAPI) AddPermission(*sqs.AddPermissionInput) (*sqs.AddPermissionOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) ChangeMessageVisibilityRequest(*sqs.ChangeMessageVisibilityInput) (*request.Request, *sqs.ChangeMessageVisibilityOutput) {
	return nil, nil
}
func (s *TestSQSAPI) ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) ChangeMessageVisibilityBatchRequest(*sqs.ChangeMessageVisibilityBatchInput) (*request.Request, *sqs.ChangeMessageVisibilityBatchOutput) {
	return nil, nil
}
func (s *TestSQSAPI) ChangeMessageVisibilityBatch(*sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) CreateQueueRequest(*sqs.CreateQueueInput) (*request.Request, *sqs.CreateQueueOutput) {
	return nil, nil
}
func (s *TestSQSAPI) CreateQueue(*sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) DeleteMessageRequest(*sqs.DeleteMessageInput) (*request.Request, *sqs.DeleteMessageOutput) {
	return nil, nil
}

func (s *TestSQSAPI) DeleteQueueRequest(*sqs.DeleteQueueInput) (*request.Request, *sqs.DeleteQueueOutput) {
	return nil, nil
}
func (s *TestSQSAPI) DeleteQueue(*sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) GetQueueAttributesRequest(*sqs.GetQueueAttributesInput) (*request.Request, *sqs.GetQueueAttributesOutput) {
	return nil, nil
}
func (s *TestSQSAPI) GetQueueAttributes(*sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) GetQueueUrlRequest(*sqs.GetQueueUrlInput) (*request.Request, *sqs.GetQueueUrlOutput) {
	return nil, nil
}
func (s *TestSQSAPI) GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) ListDeadLetterSourceQueuesRequest(*sqs.ListDeadLetterSourceQueuesInput) (*request.Request, *sqs.ListDeadLetterSourceQueuesOutput) {
	return nil, nil
}
func (s *TestSQSAPI) ListDeadLetterSourceQueues(*sqs.ListDeadLetterSourceQueuesInput) (*sqs.ListDeadLetterSourceQueuesOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) ListQueuesRequest(*sqs.ListQueuesInput) (*request.Request, *sqs.ListQueuesOutput) {
	return nil, nil
}
func (s *TestSQSAPI) ListQueues(*sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) PurgeQueueRequest(*sqs.PurgeQueueInput) (*request.Request, *sqs.PurgeQueueOutput) {
	return nil, nil
}
func (s *TestSQSAPI) PurgeQueue(*sqs.PurgeQueueInput) (*sqs.PurgeQueueOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) ReceiveMessageRequest(*sqs.ReceiveMessageInput) (*request.Request, *sqs.ReceiveMessageOutput) {
	return nil, nil
}

func (s *TestSQSAPI) RemovePermissionRequest(*sqs.RemovePermissionInput) (*request.Request, *sqs.RemovePermissionOutput) {
	return nil, nil
}
func (s *TestSQSAPI) RemovePermission(*sqs.RemovePermissionInput) (*sqs.RemovePermissionOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) SendMessageRequest(*sqs.SendMessageInput) (*request.Request, *sqs.SendMessageOutput) {
	return nil, nil
}
func (s *TestSQSAPI) SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) SendMessageBatchRequest(*sqs.SendMessageBatchInput) (*request.Request, *sqs.SendMessageBatchOutput) {
	return nil, nil
}
func (s *TestSQSAPI) SendMessageBatch(*sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	return nil, errNotImpl
}
func (s *TestSQSAPI) SetQueueAttributesRequest(*sqs.SetQueueAttributesInput) (*request.Request, *sqs.SetQueueAttributesOutput) {
	return nil, nil
}
func (s *TestSQSAPI) SetQueueAttributes(*sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error) {
	return nil, errNotImpl
}
