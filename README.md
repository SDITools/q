# Q

Tiny Go library for receiving messages from AWS Simple Queue Service.

A wrapper for the [AWS Labs SQS library](https://github.com/awslabs/aws-sdk-go/tree/master/service/sqs), Q provides facilities for "shifting" messages and polling.


## Usage

```Go
import (
  "fmt"
  "log"
  "github.com/SDITools/q"
)

var queue = q.New("[your queue url]", "us-east-1", q.QueueParams{})

var msgCh = make(chan string)
var errCh = make(chan error)

func main() {
  go queue.Poll(msgCh, errCh)

  var msg string
  var err error

  for {
    select {
    case msg = <-msgCh:
      fmt.Println(msg)
    case err = <-err:
      log.Error(err)
    }
  }
}
```
