# rest-articles
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
sudo docker build -t goapp .
```
### Create container from the Docker image created
```sh
sudo docker run -d --name=co-goapp -p 3000:3000 goapp
```

## Resources
```html
https://www.baeldung.com/ops/docker-compose
https://firehydrant.com/blog/develop-a-go-app-with-docker-compose/
```