package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	dates          map[string][]string
	activeDate     int
	activeMonth    string
	inputfield     textarea.Model
	showinputfield bool
	monthOrder     []string
}

type Change int

const (
	Decrease Change = 0
	Increase Change = 1
)

func getfileinfo(month string, day int) string {
	// user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// filename
	filename := strconv.Itoa(day) + ".txt"

	// full file path
	return filepath.Join(home, "journal", month, filename)
}

func checkMonthChange(allDaysInMonth []string, day int, change Change) bool {
	idx := slices.Index(allDaysInMonth, strconv.Itoa(day))

	changeMonth := false

	// increase month
	if change == Increase {
		if idx == len(allDaysInMonth)-1 || idx == -1 {
			changeMonth = true
		}
	}

	//decrease month
	if change == Decrease {
		if idx == 0 || idx == -1 {
			changeMonth = true
		}
	}

	log.Println("changeMonth : ", changeMonth, " idx: ", idx)

	return changeMonth
}

func InitialModel() model {

	m := make(map[string][]string)

	m["January"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["February"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28"}

	m["March"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["April"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30"}

	m["May"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["June"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30"}

	m["July"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["August"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["September"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30"}

	m["October"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	m["November"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30"}

	m["December"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

	// get todays date
	day := time.Now().Day()
	month := time.Now().Month()

	inputfield := textarea.New()
	inputfield.Focus()
	inputfield.Placeholder = "how's your day going"
	inputfield.ShowLineNumbers = false

	monthOrder := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	return model{
		dates:          m,
		activeDate:     day,
		activeMonth:    month.String(),
		inputfield:     inputfield,
		showinputfield: false,
		monthOrder:     monthOrder,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// when text editor is open all keys should be directed there
	if m.showinputfield {
		// Handle resize and special keys before passing to textarea
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.showinputfield = false

				// get input content
				content := m.inputfield.Value()

				// write to system
				filePath := getfileinfo(m.activeMonth, m.activeDate)

				os.MkdirAll(filepath.Dir(filePath), 0755)

				os.WriteFile(filePath, []byte(content), 0644)

				// reset text area
				m.inputfield.SetValue("")

				return m, nil
			}
		case tea.WindowSizeMsg:
			// Set dimensions directly and also pass to textarea so it reflows correctly
			m.inputfield.SetWidth(msg.Width)
			m.inputfield.SetHeight(msg.Height - 5)
		}

		// Pass all other messages (including keypresses) to the textarea
		var cmd tea.Cmd
		m.inputfield, cmd = m.inputfield.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()
		switch key {

		case "left":
			shouldChangeMonth := checkMonthChange(m.dates[m.activeMonth], m.activeDate, Decrease)

			if shouldChangeMonth && m.activeMonth != "January" {
				idx := slices.Index(m.monthOrder, m.activeMonth)
				m.activeMonth = m.monthOrder[idx-1]
				// date should be last day of month
				daysInMonth := m.dates[m.activeMonth]
				lastDay := daysInMonth[len(daysInMonth)-1]
				lastDayInt, _ := strconv.ParseInt(lastDay, 10, 64)
				m.activeDate = int(lastDayInt)

			} else {
				if m.activeDate == 1 && m.activeMonth == "January" {
					break
				} else {
					m.activeDate--
				}
			}
		case "right":

			shouldChangeMonth := checkMonthChange(m.dates[m.activeMonth], m.activeDate, Increase)

			if shouldChangeMonth && m.activeMonth != "December" {
				idx := slices.Index(m.monthOrder, m.activeMonth)
				m.activeMonth = m.monthOrder[idx+1]

				// date should be first day of month
				daysInMonth := m.dates[m.activeMonth]
				firstDay := daysInMonth[0]
				firstDayInt, _ := strconv.ParseInt(firstDay, 10, 64)
				m.activeDate = int(firstDayInt)

			} else {
				if m.activeDate == 31 && m.activeMonth == "December" {
					break
				} else {
					m.activeDate++
				}
			}

		case "up":

			nextValue := m.activeDate - 7
			shouldChangeMonth := false
			idx := slices.Index(m.dates[m.activeMonth], strconv.Itoa(nextValue))

			if idx == 0 || idx == -1 {
				shouldChangeMonth = true
			}

			log.Println("changeMonth : ", shouldChangeMonth, " idx: ", idx)

			if shouldChangeMonth && m.activeMonth != "January" {
				idx := slices.Index(m.monthOrder, m.activeMonth)
				m.activeMonth = m.monthOrder[idx-1]
				// date should be last day of month
				daysInMonth := m.dates[m.activeMonth]
				lastDay := daysInMonth[len(daysInMonth)-1]
				lastDayInt, _ := strconv.ParseInt(lastDay, 10, 64)
				m.activeDate = int(lastDayInt)

			} else if idx == -1 && m.activeMonth == "January" {
				break
			} else {
				m.activeDate = m.activeDate - 7
			}

		case "down":
			nextValue := m.activeDate + 7
			shouldChangeMonth := false
			idx := slices.Index(m.dates[m.activeMonth], strconv.Itoa(nextValue))

			if idx == 0 || idx == -1 {
				shouldChangeMonth = true
			}

			log.Println("changeMonth : ", shouldChangeMonth, " idx: ", idx)

			if shouldChangeMonth && m.activeMonth != "December" {
				idx := slices.Index(m.monthOrder, m.activeMonth)
				m.activeMonth = m.monthOrder[idx+1]

				// date should be first day of month
				daysInMonth := m.dates[m.activeMonth]
				firstDay := daysInMonth[0]
				firstDayInt, _ := strconv.ParseInt(firstDay, 10, 64)
				m.activeDate = int(firstDayInt)

			} else if idx == -1 && m.activeMonth == "December" {
				break
			} else {
				m.activeDate = m.activeDate + 7
			}

		case "a":

			// append previous entry
			filePath := getfileinfo(m.activeMonth, m.activeDate)
			content, err := os.ReadFile(filePath)

			if err != nil {
				// file does not exit
				log.Println("filepath not exist : " + filePath)
			} else {
				// file exists
				m.inputfield.SetValue(string(content))
			}

			m.showinputfield = true

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.inputfield.SetHeight(msg.Height - 5)
		m.inputfield.SetWidth(msg.Width - 10)
		log.Println("Height: " + strconv.Itoa(msg.Height) + "\tWidth: " + strconv.Itoa(msg.Width))

	}

	return m, nil

}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF75B5"))

	// Active/selected date
	activeDateStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FF75B5")).
			Padding(0, 1)

	// Normal date
	normalDateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)

	// Journal content
	contentStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#7B61FF")).
			Padding(1, 2).
			MarginTop(1)

	// No entry message
	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")).
			Italic(true)
)

func (m model) View() string {

	// get all days with entries in this month
	homedir, err  := os.UserHomeDir()
	if err != nil {
		log.Println("error finding home dir", err.Error() )
	}

	monthpath := homedir + "/journal/" + m.activeMonth

	files, err := os.ReadDir(monthpath)
	if err != nil {
		log.Println("error", err.Error())
	}

	log.Println("monthpath :", monthpath)

	var allentries []string

	for _, file := range files {
		name := strings.Replace(file.Name(), ".txt", "", 1)
		allentries = append(allentries, name)
	}

	log.Println("allentries: ", allentries)

	// read the file
	filePath := getfileinfo(m.activeMonth, m.activeDate)
	content, err := os.ReadFile(filePath)

	if m.showinputfield {
		return fmt.Sprintf("%s\n\n ", m.inputfield.View())
	}

	var calendarContent strings.Builder

	// Styled header
	calendarContent.WriteString(headerStyle.Render("Day: "+strconv.Itoa(m.activeDate)) + "\n")
	calendarContent.WriteString(headerStyle.Render("Month: "+m.activeMonth) + "\n\n")

	counter := 0
	for _, v := range m.dates[m.activeMonth] {
		s_v, _ := strconv.Atoi(v)

		var formatted string
		if slices.Contains(allentries, v) {
			formatted = fmt.Sprintf("%3s", "#"+v)
		} else {
			formatted = fmt.Sprintf("%3s", v)
		}

		if s_v == m.activeDate {
			formatted = activeDateStyle.Render(formatted)
		} else {
			formatted =  normalDateStyle.Render(formatted)
		}

		calendarContent.WriteString(formatted)

		counter++
		if counter == 7 {
			counter = 0
			calendarContent.WriteString("\n")
		}
	}

	var journalContent string
	if err != nil {
		journalContent = emptyStyle.Render("no entry")
	} else {
		var preview string
		if len(string(content)) > 1000 {
			preview = string(content)[:1000] + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF75B5")).
				Render("\n\npress a to view more....")
		} else {
			preview = string(content)
		}
		journalContent = "\n" + contentStyle.Render(preview)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, calendarContent.String(), journalContent)
}

func main() {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	// Set the log output to the file
	log.SetOutput(file)

	log.Println("Application started at ", time.Now())

	// new bubbletea program
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
