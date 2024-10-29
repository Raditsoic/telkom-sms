## **Endpoint**

### **Admin**

#### Register Admin

- endpoint: /api/admin/register
- JSON POST: 
```json
{
    "username":"soic",
    "password": "123"
}
```
- Curl: 
```curl
curl --location 'http://localhost:8080/api/admin/register' \
--header 'Content-Type: application/json' \
--data '{
    "username":"soic",
    "password": "123"
}'
```

#### Login Admin (GET Admin Token)
- endpoint: /api/admin/login
- JSON POST: 
```json
{
    "username":"soic",
    "password": "123"
}
```
- Curl: 
```curl
curl --location 'http://localhost:8080/api/admin/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"soic",
    "password": "123"
}'
```

### **Item**

### Get All Items
- endpoint: /api/items
- Curl: 
```curl
curl --location 'http://localhost:8080/api/items'
```

### Create Item
- endpoint: /api/item
- JSON Post:
```json
{
    "name":"Biru",
    "quantity":30,
    "category_id":1
}
``` 
- cURL:
```curl
curl --location 'http://localhost:8080/api/item' \
--header 'Content-Type: application/json' \
--data '{
    "name":"Biru",
    "quantity":30,
    "category_id":1
}
```

### Get Item by ID
- endpoint: /api/item/{id}
- cURL:
```curl
curl --location 'http://localhost:8080/api/item/1'
```

### **Category**

#### Get All Categories
- endpoint: /api/categories
- cURL: curl --location 'http://localhost:8080/api/categories'

### Create Category
- endpoint: /api/category
- JSON Post:
```json
{
    "name":"Pulpen",
    "storage_id":1
}
```
- cURL:
```
curl --location 'http://localhost:8080/api/category' \
--header 'Content-Type: application/json' \
--data '{
    "name":"Pulpen",
    "storage_id":1
}'
```

#### Get Category by ID
- endpoint: /api/category/{id}
- cURL: curl --location 'http://localhost:8080/api/category/1'

#### Get Category by ID w/ Items
- endpoint: /api/category/{id}/items
- cURL: curl --location 'http://localhost:8080/api/category/1/items'


### **Storage**

#### Get All Storage
- endpoint: /api/storages
- cURL
```curl
curl --location 'http://localhost:8080/api/storage'
```

#### Create Storage
- endpoint: /api/storage
- JSON Post:
```json
{
    "name":"ATK",
    "location":"TSO Manyar"
}
```
- cURL:
```curl
curl --location 'http://localhost:8080/api/storage' \
--header 'Content-Type: application/json' \
--data '    {
        "name":"ATK",
        "location":"TSO Manyar"
    }'
```

## **Docker**

### Prerequisite
- Docker Windows/Linux

### Run With Docker
`docker compose up --build`
