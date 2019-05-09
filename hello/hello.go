package main

import (
	"fmt"
)

type VowelsFinder interface {
	FindVowels() []rune
	ChangeString() string
}

type MyString string

func (ms MyString) ChangeString() MyString {
	// ms = "yyy"
	return MyString("yyy")
}

func (ms *MyString) FindVowels() []rune {
	var vowel []rune
	for _, v := range *ms {
		if v == 'u' || v == 'e' || v == 'o' || v == 'a' || v == 'i' {
			vowel = append(vowel, v)
		}
	}
	*ms = "xxx"
	return vowel
}

func hello() {
	fmt.Println("Hello world goroutine")
}
func main() {
	var ms MyString = "Hello World"
	//ms := MyString("Hello World")
	// var vowelsFinder VowelsFinder
	// vowelsFinder = &ms
	var r []rune = ms.FindVowels()
	fmt.Printf("Vowels are %c ", r)
	fmt.Println()
	var m MyString = ms.ChangeString()
	// fmt.Println(ms.ChangeString())
	fmt.Println(m)
}
