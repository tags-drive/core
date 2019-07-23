package tags

type Config struct {
	Debug bool

	StorageType  string
	TagsJSONFile string

	Encrypt    bool
	PassPhrase [32]byte
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
