package database

import (
	"database/sql"
	_ "strconv"

	_ "github.com/lib/pq"
)

// QUERYS FOR USERS
var (
	InsertUser          string = "INSERT INTO users (name, email, passwd, img) VALUES ($1, $2, $3, $4);"
	SearchUserByEmail   string = "SELECT * FROM users WHERE email=$1;"
	SearchUserById      string = "SELECT * FROM users WHERE id=$1;"
	SelectUserByRoom    string = "select (users.*) from users join rooms on rooms.owner = users.id where rooms.id = $1"
	SearchUserBySession string = "select (users.*) from users join users_session on users.id = users_session.who where users_session.securePS = $1;"
)

// QUERYS FOR PRODUCTS
var (
	InsertProduct string = "INSERT INTO products (list_name, list_room, name, price, description, image) values ($1, $2, $3, $4, $5, $6);"
	UpdateProduct string = "update products set name = $1, price = $2, description = $3, image = $4 where list_room = $5 and name = $6;"
	SelectProduct string = "select * from products where name = $1 and list_room = $2;"
)

// QUERYS FOR ROOMS
var (
	SelectRoomByOwner string = "select * from rooms where owner = $1"
	SelectRoomById    string = "select * from rooms where id=$1"
	SelectRoomGuests  string = "select (users.*) from users join guests on guests.user_id = users.id join rooms on rooms.id = guests.inviting_room where rooms.id = $1;"
)

// QUERY FOR INVITES
var (
	SearchInviteByTarget string = "select * from invites where target=$1;"
)

// QUERY FOR USER SESSIONS
var (
	InsertNewSession    string = "insert into users_session (who, active_room, securePS) values ($1, $2, $3);"
	SelectSessionByUId  string = "select * from users_session where who = $1;"
	SelectSessionByHash string = "select * from users_session where securePS = $1;"
	UpdateSessionAcRoom string = "update users_session set active_room = $1 where who = $2"
)

// QUERY FOR GUESTS
var (
	SelectGuestByUserAndRoom string = "select * from guests where user_id = $1 and inviting_room = $2"
)

// MISCELLANEOUS
var (
	InsertProductList string = "insert into product_list values ($1, $2);"
)

// database variables
var (
	db  *sql.DB
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
