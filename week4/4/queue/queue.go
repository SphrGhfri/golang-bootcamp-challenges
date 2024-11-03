package queue

import (
	"container/heap"

	"github.com/google/uuid"
)

type IMessageQueue interface {
	heap.Interface
	PushMessage(content string, priority int)
	PopMessage() *Message
}

type Message struct {
	ID       uuid.UUID
	Content  string
	Priority int
	Index    int
}

type MessageQueue []*Message

func (mq MessageQueue) Len() int { return len(mq) }

func (mq MessageQueue) Less(i, j int) bool {
	return mq[i].Priority < mq[j].Priority
}

func (mq MessageQueue) Swap(i, j int) {
	mq[i], mq[j] = mq[j], mq[i]
	mq[i].Index = i
	mq[j].Index = j
}

func (mq *MessageQueue) Push(x interface{}) {
	n := len(*mq)
	message := x.(*Message)
	message.Index = n
	*mq = append(*mq, message)
}

func (mq *MessageQueue) Pop() interface{} {
	old := *mq
	n := len(old)
	message := old[n-1]
	message.Index = -1
	*mq = old[0 : n-1]
	return message
}

func (mq *MessageQueue) PushMessage(content string, priority int) {
	message := &Message{
		ID:       uuid.New(),
		Content:  content,
		Priority: priority,
	}
	heap.Push(mq, message)
}

func (mq *MessageQueue) PopMessage() *Message {
	if mq.Len() == 0 {
		return nil
	}
	return heap.Pop(mq).(*Message)
}

func NewMessageQueue() IMessageQueue {
	mq := &MessageQueue{}
	heap.Init(mq)
	return mq
}