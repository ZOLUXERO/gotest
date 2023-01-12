# Como corrrer el proyecto?

este proyecto usa una version de go 1.19.1

para iniciar el proyecto simplemente usar: 
```
go run main.go
```
este proyecto tiene dos dependencias una es kafka y la otra es dynamodb ambas se instalaron de forma local con docker:
- https://hub.docker.com/r/bitnami/kafka
- https://hub.docker.com/r/amazon/dynamodb-local

para iniciar ambos contenedores se hace de la siguiente manera:

```
docker compose up -d
docker pull amazon/dynamodb-local
docker run -p 8000:8000 amazon/dynamodb-local
```

con esto ya podria acceder a dynamodb local con aws-cli si lo tiene, apuntado al --endpoint-url http://localhost:8000 

## los comentarios que vea en el proyecto son temas que no he resuelto del todo.
