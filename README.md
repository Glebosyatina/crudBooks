must be a database library
create table books(
id SERIAL PRIMARY KEY,
name VARCHAR,
author VARCHAR
);

to run app go run .

curl -X GET localhost:6759/books
curl -X POST localhost:6759/books/add -d '{"name":"Бесы", "author":"Достоевский"}'
curl -X DELETE localhost:6759/books/delete?id=1
curl -X PUT localhost:6759/books/delete?id=1 -d '{"name":"someName", "author":"someAuthor"}'
