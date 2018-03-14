package response

type MessageList []string

func NewMessageList(msgs ...string) MessageList {
	return MessageList(msgs)
}

func (ml *MessageList) Append(msgs ...string) {
	*ml = append(*ml, msgs...)
}
