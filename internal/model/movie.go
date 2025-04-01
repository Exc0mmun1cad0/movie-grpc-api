package model

type Movie struct {
	ID       string `db:"movie_id"`
	Title    string `db:"title"`
	Genre    string `db:"genre"`
	Director string `db:"director"`
	Year     uint32 `db:"year"`
}
