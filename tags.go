package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var tagSearchLabel *widget.Label

func SetupTags(w fyne.Window) *container.TabItem {
	tagListContainer := container.NewVBox()
	setTags(tags, tagListContainer)

	tagSearchEntry := widget.NewEntry()
	tagSearchEntry.OnChanged = func(s string) {
		updateTagList(s, tagListContainer)
	}

	tagListScroll := container.NewVScroll(tagListContainer)
	tagListScroll.SetMinSize(fyne.NewSize(400, 200)) // Set minimum size to increase height

	tagSearchLabel = widget.NewLabel(fmt.Sprintf("Search Tags: (%d of %d)", len(tagListContainer.Objects), len(tags)))

	return container.NewTabItem("Tags", container.NewVBox(
		tagSearchLabel,
		tagSearchEntry,
		tagListScroll,
	))
}

func setTags(tags []string, container *fyne.Container) {
	container.Objects = nil
	for _, tag := range tags {
		richText := widget.NewRichTextFromMarkdown(tag)
		container.Add(richText)
	}
	container.Refresh()

}

func updateTagList(searchString string, container *fyne.Container) {
	container.Objects = nil

	lowerSearchTerm := strings.ToLower(searchString)
	if len(searchString) > 0 {

		for _, entry := range tags {
			if strings.Contains(strings.ToLower(entry), lowerSearchTerm) {
				highlightedTerm := highlightMatch(entry, lowerSearchTerm)
				result := widget.NewRichTextFromMarkdown(highlightedTerm)
				container.Add(result)
			}
		}
	} else {

		if container.Objects == nil {
			for _, tag := range tags {
				container.Add(widget.NewRichTextFromMarkdown(tag))
			}
		}
	}

	tagSearchLabel.SetText(fmt.Sprintf("Search Tags: (%d of %d)", len(container.Objects), len(tags)))
	container.Refresh()
}
