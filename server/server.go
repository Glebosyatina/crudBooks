package server

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"
    "encoding/json"
    "strconv"
)

type Book struct{
    Id int  `json:"id"`
    Name string `json:"name"`
    Author string`json:"author"`
}

func Run(){
    
    http.HandleFunc("/books", getListBooks)
    http.HandleFunc("/books/add", addBook)
    http.HandleFunc("/books/delete", deleteBook)
    http.HandleFunc("/books/update", updateBook)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
       fmt.Fprintf(w, "Hello") 
    })

    http.ListenAndServe(":6759", nil)
}


//return books in json array [{},{}]
func getListBooks(w http.ResponseWriter, r *http.Request){
//    fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
    if r.Method != http.MethodGet{
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method not allowed")
        return
    }

    books := []Book{}
    connStr := "user=postgres password=postgres dbname=library sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil{
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM books")
    if err != nil{
        panic(err)
    }
    for rows.Next(){
        p := Book{}
        err := rows.Scan(&p.Id, &p.Name, &p.Author)
        if err != nil{
            panic(err)
        }
        books = append(books, p)
    }
    for _, b := range books{
        fmt.Println(b.Id, b.Name, b.Author)
    }
   //сделать енкодинг в json
   data, err := json.Marshal(books)
   fmt.Fprint(w, string(data))

}

//add book from request json body
func addBook(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost{
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method not allowed")
        return
    }


   //парсим тело запроса в структуру
    b := Book{} 
    err := json.NewDecoder(r.Body).Decode(&b)
    if err != nil { panic(err) }

    fmt.Println(b)
    //пишем в бд
    connStr := "user=postgres password=postgres dbname=library sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil { panic(err) }
    defer db.Close()

    res, err := db.Exec("INSERT INTO books (name, author) VALUES ($1, $2)", b.Name, b.Author)
    if err != nil { panic(err) }

    fmt.Println(res.RowsAffected())

}

//delete book by id in url &id=someid
func deleteBook(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodDelete{
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method not allowed")
        return
    }
   id, err := strconv.Atoi(r.URL.Query().Get("id")) 
   if err != nil || id < 1 { 
       http.NotFound(w, r)
       return 
   }
    connStr := "user=postgres password=postgres dbname=library sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil { panic(err) }
    defer db.Close()

    res, err := db.Exec("DELETE FROM books WHERE id=$1", id)
    if err != nil { panic(err) }
    fmt.Println(res.RowsAffected())
}

//update information about book by id in url and json body
func updateBook(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPut{
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method not allowed")
        return
    }
   //получили id из url
   id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil { panic(err) }
   //распарсили тело запроса
    b := Book{}
    err = json.NewDecoder(r.Body).Decode(&b)
    if err != nil { panic(err) }
     
    //подключаемся к бд и меняем запись о книге
    connStr := "user=postgres password=postgres dbname=library sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    
    res, err := db.Exec("UPDATE books SET name=$1, author=$2 WHERE id=$3", b.Name, b.Author, id)
    fmt.Println(res.RowsAffected())
}
