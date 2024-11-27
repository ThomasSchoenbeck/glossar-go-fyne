import hooks

```go
	hook "github.com/robotn/gohook"
```


# glossar search

- allows searching for tags by typing `tag:` and then the tag you are looking for.
- search for string contained in name
- all search is performed converted to lowercase



todo:
- functionality to update tags
- better styling
- update readme

known issues:
- app icon (window and taskbar) is not working
- if new tag is added on an element but not saved, it can still appear in the tag tab in the tag list
- list elements cannot be unselected, because a selected element does not re-trigger the select event on another click
- adding new element with two tags seems only last tag is actually saved to overall list of tags