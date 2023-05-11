package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type channels struct {
	Ok               bool      `json:"ok"`
	Channels         []channel `json:"channels"`
	RespondeMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
type channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var slackToken string = os.Args[1]
var channelStr string = os.Args[2]
var userId string = os.Args[3]

func main() {
	isFinished := false
	myChannels := []channel{}
	cursor := ""

	for !isFinished {
		currChannels, err := getChannels(cursor)
		if err != nil {
			fmt.Println(err)
			return
		}

		myChannels = append(myChannels, currChannels.Channels...)

		if currChannels.RespondeMetadata.NextCursor == "" {
			isFinished = true
			continue
		}
		cursor = currChannels.RespondeMetadata.NextCursor

	}

	for _, c := range myChannels {
		if strings.HasPrefix(c.Name, channelStr) {
			time.Sleep(1100 * time.Millisecond)
			err := inviteChannel(c.Id)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func getChannels(cursor string) (channels, error) {

	myChannels := channels{}

	url := "https://slack.com/api/conversations.list?cursor=" + cursor + "&exclude_archived=true&limit=9999&types=public_channel,%20private_channel,%20mpim,%20im"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return myChannels, err
	}
	req.Header.Add("Authorization", "Bearer "+slackToken)

	res, err := client.Do(req)
	if err != nil {
		return myChannels, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return myChannels, err
	}

	jsonErr := json.Unmarshal(body, &myChannels)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return myChannels, jsonErr
	}

	return myChannels, nil
}

func inviteChannel(chId string) error {
	url := "https://slack.com/api/conversations.invite?channel=" + chId + "&users=" + userId
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+slackToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
