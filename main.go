package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	TG_TOKEN string
)

type ChatInfo struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

func sendMsg(chatId int, text string) {
	// resp
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		TG_TOKEN, chatId, url.QueryEscape(text)))
	if err != nil {
		fmt.Println(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respData))
}

func isBingo(predict string, target string) (string, bool, bool) {
	var text string
	win := false
	isValid := true
	for i:=0; i<4; i++ {
		for j:=i+1; j<4; j++ {
			if predict[i] == predict[j] {
				isValid = false
			}
		}
	}
	if _, err := strconv.Atoi(predict); err != nil {
		text = "Invalid input, not digital..."
		isValid = false
	} else if len(predict) != 4 {
		text = "Invalid digital length..."
		isValid = false
	} else if !isValid {
		text = "Invalid digital format..."
		isValid = false
	} else {
		reta := 0
		retb := 0

		for i:= 0; i< 4; i++ {
			if target[i] == predict[i] {
				reta++
			}
		}

		for i:=0; i<4; i++ {
			for j:=0; j<4; j++ {
				if target[i] == predict[j] {
					retb++
				}
			}
		}

		retb -= reta
		text = fmt.Sprintf("%dA%dB", reta, retb)
		if reta == 4 {
			text = "Bingo! You've got the right number!"
			win = true
		}
	}
	return text, isValid, win
}

func main() {
	fmt.Println("Start TG bot server")

	rand.Seed(time.Now().UnixNano())
	chars := []rune("0123456789")

	fmt.Println("Read TG token...")
	f, err := os.Open("TG_TOKEN")
	if err != nil {
		fmt.Println("Fail to read TG token")
		return
	}
	s, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Read token file failure")
		return
	}
	TG_TOKEN = strings.Trim(string(s), " \n")
	fmt.Printf("Read token success, token value: %s\n", TG_TOKEN)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	BinggoNumber := make(map[int]string)
	players := make(map[string]int)
	twoPlayers := false
	var turn int

	r.POST("/bingo", func(c *gin.Context) {
		var chatInfo ChatInfo
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		// var data map[string]interface{}
		if err != nil {
			fmt.Println(err)
		}
		// chatId := data["message"]["chat"]["id"].int64
		json.Unmarshal(jsonData, &chatInfo)
		fmt.Println(chatInfo)

		var text string

		if strings.Trim(chatInfo.Message.Text, " \n") == "/bingo" {
			text = "Let's play a bingo game! Please input the 4 digit (no repeat) number..."

			length := 4
			var b strings.Builder

			rand.Shuffle(len(chars), func(i, j int) { chars[i], chars[j] = chars[j], chars[i] })
			for i := 0; i < length; i++ {
				b.WriteRune(chars[i])
			}
			str := b.String()
			BinggoNumber[chatInfo.Message.Chat.ID] = str
			fmt.Printf("Target number is %s\n", str)
			sendMsg(chatInfo.Message.Chat.ID, text)
			twoPlayers = false
		} else if strings.Trim(chatInfo.Message.Text, " \n") == "/bingoa" {
			text = "Let's play a bingo game with 2 people, you are A! Waiting for B!"
			players["A"] = chatInfo.Message.Chat.ID
			fmt.Println(players["A"])

			length := 4
			var b strings.Builder

			rand.Shuffle(len(chars), func(i, j int) { chars[i], chars[j] = chars[j], chars[i] })
			for i := 0; i < length; i++ {
				b.WriteRune(chars[i])
			}
			str := b.String()
			BinggoNumber[chatInfo.Message.Chat.ID] = str
			turn = chatInfo.Message.Chat.ID
			twoPlayers = true
			sendMsg(players["A"], text)
		} else if strings.Trim(chatInfo.Message.Text, " \n") == "/bingob" {
			text = "Let's play a bingo game with 2 people, you are B!"
			players["B"] = chatInfo.Message.Chat.ID

			if idA, exist := players["A"]; !exist {
				text += "\n Error! Player A does not exist!"
			} else if players["A"] == players["B"] {
				text += "\n Error! Player A and B should be different!"
			} else {
				BinggoNumber[chatInfo.Message.Chat.ID] = BinggoNumber[idA]
			}

			sendMsg(players["B"], text)
			sendMsg(players["A"], "Let's start the game, A pls input the digit")
			sendMsg(players["B"], "Let's start the game, A pls input the digit")
			twoPlayers = true

		} else if target, exist := BinggoNumber[chatInfo.Message.Chat.ID]; exist {
			predict := strings.Trim(chatInfo.Message.Text, " \n")

			if twoPlayers && chatInfo.Message.Chat.ID != turn {
				sendMsg( chatInfo.Message.Chat.ID, "Please wait for your turn...")
			} else if twoPlayers {
				text, isValid, win := isBingo(predict, target)
				if !isValid {
					sendMsg(turn, text)
				}

				if isValid {
					if players["A"] == turn {
						sendMsg(players["A"], fmt.Sprintf("Player A predict: %s, result: %s", predict, text))
						sendMsg(players["B"], fmt.Sprintf("Player A predict: %s, result: %s", predict, text))
						turn = players["B"]
					} else {
						sendMsg(players["A"], fmt.Sprintf("Player B predict: %s, result: %s", predict, text))
						sendMsg(players["B"], fmt.Sprintf("Player B predict: %s, result: %s", predict, text))
						turn = players["A"]
					}
				}

				if win {
					sendMsg(turn, text)
					if players["A"] == turn {
						sendMsg(players["A"], "Sorry you lose!")
					} else {
						sendMsg(players["B"], "Sorry you lose!")
					}
					// sendMsg(players["A"], fmt.Sprintf("The correct number is: %s", target))
					// sendMsg(players["B"], fmt.Sprintf("The correct number is: %s", target))
				}
			} else {
				text, _, _ := isBingo(predict, target)
				sendMsg(chatInfo.Message.Chat.ID, text)
			}

		}


	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}