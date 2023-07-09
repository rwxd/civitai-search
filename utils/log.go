package utils

import "github.com/sirupsen/logrus"

func InitLogrus(level string) error {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(parsedLevel)
	return nil
}
