curl -X POST  http://localhost:8080/v1/story \
-H 'Content-Type: application/json' \
--data \
'
{
  "id": "AtEtvd4",
  "title": "Standard Go Project Layout",
  "author": "GoLang",
  "votes": 12,
  "url": "https://github.com/golang-standards/project-layout"
}
'

###

curl http://localhost:8080/v1/story?page=1&limit=2

