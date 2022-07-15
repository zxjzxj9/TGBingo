# TGBingo
A bingo played on Telegram Bot

## Introduction
This is Go lang bot server used to communicate with Telegram's Bot. The scheme can be shown as follows, telegram user can send requests to the bot server through telegram app, 
then bot server reply user according to given instructions. It is basically a https sever, so I build the whole server based on gin-gonic framework.

User <--> Telegram <--> Bot Server

## Build
Just use the command `go vendor` and `go build .` to build the running binary. If you want to run the binary on different arch/os, i.e. raspberry pi, just specify the GOARCH and GOOS env to the target arch/os.  

## Serve
Edit config.json file, add something like the followings, then start the server using `./server -c config.json`

```
{
  server: "your_server_https_url",
  weather_token: "token_from_openweathermap_org",
  tg_toke: "telegram_bot_token"
}
```

## API Introduction
The bot currently support several APIs, including the following.
Help (type /help and it displays as floows):

```
####################
 NekoRoid is a bot having a lot of fun

 #    # ###### #    #  ####  
 ##   # #      #   #  #    # 
 # #  # #####  ####   #    # 
 #  # # #      #  #   #    # 
 #   ## #      #   #  #    # 
 #    # ###### #    #  ####

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
####################
```
