package appinfo

type CategoryFiter struct {
	Title string `json:"title" query:"title"`
}

type Category struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}
