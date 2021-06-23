package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string) (string, error) {
	resp, err := http.Post(url, "application/octet-stream", nil)
	if err != nil {
		return "", fmt.Errorf("Make POST request to %s: %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read body from %s: %s", url, err)
	}
	return string(body), nil
}
