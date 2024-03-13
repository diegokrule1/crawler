package walker

type Url struct {
	Id     string `db:"id"`
	Scheme string `db:"scheme"`
	Host   string `db:"host"`
	Path   string `db:"path"`
}

type State string

const (
	Created    State = "created"
	Processing State = "processing"
	Processed  State = "processed"
)
