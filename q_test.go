package q_test

import (
	. "bitbucket.org/searchdiscovery/q"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Q", func() {
	Describe("Shift", func() {

		Context("When a single message is in the queue", func() {
			It("Should return a single message containing the correct message body", func() {
				message := "Test Message!"
				sendTestMessage(message, queueURL)

				messages, err := queue.Shift()

				Expect(err).To(BeNil())
				Expect(len(messages)).To(Equal(1))

				for _, msg := range messages {
					Expect(msg.Body).To(Equal(message))
				}
			})
		})

		Context("When a request is made to an empty queue", func() {
			It("Should wait for a message then receive it when sent", func() {

				ci := make(chan []*Message)

				go func() {
					messages, _ := queue.Shift()
					ci <- messages
				}()

				time.Sleep(500 * time.Millisecond)

				message := "Another test message!"
				sendTestMessage(message, queueURL)

				messages := <-ci

				for _, msg := range messages {
					Expect(msg.Body).To(Equal(message))
				}
			})
		})
	})

	Describe("Polling", func() {
		Context("When there are messages added to queue", func() {
			It("Should pick them up, continue polling", func() {

				msgCh := make(chan string)
				errCh := make(chan error)

				go queue.Poll(msgCh, errCh)

				message1 := "Polling test #1"
				message2 := "Aaand polling test #2"

				sendTestMessage(message1, queueURL)

				Expect(<-msgCh).To(Equal(message1))

				time.Sleep(500 * time.Millisecond)

				sendTestMessage(message2, queueURL)

				Expect(<-msgCh).To(Equal(message2))

				var err error

				select {
				case err = <-errCh:
				default:
					// do nothing
				}

				Expect(err).To(BeNil())
			})
		})
	})
})
