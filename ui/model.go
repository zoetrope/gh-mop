package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	commandViewStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				MarginTop(2).
				PaddingRight(2)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
)

type Model struct {
	width, height     int
	operationViewport viewport.Model
	command           textinput.Model
}

func InitialModel() (*Model, error) {
	c := newCommandModel()
	return &Model{
		command: c,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+u", "ctrl+d":
			m.operationViewport, cmd = m.operationViewport.Update(msg)
			return m, cmd
			//default:
			//	m.operationViewport, cmd = m.operationViewport.Update(msg)
			//	return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		vp, _ := operationViewPort(m.width/2-2, m.height-2)
		m.operationViewport = vp

	}

	m.command, cmd = m.command.Update(msg)
	return m, cmd
}
func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.commandView(), m.operationView()),
		m.statusView(),
	)
}
func (m Model) commandView() string {
	return commandViewStyle.Width(m.width/2-2).Height(m.height-4).Render("Please enter a command", m.command.View())
	//return "Please enter a command\n" + m.command.View()
}

func (m Model) operationView() string {
	return m.operationViewport.View() + m.helpView()
}

func (m Model) helpView() string {
	return helpStyle("\n  Up: ctrl+u, Down: ctrl+d\n")
}
func (m Model) statusView() string {
	return statusBarStyle.Width(m.width).Render("status")
}
