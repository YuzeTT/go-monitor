package hook

import (
	"github.com/sirupsen/logrus"
	"hsott.cn/next-monitor/web"
)

type defaultHook struct {
}

func InitHook(log *logrus.Logger) {
	log.AddHook(&defaultHook{})
}

func (h *defaultHook) Fire(entry *logrus.Entry) error {
	web.AddLog(entry.Message, entry.Time)
	return nil
}
func (h *defaultHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
