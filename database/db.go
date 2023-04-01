package database

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

const (
	// QUERYS FOR USERS
	InsertUser          string = "INSERT INTO users (name, email, passwd, img) VALUES ($1, $2, $3, $4);"
	SearchUserByEmail   string = "SELECT * FROM users WHERE email=$1;"
	SearchUserById      string = "SELECT * FROM users WHERE id=$1;"
	SelectUserByRoom    string = "select (users.*) from users join rooms on rooms.owner = users.id where rooms.id = $1"
	SearchUserBySession string = "select (users.*) from users join users_session on users.id = users_session.who where users_session.securePS = $1;"

	// QUERYS FOR PRODUCTS
	InsertProduct        string = "INSERT INTO products (list_name, list_room, name, price, description, image) values ($1, $2, $3, $4, $5, $6);"
	UpdateProduct        string = "update products set name = $1, price = $2, description = $3, image = $4 where list_room = $5 and name = $6;"
	DeleteProduct        string = "delete from products where name = $1 and list_room = $2;"
	SelectProduct        string = "select * from products where name = $1 and list_room = $2;"
	SelectProductsByList string = "select (products.*) from products join product_list on products.list_name = product_list.name and products.list_room = product_list.origin_room where product_list.name = $1 and product_list.origin_room = $2;"
	SelectProductsByRoom string = "select (products.*) from products join rooms on products.list_room = rooms.id where rooms.id = $1"

	// QUERYS FOR ROOMS
	SelectAllRooms    string = "select * from rooms;"
	SelectRoomByOwner string = "select * from rooms where owner = $1"
	SelectRoomById    string = "select * from rooms where id=$1"
	SelectRoomGuests  string = "select (users.*) from users join guests on guests.user_id = users.id join rooms on rooms.id = guests.inviting_room where rooms.id = $1;"
	SelectGuestRooms  string = "select (rooms.*) from guests join rooms on guests.inviting_room = rooms.id where user_id = $1;"

	// QUERY FOR INVITES
	SearchInviteByTarget string = "select * from invites where target=$1;"

	// QUERY FOR USER SESSIONS
	InsertNewSession    string = "insert into users_session (who, active_room, securePS) values ($1, $2, $3);"
	SelectSessionByUId  string = "select * from users_session where who = $1;"
	SelectSessionByHash string = "select * from users_session where securePS = $1;"
	UpdateSessionAcRoom string = "update users_session set active_room = $1 where who = $2"

	// QUERY FOR GUESTS
	SelectGuestByUserAndRoom string = "select * from guests where user_id = $1 and inviting_room = $2"

	// QUERY FOR PRODUCTS LIST
	InsertProductList       string = "insert into product_list values ($1, $2);"
	SelectProductList       string = "select * from product_list where name = $1 and origin_room = $2;"
	SelectProductListByRoom string = "select * from product_list where origin_room = $1;"

	// QUERY FOR TABS
	SelectTabsInRoom string = "select * from tabs where room = $1 order by number asc"
	InsertTab        string = "insert into tabs (number, room, time_maded, table_number) values ($1, $2, $3, $4)"
	DeleteTab        string = "delete from tabs where number = $1 and room = $2;"
	SelectMaxTabId   string = "select max(number)+1 from tabs where room = $1"

	// QUERY FOR REQUESTS
	InsertRequest       string = "insert into requests values ($1, $2, $3, $4, $5);"
	DeleteRequestsInTab string = "delete from requests where tab_room = $1 and tab_number = $2;"
	SelectRequestsInTab string = "select * from requests where tab_room = $1 and tab_number = $2"
)

// database variables
var (
	db   *sql.DB
	err  error
	trys int = 1
)

func IsDuplicateKeyError(warning string) bool {
	message := strings.Split(warning, " ")

	if message[1] == "duplicate" && message[2] == "key" {
		return true
	}

	return false
}

func MustNewConnection() {
	connection := "host=localhost user=root dbname=restaurant password=bloodyroots port=5432 sslmode=disable "
	db, err = sql.Open("postgres", connection)

	if err != nil {
		panic(err)
	}
}

func GetConnection() *sql.DB {
	return db
}
