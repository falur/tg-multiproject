package bot

import tele "gopkg.in/telebot.v3"

func authMiddleware(allowedUserID int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID != allowedUserID {
				return c.Send("Доступ запрещён.")
			}
			return next(c)
		}
	}
}
