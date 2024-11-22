package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v2"
)

type GlossarFile struct {
	Tags    []string        `yaml:"tags"`
	Glossar []GlossaryEntry `yaml:"glossar"`
}

type GlossaryEntry struct {
	Id         string   `yaml:"id"`
	Term       string   `yaml:"term"`
	Definition string   `yaml:"definition"`
	Tags       []string `yaml:"tags"`
}

var glossaryFilePath = "glossary.yaml"
var glossarFile GlossarFile
var glossary []GlossaryEntry
var tags []string
var quickSearchVisible bool
var searchEntry *customEntry

// var resultsContainer *fyne.Container

type customEntry struct {
	widget.Entry
	onFocusGained func()
	onKeyDown     func(*fyne.KeyEvent)
	onFocusLost   func()
}

func newCustomEntry(onFocusGained func(), onKeyDown func(*fyne.KeyEvent), onFocusLost func()) *customEntry {
	entry := &customEntry{onFocusGained: onFocusGained, onKeyDown: onKeyDown, onFocusLost: onFocusLost}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *customEntry) FocusGained() {
	e.Entry.FocusGained()
	if e.onFocusGained != nil {
		e.onFocusGained()
	}
}

func (e *customEntry) FocusLost() {
	e.Entry.FocusLost()
	if e.onFocusLost != nil {
		e.onFocusLost()
	}
}

func (e *customEntry) TypedKey(event *fyne.KeyEvent) {
	e.Entry.TypedKey(event)
	if e.onKeyDown != nil {
		e.onKeyDown(event)
	}
}

func loadGlossary() {
	file, err := os.ReadFile(glossaryFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			glossarFile = GlossarFile{}
			return
		}
		log.Fatalf("Failed to read glossary file: %v", err)
	}
	err = yaml.Unmarshal(file, &glossarFile)
	if err != nil {
		log.Fatalf("Failed to unmarshal glossary file: %v", err)
	}
	glossary = glossarFile.Glossar
	tags = glossarFile.Tags
}

func saveGlossary() {
	glossarFile.Glossar = glossary
	glossarFile.Tags = tags
	RefreshTagList()
	data, err := yaml.Marshal(&glossarFile)
	if err != nil {
		log.Fatalf("Failed to marshal glossary: %v", err)
	}
	err = os.WriteFile(glossaryFilePath, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write glossary file: %v", err)
	}
}

func main() {
	os.Setenv("FYNE_THEME", "dark")
	loadGlossary()

	a := app.NewWithID("com.example.glossary")
	w := a.NewWindow("Glossary Manager")

	glossarTabItem := SetupGlossar(w)
	tagTabItem := SetupTags(w)

	tabs := container.NewAppTabs(
		glossarTabItem,
		tagTabItem,
	)

	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)

	quickSearch := SetupQuicksearch(a)

	// System tray setup
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Glossary",
			fyne.NewMenuItem("Open Glossary Manager", func() {
				w.Show()
				FocusGlossarSearch(w)
			}),
			fyne.NewMenuItem("Open Quick Search", func() {
				quickSearch.Show()
				quickSearchVisible = true
				FocusSearchEntry(quickSearch)
			}),
			// fyne.NewMenuItem("Quit", func() {
			// 	a.Quit()
			// }),
		)

		bytes, err := os.ReadFile("icon_64.png")
		if err != nil {
			panic(err)
		}

		// desk.SetSystemTrayIcon(theme.FyneLogo())
		img := fyne.NewStaticResource("icon_64.png", bytes)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(img)
	}

	w.SetCloseIntercept(func() {
		w.Hide()
	})

	quickSearch.SetCloseIntercept(func() {
		quickSearch.Hide()
		quickSearchVisible = false
	})

	w.Resize(fyne.NewSize(600, 600))           // Increased height of the main window
	quickSearch.Resize(fyne.NewSize(400, 200)) // Adjusted height of the quick search window
	w.Show()                                   // Open the main window when the application starts
	FocusGlossarSearch(w)

	// Load your own icon resource
	iconData, err := os.ReadFile("icon_64.png")
	if err != nil {
		panic(err)
	}
	icon := fyne.NewStaticResource("icon.png", iconData)

	// Set the app icon
	a.SetIcon(icon)
	w.SetIcon(icon)
	quickSearch.SetIcon(icon)
	a.Settings().SetTheme(&myTheme{})
	a.Run()
}
