# web-service-go-oracle

Este es un web services especifico para oracle pl/sql

## Llamada a la base de datos

```go
  var db *sql.DB
	db, errdb := sql.Open("godror", "NOTESUSER/admin@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=tcp)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SID=xe)))")
```

## Llamada a un procedimiento almacenado

```go
  db.Exec("BEGIN NewNote (:1,:2,:3,:4); end;",temp.Title, temp.Description,sql.Out{Dest:&error1},sql.Out{Dest:&error2})
```

## Llamada a un select

```go
  rows, _ := db.Query("SELECT id, title, description FROM notes")
  var note Note
  for rows.Next() {
    rows.Scan(&note.Id, &note.Title, &note.Description)
    notes=append(notes,note)
  }
```

## Aplicación de prueba

- [Aplicación hecha en flutter](https://github.com/walker1239/Notes-Aplication)
