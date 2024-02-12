host="http://127.0.0.1:8080"
curl -X POST -H "Content-Type: application/json" \
    -d '{"name": "John Doe", "login": "john@example.com", "password": "password"}' "$host/users/register"


host="http://127.0.0.1:8080"
curl -X POST \
  "$host/users/login" \
  -H 'Content-Type: application/json' \
  -d '{"login": "john@example.com", "password": "password"}'