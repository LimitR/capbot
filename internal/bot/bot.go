package bot

import (
	"fmt"
	"strconv"
	"time"

	"capbot/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	storage map[int64]*user.User
	counter map[int64]int
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	// bot.Debug = true

	if err != nil {
		panic(err)
	}
	return &Bot{
		api:     bot,
		storage: make(map[int64]*user.User, 10),
		counter: make(map[int64]int, 10),
	}, nil
}

func (b *Bot) GetInlineKeyboard(userId int64) []tgbotapi.InlineKeyboardButton {
	res := make([]tgbotapi.InlineKeyboardButton, 0, 10)
	user := b.storage[userId]
	for key, value := range user.Nums {
		res = append(res, tgbotapi.NewInlineKeyboardButtonData("["+strconv.Itoa(value)+"]", key))
	}
	return res
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	u.AllowedUpdates = []string{"message", "chat_member"}
	for update := range updates {

		if update.Message != nil {
			chatMem := update.Message.NewChatMembers
			if len(chatMem) != 0 {
				for _, userMember := range chatMem {
					b.storage[userMember.ID] = user.NewUser(userMember.ID, update.FromChat().ID)
					b.counter[userMember.ID] = 2
					msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Добро пожаловать, @%s!\n\nПодтвердите, что вы не бот, у вас есть 2 попытки и 5 минут.\n\n"+b.storage[userMember.ID].GetString(), userMember.UserName))
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							b.GetInlineKeyboard(userMember.ID)...,
						),
					)
					m, _ := b.api.Send(msg)
					b.storage[userMember.ID].MessageId = m.MessageID
				}
			}
			_, ok := b.storage[update.Message.From.ID]
			if ok {
				dm := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
				b.api.Send(dm)
			}
		}
		if update.CallbackQuery != nil {
			usr, ok := b.storage[update.SentFrom().ID]
			if ok && update.CallbackQuery.From.ID == usr.Id {
				if usr.Validate(update.CallbackQuery.Data) {
					chatId := update.CallbackQuery.Message.Chat.ChatConfig().ChatID
					msg := tgbotapi.NewMessage(chatId, "Вы успешно прошли капчу, спасибо.")
					m, _ := b.api.Send(msg)
					delete(b.storage, update.CallbackQuery.From.ID)
					delete(b.counter, update.CallbackQuery.From.ID)
					go func(chatId int64, msgId int) {
						time.Sleep(time.Second * 5)
						dm := tgbotapi.NewDeleteMessage(chatId, msgId)
						b.api.Send(dm)
					}(chatId, m.MessageID)
					dm := tgbotapi.NewDeleteMessage(chatId, usr.MessageId)
					b.api.Send(dm)
				} else {
					b.counter[update.Message.From.ID] = b.counter[update.Message.From.ID] - 1
					if b.counter[update.Message.From.ID] < 1 {
						ban := tgbotapi.BanChatMemberConfig{
							ChatMemberConfig: tgbotapi.ChatMemberConfig{
								ChatID:             update.Message.Chat.ID,
								SuperGroupUsername: "",
								ChannelUsername:    "",
								UserID:             update.Message.From.ID,
							},
							UntilDate:      0,
							RevokeMessages: false,
						}
						b.api.Send(ban)
						dm := tgbotapi.NewDeleteMessage(usr.ChatId, usr.MessageId)
						b.api.Send(dm)
						delete(b.storage, update.Message.From.ID)
						delete(b.counter, update.Message.From.ID)
					}
				}
			}
		}

		//Delete user and banned
		for _, usr := range b.storage {
			if time.Now().UnixMilli()-usr.Time.UnixMilli() > 60_000*5 {
				ban := tgbotapi.BanChatMemberConfig{
					ChatMemberConfig: tgbotapi.ChatMemberConfig{
						ChatID:             usr.ChatId,
						SuperGroupUsername: "",
						ChannelUsername:    "",
						UserID:             usr.Id,
					},
					UntilDate:      0,
					RevokeMessages: false,
				}
				b.api.Send(ban)
				delete(b.storage, usr.Id)
				delete(b.counter, usr.Id)
			}
		}
	}
}
