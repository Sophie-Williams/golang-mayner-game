package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	HOUR        = 2
	MONEY_VALUE = 10
)

var (
	insert_new_user            = "INSERT INTO users (name, mayner1, mayner2, mayner3, mayner4, score, money, time) VALUES (?, 1, 0, 0, 0, 0, 300, ?);"
	find_user                  = "SELECT id FROM users WHERE name=?"
	find_score_by_username     = "SELECT score, time FROM users WHERE name=?"
	update_user_score_and_time = "UPDATE users SET score=?, time=? WHERE name=?"
	get_all_video              = "SELECT mayner1, mayner2, mayner3, mayner4 FROM users WHERE name=?"
	get_new_money              = "SELECT money, score FROM users WHERE name=?"
)

func renderScore(username string) (money int64) {
	row := db.QueryRow(
		find_score_by_username,
		username,
	)

	var clock int64

	err := row.Scan(&money, &clock)
	if err != nil {
		err.Error()
	}

	var videocarts [4]int
	col := db.QueryRow(
		get_all_video,
		username,
	)

	err = col.Scan(
		&videocarts[0],
		&videocarts[1],
		&videocarts[2],
		&videocarts[3],
	)
	if err != nil {
		err.Error()
	}

	timeBefore := clock
	timeNow := time.Now().Unix()

	for i, el := range videos {
		money += (timeNow - timeBefore) * int64(videocarts[i]*el.Power/HOUR)
	}

	_, err = db.Exec(
		update_user_score_and_time,
		money,
		time.Now().Unix(),
		username,
	)
	return
}

func menu(msg *tgbotapi.Message) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/score"),
			tgbotapi.NewKeyboardButton("/video"),
			tgbotapi.NewKeyboardButton("/shop"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/donate"),
			tgbotapi.NewKeyboardButton("/help"),
			tgbotapi.NewKeyboardButton("/sell"),
		),
	)
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Menu")
	reply.ReplyMarkup = keyboard
	_, err := bot.Send(reply)
	if err != nil {
		err.Error()
	}

}

func video(msg *tgbotapi.Message) {
	var videos [4]int

	row := db.QueryRow(
		get_all_video,
		msg.From.UserName,
	)

	err := row.Scan(
		&videos[0],
		&videos[1],
		&videos[2],
		&videos[3],
	)

	if err != nil {
		err.Error()
	}

	fmt.Println(videos)

	for i, el := range videos {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Количество видеокарт %d - %d", i+1, el))
		_, err := bot.Send(reply)
		if err != nil {
			err.Error()
		}
	}

}

func start(msg *tgbotapi.Message) {
	var reply tgbotapi.MessageConfig

	var res sql.NullString
	row := db.QueryRow(find_user, msg.From.UserName)

	err := row.Scan(&res)
	if err != nil {
		err.Error()
	}
	if res.Valid {
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Ты уже зарегистрирован")
	} else {
		_, err := db.Exec(
			insert_new_user,
			msg.From.UserName,
			time.Now().Unix(),
		)
		if err != nil {
			err.Error()
		}
		reply = tgbotapi.NewMessage(msg.Chat.ID, "Ты регнулся! /help")
	}

	_, err = bot.Send(reply)
	if err != nil {
		err.Error()
	}
}

func sell(msg *tgbotapi.Message) {
	row := db.QueryRow(get_new_money, msg.From.UserName)
	var money, score int64
	_ = renderScore(msg.From.UserName)
	err := row.Scan(&money, &score)
	if err != nil {
		err.Error()
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Сейчас у тебя %dР\nОбменять можно 500b -> 1Р, от 500b\nБаланс: %db", money, score))

	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продать!", "yes"),
		),
	)
	_, err = bot.Send(reply)
	if err != nil {
		err.Error()
	}

}

func score(msg *tgbotapi.Message) {
	money := renderScore(msg.Chat.UserName)
	reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Твой баланс %d Bitcoins!", money))

	_, err := bot.Send(reply)
	if err != nil {
		err.Error()
	}
}

func donate(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Нужно больше Р, чтобы купить видеокарт?\nПиши @likipiki.\nИ за небольшое пожертвование получи бонусы!\nТак же туда можно присылать отзывы и предложения!")
	_, err := bot.Send(reply)
	if err != nil {
		err.Error()
	}
}

func shop(msg *tgbotapi.Message) {
	var reply tgbotapi.MessageConfig
	for i, el := range videos {
		reply = tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("%s\n%s", el.Name, el.Desk))

		reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Купить %dР", el.Cost), fmt.Sprintf("video %d", i)),
			),
		)
		_, err := bot.Send(reply)
		if err != nil {
			err.Error()
		}
	}
}

func help(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, helpDesc)
	_, err := bot.Send(reply)
	if err != nil {
		err.Error()
	}
}
