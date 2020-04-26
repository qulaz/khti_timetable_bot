package vk

type MessageNew struct {
	Message    *Message   `json:"message"`
	ClientInfo ClientInfo `json:"client_info"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type HandlerMessageNew func(b *Bot, u *MessageNew)

func (u MessageNew) GetType() string {
	return u.UpdateMeta.UpdateType
}

func (u *MessageNew) SetUpdateMeta(meta *UpdateMeta) {
	u.UpdateMeta = meta
}
