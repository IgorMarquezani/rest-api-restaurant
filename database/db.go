package database

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

const connection = "host=localhost user=root dbname=restaurant password=bloodyroots port=5432 sslmode=disable"

// database variables
var db *sql.DB

func MustNewConnection() *sql.DB {
	var err error

	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	return db
}

func GetConnection() *sql.DB {
	return db
}

const (
	// Not Found
	ErrNotfound string = "Not found"
)

const (
	// QUERYS FOR USERS
	// INSERT INTO users (name, email, passwd, img) VALUES ($1, $2, $3, $4);
	InsertUser string = "INSERT INTO users (name, email, passwd, img) VALUES ($1, $2, $3, $4);"
	// SELECT * FROM users WHERE email = $1;
	SearchUserByEmail string = "SELECT * FROM users WHERE email=$1;"
	// SELECT * FROM users WHERE id = $1;
	SearchUserById string = "SELECT * FROM users WHERE id=$1;"
	// select (users.*) from users join rooms on rooms.owner = users.id where rooms.id = $1
	SelectUserByRoom string = "select (users.*) from users join rooms on rooms.owner = users.id where rooms.id = $1"
	// select (users.*) from users join users_session on users.id = users_session.who where users_session.securePS = $1;
	SearchUserBySession string = "select (users.*) from users join users_session on users.id = users_session.who where users_session.securePS = $1;"

	// QUERYS FOR PRODUCTS
	// INSERT INTO products (list_name, list_room, name, price, description, image) values ($1, $2, $3, $4, $5, $6);
	InsertProduct string = "INSERT INTO products (list_name, list_room, name, price, description, image) values ($1, $2, $3, $4, $5, $6);"
	// update products set name = $1, price = $2, description = $3, image = $4 where list_room = $5 and name = $6;
	UpdateProduct string = "update products set name = $1, price = $2, description = $3, image = $4 where list_room = $5 and name = $6;"
	// delete from products where name = $1 and list_room = $2;
	DeleteProduct string = "delete from products where name = $1 and list_room = $2;"
	// select * from products where name = $1 and list_room = $2;
	SelectProduct string = "select * from products where name = $1 and list_room = $2;"
	// select (products.*) from products join product_list on products.list_name = product_list.name and products.list_room = product_list.origin_room where product_list.name = $1 and product_list.origin_room = $2;
	SelectProductsByList string = "select (products.*) from products join product_list on products.list_name = product_list.name and products.list_room = product_list.origin_room where product_list.name = $1 and product_list.origin_room = $2 order by products.name;"
	// select (products.*) from products join rooms on products.list_room = rooms.id where rooms.id = $1:
	SelectProductsByRoom string = "select (products.*) from products join rooms on products.list_room = rooms.id where rooms.id = $1 order by products.name;"

	// QUERYS FOR ROOMS
	// insert into rooms (owner) values ($1);
	InsertRoom string = "insert into rooms (owner) values ($1);"
	// select * from rooms;
	SelectAllRooms string = "select * from rooms order by id;"
	// select * from rooms where owner = $1;
	SelectRoomByOwner string = "select * from rooms where owner = $1"
	// select * from rooms where id = $1;
	SelectRoomById string = "select * from rooms where id=$1"
	// select (users.*) from users join guests on guests.user_id = users.id join rooms on rooms.id = guests.inviting_room where rooms.id = $1;
	SelectRoomGuests string = "select (users.*) from users join guests on guests.user_id = users.id join rooms on rooms.id = guests.inviting_room where rooms.id = $1;"
	// select (rooms.*) from guests join rooms on guests.inviting_room = rooms.id where user_id = $1;
	SelectGuestRooms string = "select (rooms.*) from guests join rooms on guests.inviting_room = rooms.id where user_id = $1;"

	// QUERY FOR INVITES
	// insert into invites (target, inviting_room, permission) values ($1, $2, $3);
	InsertInvite string = "insert into invites (target, inviting_room, permission) values ($1, $2, $3)"
	// select * from invites where target=$1;
	SearchInviteByTarget string = "select * from invites where target=$1;"
	// select * from invites where target=$1 and inviting_room = $2;
	SearchInvite string = "select * from invites where target=$1 and inviting_room = $2;"

	// QUERY FOR USER SESSIONS
	// insert into users_session (who, active_room, securePS) values ($1, $2, $3);
	InsertNewSession string = "insert into users_session (who, active_room, securePS) values ($1, $2, $3);"
	// select * from users_session where who = $1;
	SelectSessionByUId string = "select * from users_session where who = $1;"
	// select * from users_session where securePS = $1;
	SelectSessionByHash string = "select * from users_session where securePS = $1;"
	// update users_session set active_room = $1 where who = $2;
	UpdateSessionAcRoom string = "update users_session set active_room = $1 where who = $2"

	// QUERY FOR GUESTS
	// insert into guests values ($1, $2, $3);
	InsertGuest string = "insert into guests values ($1, $2, $3);"
	// select * from guests where user_id = $1 and inviting_room = $2;
	SelectGuestByUserAndRoom string = "select * from guests where user_id = $1 and inviting_room = $2"

	// QUERY FOR PRODUCTS LIST
	// insert into product_list values ($1, $2);
	InsertProductList string = "insert into product_list values ($1, $2);"
	// select * from product_list where name = $1 and origin_room = $2;
	SelectProductList string = "select * from product_list where name = $1 and origin_room = $2;"
	// select * from product_list where origin_room = $1;
	SelectProductListByRoom string = "select * from product_list where origin_room = $1;"

	// QUERY FOR TABS
	// select * from tabs where number = $1 and room = $2;
	SelectTab string = "select * from tabs where number = $1 and room = $2 order by number asc;"
	// select * from tabs where room = $1 order by number asc;
	SelectTabsInRoom string = "select * from tabs where room = $1 order by number asc"
	// insert into tabs (number, room, pay_value, time_maded, table_number) values ($1, $2, $3, $4, $5);
	InsertTab string = "insert into tabs (number, room, pay_value, time_maded, table_number) values ($1, $2, $3, $4, $5)"
	// update tabs set table_number = $1 where number = 2$ and room = 3$;
	UpdateTab string = "update tabs set table_number = $1 where number = $2 and room = $3;"
	// delete from tabs where number = $1 and room = $2;
	DeleteTab string = "delete from tabs where number = $1 and room = $2;"
	// select max(number)+1 from tabs where room = $1;
	SelectMaxTabId string = "select max(number)+1 from tabs where room = $1"

	// QUERY FOR REQUESTS
	// insert into requests values ($1, $2, $3, $4, $5);
	InsertRequest string = "insert into requests values ($1, $2, $3, $4, $5);"
	// select * from requests where product_name = $1 and tab_number = $2 and tab_room = $3;
	SelectRequest string = "select * from requests where product_name = $1 and tab_number = $2 and tab_room = $3;"
	// delete from requests where tab_room = $1 and tab_number = $2 and product_name = $3;
	DeleteRequest string = "delete from requests where tab_room = $1 and tab_number = $2 and product_name = $3;"
	// delete from requests where tab_room = $1 and tab_number = $2;
	DeleteRequestsInTab string = "delete from requests where tab_room = $1 and tab_number = $2;"
	// select * from requests where tab_room = $1 and tab_number = $2 order by tab_number asc;
	SelectRequestsInTab string = "select * from requests where tab_room = $1 and tab_number = $2 order by tab_number asc;"
	// update requests set quantity = $1 where product_name = $2 and tab_number = $3 and tab_room = $4;
	UpdateRequestQuantity string = "update requests set quantity = $1 where product_name = $2 and tab_number = $3 and tab_room = $4;"
)

func IsDuplicateKeyError(warning string) bool {
	msg := strings.Split(warning, " ")
	if len(msg) < 2 {
		return false
	}

	if msg[1] == "duplicate" && msg[2] == "key" {
		return true
	}

	return false
}

func IsAuthenticantionFailedError(warning string) bool {
	msg := strings.Split(warning, " ")
	if len(msg) < 3 {
		return false
	}

	if msg[1] == "Ident" && msg[2] == "authentication" && msg[3] == "failed" {
		return true
	}

	return false
}

func IsForeignKeyConstraintError(warning string) bool {
	error := [4]string{"violates", "foreign", "key", "constraint"}
	msg := strings.Split(warning, " ")
	if len(msg) < 3 {
		return false
	}

	for i := range error {
		if msg[i+7] != error[i] {
			return false
		}
	}

	return true
}
