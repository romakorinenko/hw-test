package model

type Book struct {
	id     int
	title  string
	author string
	year   int
	size   int
	rate   float32
}

func NewBook(id int, title string, author string, year int, size int, rate float32) *Book {
	return &Book{id: id, title: title, author: author, year: year, size: size, rate: rate}
}

func (b *Book) Id() int {
	return b.id
}

func (b *Book) SetId(id int) {
	b.id = id
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) SetTitle(title string) {
	b.title = title
}

func (b *Book) Author() string {
	return b.author
}

func (b *Book) SetAuthor(author string) {
	b.author = author
}

func (b *Book) Year() int {
	return b.year
}

func (b *Book) SetYear(year int) {
	b.year = year
}

func (b *Book) Size() int {
	return b.size
}

func (b *Book) SetSize(size int) {
	b.size = size
}

func (b *Book) Rate() float32 {
	return b.rate
}

func (b *Book) SetRate(rate float32) {
	b.rate = rate
}
