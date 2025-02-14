
## Step to Run Backend
Clone the Repository
    
    git clone https://github.com/ReynardChristiansen/BE_TEST_TERRA_DISCOVER.git

Create .env File

    DB_USER="YOUR_DB_USER"
    DB_PASSWORD="YOUR_DB_PASSWORD"
    DB_HOST="YOUR_DB_HOST"
    DB_PORT="YOUR_DB_PORT"
    DB_NAME=article

Migrate Database

    go run migrate.go

Run the Server

    go run main.go

## Endpoints

- Login (POST): http://localhost:8080/login
- Register (GET): http://localhost:8080/register
- Get Article (GET): http://localhost:8080/getArticle
- Create Article (POST): http://localhost:8080/createArticle

## Authorization

Get Article and Create Article require an authorization token provided by the **Login**. The token must be included in the request headers as follows:

    Authorization: Bearer <token>

## Request Body

- login

    To Login, send a POST request with the following body:

        {
            "username": STRING,
            "password": STRING
        }


- Register

    To Register, send a POST request with the following body:

        {
            "username": STRING,
            "password": STRING,
            "email" : STRING
        }

- Create Article

    To Create Article, send a POST request with the following body:

        {
            "title": STRING,
            "content": STRING,
            "category": STRING,
            "status": STRING
        }
