package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"errors"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	focusIndex   int
	inputs       []textinput.Model
	inputsByPage [][]textinput.Model
	currentPage  int
	err          error
}

func initialModel() model {
	// Create text input fields
	inputs := make([]textinput.Model, 13)
	inputsByPage := [][]textinput.Model{
		// inputs[a:b] not inclusive of b
		inputs[0:3],
		inputs[3:8],
		inputs[8:13],
	}
	model := model{
		inputsByPage: inputsByPage,
		inputs:       inputs,
		currentPage:  0, // Start on the first page
	}
	// Title field
	dive_title := textinput.New()
	dive_title.Placeholder = "Name your dive"
	dive_title.Focus()
	dive_title.Width = 40
	inputs[0] = dive_title

	// Name field
	dive_site := textinput.New()
	dive_site.Placeholder = "Where did you dive?"
	dive_site.Width = 40
	inputs[1] = dive_site

	when_was_the_dive := textinput.New()
	when_was_the_dive.Placeholder = "When did you dive?"
	when_was_the_dive.Width = 40
	inputs[2] = when_was_the_dive

	dive_type := textinput.New()
	dive_type.Placeholder = "Enter your dive type e.g. boat"
	dive_type.Width = 40
	inputs[3] = dive_type

	water_body := textinput.New()
	water_body.Placeholder = "What type of water were you diving in?"
	water_body.Width = 40
	inputs[4] = water_body

	// Depth / time fields
	bottom_time := textinput.New()
	bottom_time.Placeholder = "Enter your bottom time"
	bottom_time.Width = 40
	inputs[5] = bottom_time

	max_depth := textinput.New()
	max_depth.Placeholder = "Enter your max depth"
	max_depth.Width = 40
	inputs[6] = max_depth

	// Temperature fields
	surface_temp := textinput.New()
	surface_temp.Placeholder = "Enter the surface temp"
	surface_temp.Width = 40
	inputs[7] = surface_temp

	bottom_temp := textinput.New()
	bottom_temp.Placeholder = "Enter the bottom temp"
	bottom_temp.Width = 40
	inputs[8] = bottom_temp

	// Gear fields
	weight := textinput.New()
	weight.Placeholder = "Enter your weight"
	weight.Width = 40
	inputs[9] = weight

	suit_type := textinput.New()
	suit_type.Placeholder = "Enter your suit type e.g. wetsuit"
	suit_type.Width = 40
	inputs[10] = suit_type

	// Gas fields
	start_gas := textinput.New()
	start_gas.Placeholder = "Enter your start gas"
	start_gas.Width = 40
	inputs[11] = start_gas
	end_gas := textinput.New()
	end_gas.Placeholder = "Enter your end gas"
	end_gas.Width = 40
	inputs[12] = end_gas

	// Notes field
	notes := textinput.New()
	notes.Placeholder = "Enter any notes"
	notes.Width = 30
	inputs = append(inputs, notes)
	return model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	calculateGasRemaining := func(startGasStr, endGasStr string) (int, error) {
		startGas, err := strconv.Atoi(startGasStr)
		if err != nil {
			return 0, err
		}
		endGas, err := strconv.Atoi(endGasStr)
		if err != nil {
			return 0, err
		}
		return startGas - endGas, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRight:
			// Navigate to the next page
			if m.currentPage < len(m.inputsByPage)-1 {
				// Save or validate inputs from the current page before switching
				for _, input := range m.inputsByPage[m.currentPage] {
					if input.Value() == "" {
						m.err = errMsg(errors.New("Please fill out all fields on this page before proceeding."))
						return m, nil
					}
				}
				m.currentPage++
			}
		case tea.KeyLeft:
			// Navigate to the previous page
			if m.currentPage > 0 {
				m.currentPage--
			}
		case tea.KeyEnter:
			// Handle Enter key presses
			if m.focusIndex == len(m.inputs)-1 {
				// If the user presses Enter on the last input, calculate gas remaining and quit
				startGas := m.inputs[9].Value()
				endGas := m.inputs[10].Value()
				gasRemaining, err := calculateGasRemaining(startGas, endGas)
				if err != nil {
					m.err = errMsg(err)
				} else {
					log.Printf("Gas remaining: %d", gasRemaining)
				}
				return m, tea.Quit
			} else {
				// Navigate to the next field or page
				if m.focusIndex < len(m.inputsByPage[m.currentPage])-1 {
					// Move to the next input field in the current page
					m.focusIndex++
				} else {
					// Move to the next page if at the end of current page
					if m.currentPage < len(m.inputsByPage)-1 {
						m.currentPage++
						m.focusIndex = 0 // Reset focus index for new page
					}
				}
				// Update focus states
				for i := range m.inputsByPage[m.currentPage] {
					if i == m.focusIndex {
						m.inputsByPage[m.currentPage][i].Focus()
					} else {
						m.inputsByPage[m.currentPage][i].Blur()
					}
				}
			}
		case tea.KeyTab:
			// Navigate to the next input field
			if m.focusIndex < len(m.inputsByPage[m.currentPage])-1 {
				// Move to the next input field in the current page
				m.focusIndex++
			} else {
				// Move to the next page if at the end of the current page
				if m.currentPage < len(m.inputsByPage)-1 {
					m.currentPage++
					m.focusIndex = 0 // Reset focus to the first field of the next page
				}
			}
			for i := range m.inputsByPage[m.currentPage] {
				if i == m.focusIndex {
					m.inputsByPage[m.currentPage][i].Focus()
				} else {
					m.inputsByPage[m.currentPage][i].Blur()
				}
			}
		case tea.KeyShiftTab:
			// Navigate to the previous input field
			if m.focusIndex > 0 {
				// Move to the previous input field in the current page
				m.focusIndex--
			} else {
				// Move to the previous page if at the start of the current page
				if m.currentPage > 0 {
					m.currentPage--
					m.focusIndex = len(m.inputsByPage[m.currentPage]) - 1 // Set focus to the last field of the previous page
				}
			}
			for i := range m.inputsByPage[m.currentPage] {
				if i == m.focusIndex {
					m.inputsByPage[m.currentPage][i].Focus()
				} else {
					m.inputsByPage[m.currentPage][i].Blur()
				}
			}

		case tea.KeyCtrlC, tea.KeyCtrlD, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Update the focused text input field
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m model) View() string {
	// Render the form
	form := "Enter your dive details\n\n"
	for _, input := range m.inputsByPage[m.currentPage] {
		form += input.View() + "\n"
	}

	// Add navigation indicator
	form += "\nPage " + strconv.Itoa(m.currentPage+1) + " of " + strconv.Itoa(len(m.inputsByPage)) + "\n"

	// Display gas remaining if calculated
	if m.err == nil && m.inputs[9].Value() != "" && m.inputs[10].Value() != "" {
		startGas := m.inputs[9].Value()
		endGas := m.inputs[10].Value()
		gasRemaining, err := strconv.Atoi(startGas)
		if err == nil {
			if endGasInt, err := strconv.Atoi(endGas); err == nil {
				gasRemaining -= endGasInt
			}
			form += "\nGas remaining: " + strconv.Itoa(gasRemaining) + "\n"
		}
	}

	form += "\n(Press Enter to submit, Tab to switch fields, Esc to quit)"
	return form
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
