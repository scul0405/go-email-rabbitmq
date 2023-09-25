# Simple project using rabbitMQ for sending emails

## Quick start

### Run rabbitMQ in docker
```sh
docker run -d --hostname rmq --name rabbit-server -p 8080:15672 -p 5672:5672 rabbitmq
```

### Run server
```sh
make run
```

### Send a request with POST method
```sh
curl --location 'localhost:3000/emails' \
--header 'Content-Type: application/json' \
--data-raw '{
    "to": ["username@gmail.com"],
    "body": "Hello, this is a test msg",
    "subject": "test subject",
    "contentType": "text"
}'
```
