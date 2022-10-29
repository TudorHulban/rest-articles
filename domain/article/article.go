package domain

import (
	"fmt"
	"strings"
	"time"
)

type Article struct {
	Title string `db:"title" json:"title"`
	URL   string `db:"url" json:"url"`
	ID    int64  `db:"id" json:"id"`

	CreatedOn time.Time  `db:"created_on" json:"-"`
	UpdatedOn *time.Time `db:"updated_on" json:"-"`
	DeletedOn *time.Time `db:"deleted_on" json:"-"`
}

type Articles []*Article

func (art Articles) String() string {
	res := []string{"Articles:"}

	for _, article := range art {
		res = append(res, fmt.Sprintf("Title: %s", (*article).Title))
	}

	return strings.Join(res, "\n")
}
