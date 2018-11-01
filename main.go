package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var logger = log.New(os.Stderr, "", 0)

type Title string
type Story string
type AnswerBook map[string]string
type MadLib struct {
	Title
	Story
	AnswerBook
}

func (s Story) toTemplate() (*template.Template, error) {
	templ := template.New("madlib")
	templ, err := templ.Parse(string(s))
	if err != nil {
		logger.Println(`error: unable to parse template from string`)
		return templ, err
	}

	return templ, err
}

func (ab AnswerBook) Fill() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n----------")
	for k := range ab {
		s := getMessageForUser(k)
		fmt.Print(s)
		scanner.Scan()
		ab[k] = scanner.Text()
	}
}

func getMessageForUser(s string) string {
	ss := strings.Split(s, " ")
	ss = ss[:len(ss)-1] // Reslice without last element (i.e., number)
	s = strings.Join(ss, " ")
	s = s + ": "

	return s
}

func NewTitle(title string) Title {
	return Title(strings.TrimSpace(title))
}

func NewStory(story string, descriptions []string) Story {
	story = strings.TrimSpace(story)

	for _, description := range descriptions {
		templateString := `{{index . "` + description + `"}}`
		story = strings.Replace(story, "_____", templateString, 1)
	}

	return Story(story)
}

func NewAnswerBook(descriptions []string) AnswerBook {
	ab := make(AnswerBook)

	shuffleDescriptions(descriptions)
	for _, description := range descriptions {
		ab[description] = ""
	}

	return ab
}

func shuffleDescriptions(descriptions []string) {
	rand.Shuffle(len(descriptions), func(i, j int) {
		descriptions[i], descriptions[j] = descriptions[j], descriptions[i]
	})
}

func sanitiseDescriptions(descriptions string) []string {
	return strings.Split(strings.TrimSpace(descriptions), "\n")
}

func readTemplateFromFile(filename string) (string, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Println(`error:`, filename, `could not be read`)
		return "", err
	}

	return string(bs), nil
}

func splitTemplateToParts(template string) (title, story, descriptions string) {
	ss := strings.Split(template, "-----")
	return ss[0], ss[1], ss[2]
}

func (ml *MadLib) Parse(filename string) error {
	st, err := readTemplateFromFile(filename)
	if err != nil {
		return err
	}

	title, story, desc := splitTemplateToParts(st)
	ml.Title = NewTitle(title)
	descriptions := sanitiseDescriptions(desc)
	ml.Story = NewStory(story, descriptions)
	ml.AnswerBook = NewAnswerBook(descriptions)

	return nil
}

func (ml *MadLib) Print() error {
	t, err := ml.Story.toTemplate()
	if err != nil {
		return err
	}

	fmt.Print("\n")
	fmt.Println(ml.Title)
	fmt.Println("----------")
	t.Execute(os.Stdout, ml.AnswerBook)
	fmt.Print("\n")

	return nil
}

func chooseRandomStory(dir string) string {
	files, err := filepath.Glob(filepath.Join(dir, "*.mdlb"))
	if err != nil {
		logger.Println("No stories found")
		os.Exit(0)
	}

	numF := len(files)
	if numF == 0 {
		logger.Printf("No story files found in '%v'\n", dir)
		os.Exit(0)
	}

	return files[rand.Intn(numF)]
}

func isValidTemplate(filename string) (bool, int, int, error) {
	st, err := readTemplateFromFile(filename)
	if err != nil {
		return false, 0, 0, err
	}

	_, story, descriptions := splitTemplateToParts(st)
	numBlanks := strings.Count(story, "_____")
	numDescriptions := len(sanitiseDescriptions(descriptions))
	isValid := numBlanks == numDescriptions

	return isValid, numBlanks, numDescriptions, nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fs := flag.NewFlagSet("gomadlibs", flag.ExitOnError)

	var verifyIntegrity bool
	fs.BoolVar(&verifyIntegrity, "verify-integrity", false, "verifies the integrity of the story template")

	var storiesDir string
	fs.StringVar(&storiesDir, "stories-dir", "./stories", "where the app will look for story templates")

	var help bool
	fs.BoolVar(&help, "help", false, "show this help message")

	args := os.Args[1:]
	fs.Parse(args)

	if help {
		logger.Println("Usage of gomadlibs:")
		fs.PrintDefaults()
		logger.Println("\nExamples:")
		logger.Println("\tgomadlibs story1.mdlb")
		logger.Println("\tgomadlibs /home/user/stories/astory.mdlb")
		logger.Println("\tgomadlibs -verify-integrity story2.mdlb")
		logger.Println("\tgomadlibs -stories-dir mystories astory.mdlb")
		os.Exit(0)
	}

	rargs := args[len(args)-fs.NArg():]
	if verifyIntegrity && len(rargs) == 0 {
		logger.Println("Path to story template must be provided with the '-verify-integrity' flag")
		os.Exit(0)
	}

	var pt string
	if len(rargs) == 0 {
		pt = chooseRandomStory(storiesDir)
	} else {
		if filepath.IsAbs(rargs[0]) {
			pt = rargs[0]
		} else {
			pt = filepath.Join(storiesDir, rargs[0])
		}
	}

	isValid, numBlanks, numDescriptions, err := isValidTemplate(pt)
	if err != nil {
		os.Exit(0)
	}

	logger.Printf("Template at '%v' has %d blanks and %d descriptions", pt, numBlanks, numDescriptions)
	if !isValid {
		logger.Println("Template is invalid. Exiting...")
		os.Exit(0)
	}

	if verifyIntegrity {
		logger.Println("Template is valid. Exiting...")
		os.Exit(0)
	}

	ml := MadLib{}

	err = ml.Parse(pt)
	if err != nil {
		os.Exit(0)
	}

	ml.AnswerBook.Fill()

	err = ml.Print()
	if err != nil {
		os.Exit(0)
	}
}
