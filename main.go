package main

import (
	"encoding/json"
	"flag"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
	"go.eqrx.net/mauzr/pkg/bme/bme680"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	ConfigData *Config
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

type WebHookInfo struct {
	Ok     bool `json:"ok"`
	Result struct {
		URL                  string `json:"url"`
		HasCustomCertificate bool   `json:"has_custom_certificate"`
		PendingUpdateCount   int    `json:"pending_update_count"`
		MaxConnections       int    `json:"max_connections"`
		IPAddress            string `json:"ip_address"`
	} `json:"result"`
}

func sendMsg(chatId int, text string) {
	// resp
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		ConfigData.TGToken, chatId, url.QueryEscape(text)))
	if err != nil {
		fmt.Println(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respData))
}

func sendMarkdown(chatId int, text string) {
	// resp
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&parse_mode=MarkdownV2",
		ConfigData.TGToken, chatId, url.QueryEscape(text)))
	if err != nil {
		fmt.Println(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respData))
}

func sendHTML(chatId int, text string) {
	// resp
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&parse_mode=HTML",
		ConfigData.TGToken, chatId, url.QueryEscape(text)))
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
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
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

		for i := 0; i < 4; i++ {
			if target[i] == predict[i] {
				reta++
			}
		}

		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
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

func sendFile(filePath string, chatId int, caption string) error {
	// Add client
	client := &http.Client{}
	//prepare the reader instances to encode
	fout, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	values := map[string]io.Reader{
		"document": fout, // lets assume its this file
		"chat_id":  strings.NewReader(strconv.Itoa(chatId)),
		"caption":  strings.NewReader(caption),
	}

	err = upload(client, fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", ConfigData.TGToken), values)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	fmt.Println("Start TG bot server")
	configFile := flag.String("config", "./config.json", "Config file path")
	flag.Parse()

	var err error
	ConfigData, err = loadConfig(*configFile)
	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("0123456789")

	fmt.Printf("Running server %s\n", ConfigData.Server)

	fmt.Println("Remove all previous webhooks...")
	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getWebhookInfo", ConfigData.TGToken))
	if err != nil {
		fmt.Println(err)
		return
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(respData))

	webHookInfo := WebHookInfo{}
	json.Unmarshal(respData, &webHookInfo)
	fmt.Printf("Current webhook info: %v", webHookInfo)
	resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook?url=%s",
		ConfigData.TGToken, webHookInfo.Result.URL))
	if err != nil {
		fmt.Println(err)
		return
	}
	respData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(respData))
	fmt.Println("Remove success, adding new webhook...")
	resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s",
		ConfigData.TGToken, ConfigData.Server))
	if err != nil {
		fmt.Println(err)
		return
	}
	respData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Register new webhook success...")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Static("/cache", "./cache")

	BinggoNumber := make(map[int]string)
	players := make(map[string]int)
	twoPlayers := false
	var turn int

	sensor := bme680.New("/dev/i2c-1", 0x77)

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

		spice := strings.Split("肉桂 丁香 茴香 陈皮 草果 豆蔻 鼠尾 香叶 甘草 百里香 孜然 香茅草 迷迭香", " ")
		cook := strings.Split("煎 炒 烹 炸 煮 熬 炖 溜 烧 汆 烤", " ")

		var text string

		if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/help") {
			sendHTML(chatInfo.Message.Chat.ID, `
Help:
####################
	NekoRoid is a bot having a lot of fun
<code>
 #    # ###### #    #  ####  
 ##   # #      #   #  #    # 
 # #  # #####  ####   #    # 
 #  # # #      #  #   #    # 
 #   ## #      #   #  #    # 
 #    # ###### #    #  ####
</code>
####################
	/cat -- post a random cute cat image
    /dice -- cast a dice
    /dart -- cast a dart
    /bingo -- play a bingo game
    /nhknews -- get nhk news from rss
    /investing -- get investing news from rss
    /googleai -- get google AI news
    /book {bookid} -- download a book
	/sensor -- current room status
	/search_song {song info} -- search a song and generate download link
	/weather {city} -- current weather in the city
    /stock {name} -- display stock market
####################
`)
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/stock") {
			// symbol := strings.Trim(strings.Replace(chatInfo.Message.Text, "/stock", "", -1), " \n")
			// 	csv, err := GetQuote(symbol, time.Now().Format("2006-01-02"))
			// 	if err != nil {
			// 		sendMsg(chatInfo.Message.Chat.ID, csv)
			// 	} else {
			// 		sendMsg(chatInfo.Message.Chat.ID, err.Error())
			// 	}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/weather") {
			city := strings.Trim(strings.Replace(chatInfo.Message.Text, "/weather", "", -1), " \n")
			sendMsg(chatInfo.Message.Chat.ID, GetWeather(city, ConfigData.WeatherToken))
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/aqi") {
			city := strings.Trim(strings.Replace(chatInfo.Message.Text, "/aqi", "", -1), " \n")
			sendMsg(chatInfo.Message.Chat.ID, GetAQI(city, ConfigData.WeatherToken))
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/nhknews") {
			for _, rss := range getFeed("http://www3.nhk.or.jp/rss/news/cat0.xml") {
				sendMsg(chatInfo.Message.Chat.ID, rss)
			}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/dw") {
			for _, rss := range getFeed("https://rss.dw.com/xml/rss-de-all") {
				sendMsg(chatInfo.Message.Chat.ID, rss)
			}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/stat") {
			// stat1, err := linuxproc.ReadStat("/proc/stat")
			// 	if err != nil {
			// 		sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("Error read cpu info %s", err.Error()))
			// 		return
			// 	}
			stat2, err := linuxproc.ReadMemInfo("/proc/meminfo")
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("Error read mem info %s", err.Error()))
				return
			}
			sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("Memory usage: %4.3f",
				float64(stat2.MemTotal-stat2.MemFree)/float64(stat2.MemTotal)))
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/googleai") {
			for _, rss := range getFeed("https://blog.google/technology/ai/rss/") {
				sendMsg(chatInfo.Message.Chat.ID, rss)
			}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/investing") {
			for _, rss := range getFeed("https://www.investing.com/rss/news.rss") {
				sendMsg(chatInfo.Message.Chat.ID, rss)
			}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/sensor") {
			sensor.Reset()
			measure, err := sensor.Measure()
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, err.Error())
			} else {
				sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf(
					"Pressure: %6.3f, Temperature: %4.2f, Humidity: %4.2f, Gas: %6.3f",
					measure.Pressure, measure.Temperature, measure.Humidity, measure.GasResistance))
			}
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/search_song") {
			platform := "kg"
			client, err := provider.GetClient(platform)
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, err.Error())
				return
			}
			keyword := strings.Trim(strings.Replace(chatInfo.Message.Text, "/search_song", "", -1), " \n")
			result, err := client.SearchSongs(context.Background(), keyword)
			var sb strings.Builder
			for i, s := range result {
				fmt.Fprintf(&sb, "[%02d] %s - %s - /download_song%s\n", i+1, s.Name, s.Artist, s.Id)
			}
			sendMsg(chatInfo.Message.Chat.ID, sb.String())
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/download_song") {
			platform := "kg"
			client, err := provider.GetClient(platform)
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, err.Error())
				return
			}
			songId := strings.Trim(strings.Replace(chatInfo.Message.Text, "/download_song", "", -1), " \n")
			ctx := context.Background()
			song, err := client.GetSong(ctx, songId)
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, err.Error())
				return
			}

			sendMsg(chatInfo.Message.Chat.ID, "Start downloading music...")
			go func() {
				ctx := context.Background()
				mp3FilePath, err := ConcurrentDownload(ctx, client, "./cache", song)
				if err != nil {
					sendMsg(chatInfo.Message.Chat.ID, "Error: "+err.Error())
					return
				}
				songInfo := fmt.Sprintf("%s - %s", song.Artist, song.Name)
				// filePath := filepath.Join("./cache", utils.TrimInvalidFilePathChars(songInfo))
				// filePath := utils.TrimInvalidFilePathChars(songInfo)
				// mp3FilePath := filePath + ".mp3"
				sendFile(mp3FilePath, chatInfo.Message.Chat.ID, songInfo)
				// sendMsg(chatInfo.Message.Chat.ID, mp3FilePath)
				if err != nil {
					sendMsg(chatInfo.Message.Chat.ID, "Error: "+err.Error())
					return
				}
			}()

		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/dice") {
			resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendDice?chat_id=%d&emoji=%s",
				ConfigData.TGToken, chatInfo.Message.Chat.ID, "🎲"))
			if err != nil {
				fmt.Println(err)
			}
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(respData))
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/dart") {
			resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendDice?chat_id=%d&emoji=%s",
				ConfigData.TGToken, chatInfo.Message.Chat.ID, "🎯"))
			if err != nil {
				fmt.Println(err)
			}
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(respData))
		} else if strings.HasPrefix(strings.Trim(chatInfo.Message.Text, " \n"), "/book") {
			bookId, err := strconv.Atoi(strings.Trim(strings.Replace(chatInfo.Message.Text, "/book", "", -1), " \n"))
			if err != nil {
				sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("Invalid bookId: %d", bookId))
				sendMsg(chatInfo.Message.Chat.ID, err.Error())
				return
			}
			c := make(chan float32)
			go crawl(bookId, c)
			go func(c <-chan float32) {
				sendMsg(chatInfo.Message.Chat.ID, "Generating book link...")
				u, err := url.Parse(ConfigData.Server)
				if err != nil {
					fmt.Println("url parse error...")
				}
				bookPath := fmt.Sprintf("cache/%d.txt", bookId)
				u.Path = path.Join(u.Path, bookPath)

				for {
					val := <-c
					if int(val*10000)%500 == 0 {
						sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("Downloading progress: %.2f", val))
					}
					if val == 1.0 {
						// Add client
						client := &http.Client{}
						//prepare the reader instances to encode
						fout, err := os.Open(bookPath)
						values := map[string]io.Reader{
							"document": fout, // lets assume its this file
							"chat_id":  strings.NewReader(strconv.Itoa(chatInfo.Message.Chat.ID)),
						}

						err = upload(client, fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", ConfigData.TGToken), values)

						if err != nil {
							fmt.Println(err)
						}
						break
					}
				}
			}(c)
		} else if strings.Trim(chatInfo.Message.Text, " \n") == "/pigeon" {
			sendMsg(chatInfo.Message.Chat.ID, "咕的十种家常做法：")
			for i := 0; i < 10; i++ {
				idx1 := rand.Intn(len(spice))
				idx2 := rand.Intn(len(cook))
				sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("咕的十种家常做法：%s%s咕", spice[idx1], cook[idx2]))
			}
		} else if strings.Trim(chatInfo.Message.Text, " \n") == "/cat" {
			// sendMsg(chatInfo.Message.Chat.ID, "喵的十种家常做法：")
			// for i:=0; i<10; i++ {
			// 	idx1 := rand.Intn(len(spice))
			// 	idx2 := rand.Intn(len(cook))
			// 	sendMsg(chatInfo.Message.Chat.ID, fmt.Sprintf("喵的十种家常做法：%s%s喵", spice[idx1], cook[idx2]))
			// }
			go func() {
				url, err := randImage("cat")
				if err != nil {
					return
				}
				resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument?chat_id=%d&document=%s",
					ConfigData.TGToken, chatInfo.Message.Chat.ID, url))
				if err != nil {
					data, err := ioutil.ReadAll(resp.Body)
					fmt.Println(data)
					fmt.Println(err)
					return
				}
			}()
		} else if strings.Trim(chatInfo.Message.Text, " \n") == "/bingo" {
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
				sendMsg(chatInfo.Message.Chat.ID, "Please wait for your turn...")
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
				}
			} else {
				text, _, _ := isBingo(predict, target)
				sendMsg(chatInfo.Message.Chat.ID, text)
			}

		} else {
			// go to chatgpt
		}

	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
