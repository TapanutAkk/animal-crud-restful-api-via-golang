package utils

import (
	"math/rand"
	"time"
)

var AnimalNames = []string{"Leo", "Bao", "Koko", "Zara", "Milo", "Luna", "Simba", "Nala", "Rocky", "Bella"}
var AnimalSpecies = []string{"Lion", "Tiger", "Bear", "Panda", "Elephant", "Monkey", "Dolphin", "Eagle", "Wolf", "Fox"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomName() string {
	index := rand.Intn(len(AnimalNames))
	return AnimalNames[index]
}

func RandomSpecies() string {
	index := rand.Intn(len(AnimalSpecies))
	return AnimalSpecies[index]
}

func RandomAge() int {
	return rand.Intn(15) + 1
}

func RandomAnimal() (string, string, int) {
	name := RandomName()
	species := RandomSpecies()
	age := RandomAge()
	return name, species, age
}