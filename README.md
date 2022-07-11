# TGBingo
A bingo played on Telegram Bot

## Introduction
This is Go lang bot server used to communicate with Telegram's Bot. The scheme can be shown as follows, telegram user can send requests to the bot server through telegram app, 
then bot server reply user according to given instructions. It is basically a https sever, so I build the whole server based on gin-gonic framework.

User <--> Telegram <--> Bot Server
