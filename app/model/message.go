package model

type MessageContent struct {
	ID            int64
	Content       string
	CreateTime    int64
	OriginContent string
}

func (m *MessageContent) TableName() string {

	return "message_content"
}
