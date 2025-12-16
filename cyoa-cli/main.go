package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

type Story = map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type model struct {
	chapter Chapter
	cursor  int
}

var story Story

func initialModel(story *Story) model {
	chapter, ok := (*story)["intro"]
	if !ok {
		fmt.Println("Alas, there's no intro chapter in the page")
		os.Exit(1)
	}

	return model{
		chapter: chapter,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.chapter.Options)-1 {
				m.cursor++
			}

		case "enter", " ":
			selectedOption := m.chapter.Options[m.cursor]
			m.chapter = story[selectedOption.Arc]
		}

	}
	return m, nil
}

func (m model) View() string {
	width := 100

	s := "\n                               GOPHER STORY                            \n"
	s += "=======================================================================\n\n"

	// paragraphs
	for _, p := range m.chapter.Story {
		s += fmt.Sprintf("%v\n\n", wordwrap.String(p, width))
	}

	// options
	s += "\n\nSelect options\n\n"

	for i, option := range m.chapter.Options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, wordwrap.String(option.Text, width))
	}

	s += "\nPress q to quit."
	return s
}

func main() {
	filename := flag.String("json", "gopher.json", "The path to the JSON file containing the CYOA story")
	jsonBytes, err := os.ReadFile(*filename)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonBytes, &story); err != nil {
		panic(err)
	}

	p := tea.NewProgram(initialModel(&story))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
