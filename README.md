# rest-articles
## Database Migrations
a. create Docker container
```sh
sudo docker run -d --name=co-postgres -p 5432:5432  -e POSTGRES_PASSWORD=thepassword postgres
```
b. create database `rest` with DBeaver
c. export Postgres connection
```sh
export POSTGRESQL_URL='postgres://postgres:thepassword@localhost:5432/rest?sslmode=disable'
```
d. run migrations
```sh
migrate create -ext sql -dir migrations -seq create_articles_table
```