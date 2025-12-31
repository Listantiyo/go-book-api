# API Testing Collection

## 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

## 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the token from response!

## 3. Create Book (with token)
```bash
TOKEN="your_token_here"

curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "isbn": "9780132350884",
    "description": "A handbook of agile software craftsmanship",
    "stock": 10
  }'
```

## 4. Get All Books (public)
```bash
curl http://localhost:8080/api/v1/books?page=1&page_size=10
```

## 5. Borrow Book (with token)
```bash
curl -X POST http://localhost:8080/api/v1/borrows \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": 1}'
```

## 6. Get My Borrows (with token)
```bash
curl http://localhost:8080/api/v1/borrows/me \
  -H "Authorization: Bearer $TOKEN"
```

## 7. Return Book (with token)
```bash
curl -X POST http://localhost:8080/api/v1/borrows/return \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"borrow_id": 1}'
```