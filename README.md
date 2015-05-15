# Q

Tiny Go library for receiving messages from AWS Simple Queue Service.  Supports "shifting" messages and polling.


## Usage

```Go
import (
  "fmt"
  "log"
  "bitbucket.org/searchdiscovery/q"
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
