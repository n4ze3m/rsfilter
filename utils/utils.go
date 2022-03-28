package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func RemoveSpecialCharacters(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, str)
}

// get all .pdf files in current directory
func GetAllPdfNames() []string {
	pdfs := make([]string, 0)
	files, _ := GetAllFiles()
	for _, file := range files {
		if strings.HasSuffix(file, ".pdf") {
			pdfs = append(pdfs, file)
		}
	}
	return pdfs
}

// get all files from current working dir
func GetAllFiles() ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ReadPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	words := make([]string, 0)

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				lastTextStyle = text
			}

			words = append(words, lastTextStyle.S)
		}

	}
	return strings.ToLower(strings.Join(words, "")), nil
}

func isSameSentence(text pdf.Text, lastTextStyle pdf.Text) bool {
	return text.Font == lastTextStyle.Font &&
		text.FontSize == lastTextStyle.FontSize &&
		text.X == lastTextStyle.X &&
		text.Y == lastTextStyle.Y
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func MoveFiles(files *[]string, keywords *[]string)  string{
	dirName := strings.Join(*keywords, "-")
	dir := fmt.Sprintf("%s-%v", dirName, time.Now().UnixNano())
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
	for _, file := range *files {
		os.Rename(file, filepath.Join(dir, file))
	}

	return dir
}
