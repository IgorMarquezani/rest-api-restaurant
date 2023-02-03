package database

import (
	"database/sql"
	_ "strconv"

	_ "github.com/lib/pq"
)

// QUERYS FOR USERS
var (
  InsertUser        string = "INSERT INTO users (name, email, passwd, img) VALUES ($1, $2, $3, $4);"
	SearchUserByEmail string = "SELECT * FROM users WHERE email=$1;"
	SearchUserById    string = "SELECT * FROM users WHERE id=$1;"
  SelectUserByRoom  string = "select (users.*) from users join rooms on rooms.owner = users.id where rooms.id = $1" 
)

// QUERYS FOR PRODUCTS
var (
  InsertProduct     string = "INSERT INTO products (list_name, list_room, name, price, description, image) values ($1, $2, $3, $4, $5, $6);"
  UpdateProduct     string = "update products set name = $1, price = $2, description = $3, image = $4 where list_room = $5 and name = $6;"
  SelectProduct     string = "select * from products where name = $1 and list_room = $2;"
)

// QUERYS FOR ROOMS
var (
  SelectRoomByOwner string = "select * from rooms where owner = $1"
  SearchRoomById    string = "select (owner) from rooms where id=$1"
  SelectRoomGuests  string = "select (users.*) from users join guests on guests.user_id = users.id join rooms on rooms.id = guests.inviting_room where rooms.id = $1;"
)

// QUERY FOR INVITES
var (
  SearchInviteByTarget string = "select * from invites where target=$1;"
)

// QUERY FOR USER SESSIONS
var (
  InsertNewSession    string = "insert into users_session (who, securePS) values ($1, $2);"
  SelectSessionByUId  string = "select * from users_session where who = $1;"
  UpdateSessionAcRoom string = "update users_session set active_room = $1"
)

// MISCELLANEOUS
var (
  InsertProductList string = "insert into product_list values ($1, $2);"
)

// database variables
var (
  db *sql.DB 
  err error 
)

func NewConnection() {
  connection := "host=localhost user=root dbname=restaurant password=bloodyroots port=5432 sslmode=disable "
  db, err = sql.Open("postgres", connection)

  if err != nil {
    panic(err)
  }
}

func GetConnection() *sql.DB {
  return db
}

/*func Search(query string, args ...any) (interface{}, error) {
  var finalQuery []byte 
  var key int = 1

  values := make(map[string]interface{})

  for i := 0; i < len(query); i++ {
    if query[i] == '$' {
      values[string(query[i]) + strconv.Itoa(key)] = args[key]

      start := []byte(query[:i])
      end := []byte(query[i+1:])


      key++
    } 
  }

  search, err := db.Query(query)  
  if err != nil {
    return nil, err
  }

  columns, err := search.ColumnTypes()
  if err != nil {
    return nil, err
  }


  return nil, nil
} */
