package tags

type Config struct {
	Debug bool

	StorageType  string
	TagsJSONFile string

	Encrypt    bool
	PassPhrase [32]byte
}

// TagStorageInterface provides methods for interactions with tags
type TagStorageInterface interface {
	// Get return tag with passed id. If a tag doesn't exist, it returns Tag{}, false
	Get(id int) (Tag, bool)

	// GetAll returns all tags
	GetAll() Tags

	// Add adds a new tag with passed name and color
	Add(name, color, group string)

	// UpdateTag changes name and color of a tag with passed id.
	// If newName/newColor is an empty string, it won't be changed.
	UpdateTag(id int, newName, newColor string) (updatedTag Tag, err error)

	// UpdateGroup changes only group a tag with passed id.
	UpdateGroup(id int, newGroup string) (updatedTag Tag, err error)

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
	Group string `json:"group"`
}