package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetUserInfo(username string) (user UserInfo, err error) {

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s", username), nil)
	if err != nil {
		return user, err
	}

	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		return user, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return user, err
		}

		return user, nil
	} else {
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return user, err
		}

		return user, errors.New(string(responseData))
	}
}