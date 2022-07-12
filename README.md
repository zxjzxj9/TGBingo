# TGBingo
A bingo played on Telegram Bot

## Introduction
This is Go lang bot server used to communicate with Telegram's Bot. The scheme can be shown as follows, telegram user can send requests to the bot server through telegram app, 
then bot server reply user according to given instructions. It is basically a https sever, so I build the whole server based on gin-gonic framework.

User <--> Telegram <--> Bot Server

## Build
Just use the command `go vendor` and `go build .` to build the running binary. If you want to run the binary on different arch/os, i.e. raspberry pi, just specify the GOARCH and GOOS env to the target arch/os.  

## Serve
