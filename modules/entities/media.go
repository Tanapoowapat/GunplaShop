package entities

type Images struct {
	Id       string `db:"id"`
	FileName string `db:"filename"`
	Url      string `db:"url"`
}
