# rest-articles (MVP)
## Introduction
The repository contains functionality that serves interacting with an Article object.  
The Article, with plural Articles, is an object with title and URL as properties.   
The object is defined as:
```go
type Article struct {
	Title string `db:"title" json:"title"`
	URL   string `db:"url" json:"url"`
	ID    int64  `db:"id" json:"id" gorm:"PRIMARY_KEY"`

	CreatedOn time.Time  `db:"created_on" json:"-"`
	UpdatedOn *time.Time `db:"updated_on" json:"-"`
	DeletedOn *time.Time `db:"deleted_on" json:"-"`
}
```
An Article has an ID primary key, and tags for JSON and database purposes are introduced.  The ID field is not placed as first in the structure for memory allignment reasons. Unicity indexes were not placed on title or URL due to time constraints.   
An Article can be persisted using a PostgreSQL repository and served to a transport by the application service. The curent transport is done as REST with the Fiber web framework.  
Appication configuration is limited to database name and host due to time constraints.  
Managing the application can be done using the make targets defined.
## Error Handling
The app mainly uses an error type as below with error messages concentrated in an error package. 
```go
type ErrorApplication struct {
	Area      ErrorArea
	AreaError error
	Code      string
	OSExit    *int
}
```
In each area helpers should exist, example:
```go
func (repo *Repository) Errors(repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaRepository],
	}
}

func (repo *Repository) ErrorsWCode(code string, repoError error) *apperrors.ErrorApplication {
	return &apperrors.ErrorApplication{
		Area: apperrors.Areas[apperrors.ErrorAreaRepository],
		Code: code,
	}
}
```
This approach should improve:  
a. the understanding of the application as high level information about the application areas is concentrated in one package  
b. error messages updates as there is only one place which offers a full view of the errors  
c. documentation as the app errors package can be easily shared with other teams
## Logging
Was not added due to time constraints.
## Open API 
Was not added due to time constraints.
## Unit testing
### Prerequisites
Go starting with version 1.19.  
Docker Compose starting with version 3.9.
### Start the database container
```sh
make database-unit
```
### Run unit tests
```sh
make test
```
## Run
Perform clean up if unit testing was done with:
```sh
make infra-cleanup
```
Run:
```sh
make run
```
### Endpoints
Below endpoints would be available on run:
| Endpoint URL |HTTP Verb |Info |
| --------------- | --------------- |------|
| http://localhost:3000/api/v1/article |POST |Create Item|
| http://localhost:3000/api/v1/article/:id |GET |Get Item wth ID|
| http://localhost:3000/api/v1/article/:id |PUT |Update Item wth ID|
| http://localhost:3000/api/v1/article/:id |DELETE |Soft Delete Item wth ID|
| http://localhost:3000/api/v1/articles/all/ |GET |Get all items|
| http://localhost:3000/api/v1/articles?limit=2&page=1 |GET |Get items with pagination|
A `.http` file with requests to these endpoint resides in the `rest` package.
## Stopping
Can be done with Ctrl+C on Linux. Gracefull shutdown would be included in Fiber.
## Benchmark
With `hey` app from https://github.com/rakyll/hey.
```sh
tudi@pad16:~/ram$ ./hey -m GET  -n 10000 "http://localhost:3000/api/v1/article/1"

Summary:
  Total:	19.0677 secs
  Slowest:	0.6172 secs
  Fastest:	0.0002 secs
  Average:	0.0921 secs
  Requests/sec:	524.4481
  
  Total data:	860000 bytes
  Size/request:	86 bytes

Response time histogram:
  0.000 [1]	|
  0.062 [5904]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.124 [800]	|■■■■■
  0.185 [1066]	|■■■■■■■
  0.247 [1170]	|■■■■■■■■
  0.309 [696]	|■■■■■
  0.370 [234]	|■■
  0.432 [84]	|■
  0.494 [27]	|
  0.556 [10]	|
  0.617 [8]	|


Latency distribution:
  10% in 0.0057 secs
  25% in 0.0117 secs
  50% in 0.0320 secs
  75% in 0.1723 secs
  90% in 0.2519 secs
  95% in 0.2920 secs
  99% in 0.3867 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0002 secs, 0.6172 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0180 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0223 secs
  resp wait:	0.0916 secs, 0.0002 secs, 0.6168 secs
  resp read:	0.0002 secs, 0.0000 secs, 0.0204 secs

Status code distribution:
  [200]	10000 responses

```


## Resources
```html
https://www.baeldung.com/ops/docker-compose
https://firehydrant.com/blog/develop-a-go-app-with-docker-compose/
https://stackoverflow.com/questions/52115178/create-a-database-then-table-with-dockerfile
https://golangbot.com/mysql-create-table-insert-row/
https://dev.to/devsmranjan/golang-build-your-first-rest-api-with-fiber-24eh
https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f
https://www.jajaldoang.com/post/golang-function-timeout-with-context/
```