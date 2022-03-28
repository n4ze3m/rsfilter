package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/n4ze3m/rsfilter/utils"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter keywords to search for: ")
	scanner.Scan()
	text := strings.Trim(strings.ToLower(scanner.Text()), " ")
	keywords := strings.Split(utils.RemoveSpecialCharacters(text), " ")
	pdfs := utils.GetAllPdfNames()

	bestCandiates := make([]string, 0)

	for _, pdf := range pdfs {
		t, err := utils.ReadPDF(pdf)
		if err != nil {
			fmt.Println(err)
		}

		for _, keyword := range keywords {
			if strings.Contains(t, keyword) {
				if !utils.Contains(bestCandiates, pdf) {
					bestCandiates = append(bestCandiates, pdf)
				}
			}
		}
	}

	dir := utils.MoveFiles(&bestCandiates, &keywords)

	fmt.Printf("Suitable candidates resumés are in %s\n", dir)
}
