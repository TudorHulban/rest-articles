# rest-articles
## Unit testing
### Docker Postgres
### Create Docker image
```sh
docker build - < Dockerfile_Postgres -t database
```
```sh
sudo docker run -d --name=co-postgres -p 5432:5432 -e POSTGRES_PASSWORD=thepassword postgres
```
### Database Objects
Run migrations:
```sh
migrate -path migrations/ -database postgres://postgres:thepassword@127.0.0.1:5432/rest?sslmode=disable -verbose up 2
```
## Prerequisites
### Docker ( Compose )
## Infrastructure
Use below command to create infrastructure:
```sh
make infra
```
## Test Go Docker file
### Create Docker image
```sh
sudo docker build -t goapp Dockerfile_GoApp
```
### Create container from the Docker image created
```sh
sudo docker run -d --name=co-goapp -p 3000:3000 goapp
```
## Benchmark
With key app from https://github.com/rakyll/hey.
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