package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
)

// Call of Cthulhu modules
func init() {
}

type Player struct {
	ChatId int `gorm:"primaryKey"`
	Name   string

	// Primary
	Strength     int
	Constitution int
	Power        int
	Dexterity    int
	Appearance   int
	Size         int
	Intelligence int
	Education    int

	// Secondary
	Luck        int
	MP          int
	DamageBonus int
	Build       int

	SP  int // skill points
	HP  int // health points
	SAN int // sanity
	Age int
}

// cat n m-facet dice and return result
func nDm(n int, m int) int {
	ret := 0
	for i := 0; i < n; i++ {
		ret += rand.Intn(m) + 1
	}
	return ret
}

func createPlayer(chatId int, name string) {
	db, err := gorm.Open(sqlite.Open("game.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Open database error: %v\n", err)
	}

	db.AutoMigrate(&Player{})

	// TODOs: check of the user exist in the database

	player := Player{}
	player.ChatId = chatId
	player.Name = name
	player.Strength = nDm(3, 6) * 5
	player.Constitution = nDm(3, 6) * 5
	player.Size = (nDm(2, 6) + 6) * 5
	player.Dexterity = nDm(3, 6) * 5
	player.Appearance = nDm(3, 6) * 5
	player.Intelligence = (nDm(2, 6) + 6) * 5
	player.Power = nDm(3, 6) * 5
	player.Education = nDm(3, 6) * 5
	player.Luck = nDm(3, 6) * 5
	player.MP = player.Power / 5

	s := player.Strength + player.Size
	if s >= 2 && s <= 64 {
		player.DamageBonus = -2
		player.Build = -2
	} else if s >= 65 && s <= 84 {
		player.DamageBonus = -1
		player.Build = -1
	} else if s >= 85 && s <= 124 {
		player.DamageBonus = 0
		player.Build = 0
	} else if s >= 125 && s <= 164 {
		player.DamageBonus = nDm(1, 4)
		player.Build = 1
	} else if s >= 165 && s <= 204 {
		player.DamageBonus = nDm(1, 6)
		player.Build = 2
	} else {
		player.DamageBonus = 0
		player.Build = 0
	}

	// result := db.Create(&player)
	// result := db.Model(&player).Updates(player)
	// result := db.Save(&player)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&player)
	fmt.Printf("%v \n", result.Error)
}
