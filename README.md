# UsersAssets https Server

This is a web service for uploading/downloading data. The user sends a login/password to the server and receives an authorization token (session-id) to access the API. Using the token, the user can upload data to the server and download data from the server.

## Task

It is necessary to develop a web service for uploading/downloading data.

- Improve the authorization mechanism so that at any given time the user has only one (last) session active.

- Limit the maximum user session time to 24 hours. 

- Add IP address data of the authorized user to the DB.

- Implement API methods to get a list of all downloaded files. 
- Implement API methods to delete files.

## Requirements:
- Implement server operation via HTTPS.
- You can only use the standard Go library (latest version) and the pgx library for database access.

## Getting started 
```bash
export  db_conn = "postgres://myuser:mypassword@localhost:5432/mydatabase"

cp your-server-crt server.crt
cp your-server-key server.key

go run cmd/main.go

```