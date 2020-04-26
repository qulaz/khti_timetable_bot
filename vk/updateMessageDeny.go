package vk

type MessageDeny struct {
	UserID int `json:"user_id"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type HandlerMessageDeny func(b *Bot, u *MessageDeny)

func (u MessageDeny) GetType() string {
	return u.UpdateMeta.UpdateType
}

func (u *MessageDeny) SetUpdateMeta(meta *UpdateMeta) {
	u.UpdateMeta = meta
}
