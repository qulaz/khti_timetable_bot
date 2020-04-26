package vk

import (
	"time"
)

// Middleware позволяет контролировать процесс выполнения обработки апдейта и выполнять определенные коллбеки на каждой
// из стадии его обработки. С их помощью можно настроить логгирование, фильтрацию, ограничение доступа и многое другое.
//
// Порядок выполнения функций Middleware следующий:
//
// 	> OnPreProcessUpdate
// 	> OnPre<UpdateType>Update (OnPreMessageNewUpdate, OnPreMessageReplyUpdate и другие в зависимости от пришедшего
// 	апдейта)
// 	> OnPost<UpdateType>Update
// 	> OnPostProcessUpdate
//
// Каждая из функций возвращает булевое значение. В случае, если любая из OnPre функций вернет false, обработка апдейта
// останавливается. В OnPost функциях на данный момент возвращаемое значение ни на что не влияет.
type Middleware struct {
	OnPreProcessUpdate       func(b *Bot, u *Update) bool
	OnPostProcessUpdate      func(b *Bot, u *Update) bool
	OnPreMessageNewUpdate    func(b *Bot, u *MessageNew) bool
	OnPostMessageNewUpdate   func(b *Bot, u *MessageNew) bool
	OnPreMessageReplyUpdate  func(b *Bot, u *MessageReply) bool
	OnPostMessageReplyUpdate func(b *Bot, u *MessageReply) bool
	OnPreMessageEditUpdate   func(b *Bot, u *MessageEdit) bool
	OnPostMessageEditUpdate  func(b *Bot, u *MessageEdit) bool
	OnPreMessageAllowUpdate  func(b *Bot, u *MessageAllow) bool
	OnPostMessageAllowUpdate func(b *Bot, u *MessageAllow) bool
	OnPreMessageDenyUpdate   func(b *Bot, u *MessageDeny) bool
	OnPostMessageDenyUpdate  func(b *Bot, u *MessageDeny) bool
}

const (
	onPreProcessMiddleware       = "OnPreProcessUpdate"
	onPostProcessMiddleware      = "OnPostProcessUpdate"
	onPreMessageNewMiddleware    = "OnPreMessageNewUpdate"
	onPostMessageNewMiddleware   = "OnPostMessageNewUpdate"
	onPreMessageReplyMiddleware  = "OnPreMessageReplyUpdate"
	onPostMessageReplyMiddleware = "OnPostMessageReplyUpdate"
	onPreMessageEditMiddleware   = "OnPreMessageEditUpdate"
	onPostMessageEditMiddleware  = "OnPostMessageEditUpdate"
	onPreMessageAllowMiddleware  = "OnPreMessageAllowUpdate"
	onPostMessageAllowMiddleware = "OnPostMessageAllowUpdate"
	onPreMessageDenyMiddleware   = "OnPreMessageDenyUpdate"
	onPostMessageDenyMiddleware  = "OnPostMessageDenyUpdate"
)

var logUpdateHandleStatus = map[bool]string{true: "Handled", false: "Unhandled"}

// Middleware для логгирования всех апдейтов
var LoggingMiddleware = Middleware{
	OnPreProcessUpdate: func(b *Bot, u *Update) bool {
		b.Logger.Infof("Received update <ID: %s>", u.UpdateMeta.EventID)
		return true
	},
	OnPostProcessUpdate: func(b *Bot, u *Update) bool {
		duration := time.Now().Sub(u.UpdateMeta.ReceiveTime)
		b.Logger.Infof("Update <ID: %s> processed in %v ms", u.UpdateMeta.EventID, duration.Milliseconds())
		return true
	},
	OnPreMessageNewUpdate: func(b *Bot, u *MessageNew) bool {
		b.Logger.Infof(
			"Received message [ID: %d] from [VKID: %d (PeerID: %d)]",
			u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPostMessageNewUpdate: func(b *Bot, u *MessageNew) bool {
		b.Logger.Debugf(
			"%s new message [ID: %d] from [VKID: %d (PeerID: %d)]",
			logUpdateHandleStatus[u.UpdateMeta.Handled], u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPreMessageReplyUpdate: func(b *Bot, u *MessageReply) bool {
		b.Logger.Infof("Received reply message [ID: %d] to [VKID: %d (PeerID: %d)]",
			u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPostMessageReplyUpdate: func(b *Bot, u *MessageReply) bool {
		b.Logger.Infof("%s reply message [ID: %d] to [VKID: %d (PeerID: %d)]",
			logUpdateHandleStatus[u.UpdateMeta.Handled], u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPreMessageEditUpdate: func(b *Bot, u *MessageEdit) bool {
		b.Logger.Infof("Received edited message [ID: %d] in [VKID: %d (PeerID: %d)]",
			u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPostMessageEditUpdate: func(b *Bot, u *MessageEdit) bool {
		b.Logger.Infof("%s edited message [ID: %d] in [VKID: %d (PeerID: %d)]",
			logUpdateHandleStatus[u.UpdateMeta.Handled], u.Message.ID, u.Message.FromID, u.Message.PeerID,
		)
		return true
	},
	OnPreMessageAllowUpdate: func(b *Bot, u *MessageAllow) bool {
		b.Logger.Infof("Received message allow update from [VKID: %d]", u.UserID)
		return true
	},
	OnPostMessageAllowUpdate: func(b *Bot, u *MessageAllow) bool {
		b.Logger.Infof(
			"%s message allow update from [VKID: %d]", logUpdateHandleStatus[u.UpdateMeta.Handled], u.UserID,
		)
		return true
	},
	OnPreMessageDenyUpdate: func(b *Bot, u *MessageDeny) bool {
		b.Logger.Infof("Received message deny update from [VKID: %d]", u.UserID)
		return true
	},
	OnPostMessageDenyUpdate: func(b *Bot, u *MessageDeny) bool {
		b.Logger.Infof(
			"%s message deny update from [VKID: %d]", logUpdateHandleStatus[u.UpdateMeta.Handled], u.UserID,
		)
		return true
	},
}
