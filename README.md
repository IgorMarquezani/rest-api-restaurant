API Rest 

# User Register:
To create a new user send a http POST request containing a JSON that especifies the user that should be created.
\nThere is a example of how a user JSON should be on the directory "JSON_examples/user/register.JSON"

endpoint: localhost:8000/api/user/register

# User Login:
To make login send a http POST request with the same JSON as the one that you should use im Register but,
\nit should only contain the user's password and email. This operation will create a user session if everything ocurrs correctly,
\nso don't forget to add the session cookie named "_SecurePS" that latter will be used for every operation on the API
\nthis endpoint will return a series of informations, if everything have been done correctly, of couse, so beyond the session cookie
\nit also returs a JSON informing the user credentials as well as the room that hi owns (for now it dosen't contain a loot that will be implemented latter on)
\nThere is a example of how a user JSON should be on the directory: "JSON_examples/user/login.JSON"

\nendpoint: localhost:8000/api/user/login

FOR NOW, THE API SUPPORTS THE FOLLOWING REQUESTS COMING FROM A GUEST ON THAT ROOM, BUT THE "ROOM" IS A CONCEPT THAT I SHOULD EXAMPLE BETTER IN THE FUTURE
\nFOR KNOW, JUST SEND REQUESTS AS THE OWNER THAT ROOM
# Product List Register:
To create a new product list send a http POST request containing a JSON that especifies the product list that should be create.
\nThere is an exemple of how a product list JSON should look like on this directory: JSON_example/product_list/register.JSON

endpoint: localhost:8000/api/product_list/register

# Product Register:
To create a new product send a http POST request containing a JSON that especifies the product that will be created
\nThere is an example of how the product JSON should look like on the directory: JSON_examples/products/register.JSON

\nendpoint: localhost:8000/api/product/register
