package pagination

const (
	DefaultLimit = 10
	MaxLimit     = 100

	DefaultOffset = 0
)

type Params struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}
