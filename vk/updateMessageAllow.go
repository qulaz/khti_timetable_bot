package vk

type MessageAllow struct {
	UserID int    `json:"user_id"`
	Key    string `json:"key"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type HandlerMessageAllow func(b *Bot, u *MessageAllow)

func (u MessageAllow) GetType() string {
	return u.UpdateMeta.UpdateType
}

func (u *MessageAllow) SetUpdateMeta(meta *UpdateMeta) {
	u.UpdateMeta = meta
}
