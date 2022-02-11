package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
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

	MOV int
	SP  int // skill points
	HP  int // health points
	SAN int // sanity
	Age int
}

type Occupation struct {
	Id int
}

// cat n m-facet dice and return result
func nDm(n int, m int) int {
	ret := 0
	for i := 0; i < n; i++ {
		ret += rand.Intn(m) + 1
	}
	return ret
}

func enhance(val int) int {
	ret := val
	if nDm(1, 100) > val {
		ret += nDm(1, 10)
	}
	if ret >= 99 {
		return 99
	}
	return ret
}

func reduce(rval int, vals ...*int) {
	sz := len(vals)
	frac := make([]float64, sz)
	sum := 0.0
	for i := 0; i < sz; i++ {
		frac[i] = rand.Float64()
		sum += frac[i]
	}
	for i := 0; i < sz-1; i++ {
		frac[i] = math.Round(frac[i] * float64(rval) / sum)
		rval -= int(frac[i])
		*vals[i] = int(frac[i])
	}
	*vals[sz-1] = rval
}

func createPlayer(chatId int, name string) *Player {
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

	player.Age = rand.Intn(75) + 15

	// * 5, yards per second
	if player.Dexterity < player.Size && player.Strength < player.Size {
		player.MOV = 7
	} else if player.Dexterity >= player.Size && player.Strength >= player.Size {
		player.MOV = 9
	} else if player.Dexterity >= player.Size || player.Strength >= player.Size {
		player.MOV = 8
	}

	switch {
	case player.Age >= 15 && player.Age <= 19:
		player.Strength -= 5
		player.Size -= 5
		player.Education -= 5
		t := nDm(3, 6) * 5
		if t >= player.Luck {
			player.Luck = t
		}
	case player.Age >= 20 && player.Age <= 39:
		player.Education = enhance(player.Education)
	case player.Age >= 40 && player.Age <= 49:
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		reduce(5, &player.Strength, &player.Constitution, &player.Dexterity)
		reduce(5, &player.Appearance)
		player.MOV -= 1
	case player.Age >= 50 && player.Age <= 59:
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		reduce(10, &player.Strength, &player.Constitution, &player.Dexterity)
		reduce(10, &player.Appearance)
		player.MOV -= 2
	case player.Age >= 60 && player.Age <= 69:
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		reduce(20, &player.Strength, &player.Constitution, &player.Dexterity)
		reduce(15, &player.Appearance)
		player.MOV -= 3
	case player.Age >= 70 && player.Age <= 79:
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		reduce(40, &player.Strength, &player.Constitution, &player.Dexterity)
		reduce(20, &player.Appearance)
		player.MOV -= 4
	case player.Age >= 80 && player.Age <= 89:
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		player.Education = enhance(player.Education)
		reduce(80, &player.Strength, &player.Constitution, &player.Dexterity)
		reduce(25, &player.Appearance)
		player.MOV -= 5
	}

	s := player.Strength + player.Size
	switch {
	case s >= 2 && s <= 64:
		player.DamageBonus = -2
		player.Build = -2
	case s >= 65 && s <= 84:
		player.DamageBonus = -1
		player.Build = -1
	case s >= 85 && s <= 124:
		player.DamageBonus = 0
		player.Build = 0
	case s >= 125 && s <= 164:
		player.DamageBonus = nDm(1, 4)
		player.Build = 1
	case s >= 165 && s <= 204:
		player.DamageBonus = nDm(1, 6)
		player.Build = 2
	case s >= 205 && s <= 284:
		player.DamageBonus = nDm(2, 6)
		player.Build = 3
	case s >= 285 && s <= 364:
		player.DamageBonus = nDm(3, 6)
		player.Build = 4
	case s >= 365 && s <= 444:
		player.DamageBonus = nDm(4, 6)
		player.Build = 5
	case s >= 445 && s <= 524:
		player.DamageBonus = nDm(5, 6)
		player.Build = 6
	default:
		player.DamageBonus = 0
		player.Build = 0
	}

	player.HP = (player.Constitution + player.Size) / 10

	// result := db.Create(&player)
	// result := db.Model(&player).Updates(player)
	// result := db.Save(&player)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&player)
	fmt.Printf("%v \n", result.Error)
	if result.Error != nil {
		return nil
	}
	return &player
}

func createOccupation(id int) {
	switch id {
	case 0:
		// ANTIQUARIAN

	}
}
