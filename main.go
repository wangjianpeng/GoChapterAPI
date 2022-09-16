package main

import (
	"GoChaptersAPI/learnginpkg"
	"fmt"
)

func main() {
	// api.DoConsume(1, 1)
	// api.Do1()
	// learnginpkg.DoPingGin()
	// learnginpkg.DoBuildChapterResponseText()
	fmt.Println(learnginpkg.DoReadFile("./files/ugcshop.json"))
}
