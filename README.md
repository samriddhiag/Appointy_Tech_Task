# Appointy Technical Task

### All required routes are working.

Setup the latest version of golang on your machine.

Commands to run:
go run insta_users_posts.go

After initialising the localhost server:
curl localhost:8080/users -X -d '{"name" : "Appointy", "id" : "app1", "email" : "appointy@gmail.com", "password" : "apptechtask"}' -H "Content-Type: application/json"

JSON object will be returned. Eg:
User object - {"name" : "Appointy", "id" : "app1", "email" : "appointy@gmail.com", "password" : "apptechtask"}






