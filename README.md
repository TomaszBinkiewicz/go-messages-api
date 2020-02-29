# Messages API
**This is a Golang application which allows to**

* Display all messages
* Display messages filtered by email address
* Create new message
* Send messages

Moreover messages older than 5 minutes are automatically deleted.

## REST API endpoints

### Get list of messages
*Request*

`GET /api/messages`


Returns all messages from database.

### Get list of messages filtered by email address
*Request*

`GET /api/messages/{emailValue}`

Returns all messages with given email address.

If email address is incorrect (e.g. jan.kowalski.example.com) returns `HTTP 400 response status code`

### Create new message
*Request*

`POST /api/message`
```
curl -X POST localhost:8000/api/message -d '{"email":"jan.kowalski@example.com","title":"Interview","content":"simple text","magic_number":101}'
```

The App will store new message into cassandra database.

If data is incorrect returns `HTTP 400 response status code`

### Send message(s)
*Request*

`POST /api/send`
```
curl -X POST localhost:8000/api/send -d '{"magic_number":101}'
```

The App will send messages with given magic_number and delete them from database.

If data is incorrect (e.g. {"magic_number":"not_a_number"}) returns `HTTP 400 response status code`

