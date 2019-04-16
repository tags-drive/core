package cmd

// TagStorageInterface provides methods for interactions with tags
type TagStorageInterface interface {
	// Get return tag with passed id. If a tag doesn't exist, it returns Tag{}, false
	Get(id int) (Tag, bool)

	// GetAll returns all tags
	GetAll() Tags

	// Add adds a new tag with passed name and color
	Add(name, color string)

	// Change changes a tag with passed id.
	// If pass empty newName (or newColor), field Name (or Color) won't be changed.
	Change(id int, newName, newColor string) (updatedTag Tag, err error)

	// Delete deletes a tag with passed id
	Delete(id int)

	// Check checks is there tag with passed id
	Check(id int) bool

	// Shutdown gracefully shutdown TagStorage
	Shutdown() error
}

// Tags is a map of Tag
type Tags map[int]Tag

// Tag contains the information about a tag
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
