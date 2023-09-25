# Simple project using rabbitMQ for sending emails

## Quick start

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
