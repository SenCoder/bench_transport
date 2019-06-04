package util

import (
	"io/ioutil"
	"os"
	"strings"
)

func IsExits(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func ToTrimString(filePath string) (string, error) {
	str, err := ToString(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

func ToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
