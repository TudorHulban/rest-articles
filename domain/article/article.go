package domain

import (
	"fmt"
	"strings"
	"time"
)

type Article struct {
	Title string `db:"title"`
	URL   string `db:"url"`
	ID    int    `db:"id"`

	CreatedOn time.Time  `db:"created_on"`
	UpdatedOn *time.Time `db:"updated_on"`
	DeletedOn *time.Time `db:"deleted_on"`
}

type Articles []*Article

func (art Articles) String() string {
	res := make([]string, len(art)+1, len(art)+1)
	res = append(res, "Articles:")

	for _, article := range art {
		res = append(res, fmt.Sprintf("Title: %s", (*article).Title))
	}

	return strings.Join(res, "\n")
}
