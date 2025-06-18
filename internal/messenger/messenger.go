package messenger

import (
	"friend_graphql/internal/logger"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"strings"
	"sync"
)

type Messenger struct {
	m       *melody.Melody
	mu      *sync.Mutex
	clients map[string]*melody.Session
}

type Message struct {
	//todo: парсить в структуру ?
}

func MessageConvert(msg []byte) string {
	//todo: сообщение типа кто/кому/что/когда--> ?
	messageSplit := strings.Split(string(msg), "/")
	return messageSplit[1]
}

func NewMessangerDomain() *Messenger {
	return &Messenger{m: melody.New(), mu: &sync.Mutex{}, clients: make(map[string]*melody.Session, 30)}
}

func (m *Messenger) MessangerHander() echo.HandlerFunc {
	return func(c echo.Context) error {
		op := "upgrade connection"
		messengerLogger := logger.GetLogger().With(op)
		err := m.m.HandleRequest(c.Response(), c.Request())
		if err != nil {
			messengerLogger.Error("upgrade connection err:", "err", err)
			return err
		}
		m.m.HandleConnect(func(s *melody.Session) {
			userID, ok := c.Get("userID").(string)
			if !ok {
				messengerLogger.Error("problem to convert userID data:", "userID", userID)
				if err = s.Close(); err != nil {
					messengerLogger.Error("problem to close session:", "melody err", err.Error())
				}
			} else {
				main, ok := c.Get("main").(string)
				if ok && main == "main" {
					s.Set("userID", userID)
					m.mu.Lock()
					m.clients[userID] = s
					m.mu.Unlock()
				} else if ok && main != "main" {
					//todo: можно прописать чтобы при отключении все пользователи которые пренадлежали этому пользователю отключались
				} else {
					if err = s.Close(); err != nil {
						messengerLogger.Error("problem to close session", "melody err", err.Error())
					}
				}

				//todo:написать логику отправки сообщения что пользователь подключился на сервер

			}
		})

		m.m.HandleMessage(func(s *melody.Session, msg []byte) {
			reader := MessageConvert(msg)
			client, ok := m.clients[reader]
			if !ok {
				//todo:реализовать отправку cообщения на kafka
			} else {
				if err := client.Write(msg); err != nil {
					messengerLogger.Error("problem to write msg:", "melody err", err.Error())
					//todo:реализовать отправку cообщения на kafka для оповещения и может удаление пользователя
				}
			}

		})

		m.m.HandleDisconnect(func(s *melody.Session) {
			userID, _ := s.Get("userID")
			m.mu.Lock()
			if err := m.clients[userID.(string)].Close(); err != nil {
				messengerLogger.Error("problem to close session:", "melody err", err.Error())
			}
			m.mu.Unlock()
		})
		return err
	}
}
