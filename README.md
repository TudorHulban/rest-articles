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


## Resources
```html
https://www.baeldung.com/ops/docker-compose
https://firehydrant.com/blog/develop-a-go-app-with-docker-compose/
https://stackoverflow.com/questions/52115178/create-a-database-then-table-with-dockerfile
https://golangbot.com/mysql-create-table-insert-row/
https://dev.to/devsmranjan/golang-build-your-first-rest-api-with-fiber-24eh
https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f
```