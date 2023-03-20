package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	commandViewStyle = lipgloss.NewStyle().
				PaddingRight(1).
				MarginRight(1).
				Border(lipgloss.RoundedBorder(), false, true, false, false)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
)

type Model struct {
	width, height     int
	operationViewport viewport.Model
}

func InitialModel() (*Model, error) {
	return &Model{}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.operationViewport, cmd = m.operationViewport.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		vp, _ := operationViewPort(m.width/2-2, m.height-2)
		m.operationViewport = vp

	}
	return m, nil
}
func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.commandView(), m.operationView()),
		m.statusView(),
	)
}
func (m Model) commandView() string {
	return commandViewStyle.Width(m.width / 2).Height(m.height).Render("list view")
}

func (m Model) operationView() string {
	return m.operationViewport.View() + m.helpView()
}

func (m Model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate\n")
}
func (m Model) statusView() string {
	return statusBarStyle.Width(m.width).Render("status")
}
