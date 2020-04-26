package vk

type MessageEdit struct {
	Message *Message `json:"object"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type HandlerMessageEdit func(b *Bot, u *MessageEdit)

func (u MessageEdit) GetType() string {
	return u.UpdateMeta.UpdateType
}

func (u *MessageEdit) SetUpdateMeta(meta *UpdateMeta) {
	u.UpdateMeta = meta
}
