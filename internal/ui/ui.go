package ui

import (
	"fmt"

	"github.com/aletomasella/namepicker-cli/internal/lang"
	"github.com/aletomasella/namepicker-cli/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	FILE   = "file"
	RANDOM = "random"
	MANUAL = "manual"
)

type Model struct {
	choices          []string
	cursor           int
	selected         map[int]struct{}
	selectedLanguage string
	selectedSource   string
}

type Header lang.Label

type Footer lang.Label

// First the user must select the language from the available languages
// Then the user must select from where do they want to get the names from:
// - Random names
// - Names from a file
// - Write the names manually separated by commas
// Then we execute our logic randomly selecting a name from the selected source and displaying it until the user quits or the source is empty

func inicializeNameChoices() map[string][]string {
	nameChoices := make(map[string][]string)

	nameSeed := []string{"John", "Jane", "Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy", "Kevin", "Laura", "Michael", "Nancy", "Oliver", "Peggy", "Quincy", "Rita", "Steve", "Tina", "Ursula", "Victor", "Wendy", "Xavier", "Yvonne", "Zack"}

	// Randomize the names
	randomSeed := utils.RandomizeSlice(nameSeed)

	nameChoices[RANDOM] = randomSeed
	nameChoices[FILE] = make([]string, 0)
	nameChoices[MANUAL] = make([]string, 0)

	return nameChoices
}

func InitialModel() Model {

	nameChoices := inicializeNameChoices()

	// Default name choices are the random names

	return Model{
		choices:  nameChoices[RANDOM],
		selected: make(map[int]struct{}),
		cursor:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// The r key randomizes the order of the choices in the list
		case "r":
			m.choices = utils.RandomizeSlice(m.choices)

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header

	// Default language is English
	// The header

	header := Header{
		Text: map[string]string{
			lang.English: "Name Picker\n\n",
			lang.Spanish: "Selector de Nombres\n\n",
		},
		DefineLanguages: lang.GetAvailableLanguages(),
	}

	view := header.Text[lang.Spanish]

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	footer := Footer{
		Text: map[string]string{
			lang.English: "Press q to quit.\n",
			lang.Spanish: "Presiona q para salir.\n",
		},
		DefineLanguages: lang.GetAvailableLanguages(),
	}

	view += footer.Text[lang.Spanish]

	// Send the UI for rendering
	return view
}
