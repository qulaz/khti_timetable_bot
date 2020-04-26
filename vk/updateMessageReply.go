package vk

type MessageReply struct {
	Message *Message `json:"object"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type HandlerMessageReply func(b *Bot, u *MessageReply)

func (u MessageReply) GetType() string {
	return u.UpdateMeta.UpdateType
}

func (u *MessageReply) SetUpdateMeta(meta *UpdateMeta) {
	u.UpdateMeta = meta
}
