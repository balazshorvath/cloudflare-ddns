package ipify

import (
	"io/ioutil"
	"net/http"
)

func GetIp() (string, error) {
	ip, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(ip.Body)
	return string(bytes), err
}
