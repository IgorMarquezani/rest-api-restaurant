API Rest

# Starting: 
First of all you will need to have Postgres on your machine as well as Golang and Docker.<br>
Enter in the directory "container_files" and the use the command sudo docker-compose build and then sudo docker-compose up,
but you should check if the Postgres isn't running on tou machine already, if it is and is using the same port as the
Postgres that will be created on a container by Docker as especified on: container_files/docker-compose.yml, you should disable the Postgres
already running on you Machine

# User Register:
To create a new user send a http POST request containing a JSON that especifies the user that should be created.

There is a example of how a user JSON should be on the directory "JSON_examples/user/register.JSON"

endpoint: localhost:6000/api/user/register

# User Login:
To make login send a http POST request with the same JSON as the one that you should use im Register but,
it should only contain the user's password and email. This operation will create a user session if everything ocurrs correctly,
so don't forget to add the session cookie named "_SecurePS" that latter will be used for every operation on the API.
this endpoint will return a series of informations, if everything have been done correctly, of couse, so beyond the session cookie
it also returs a JSON informing the user credentials as well as the room that hi owns (for now it dosen't contain a lot of things that will be implemented latter on).

There is a example of how a user JSON should be on the directory: "JSON_examples/user/login.JSON"

endpoint: localhost:6000/api/user/login

FOR NOW, THE API SUPPORTS THE FOLLOWING REQUESTS COMING FROM A GUEST ON THAT ROOM, BUT THE "ROOM" IS A CONCEPT THAT I SHOULD EXPLAIN BETTER IN THE FUTURE,
FOR KNOW, JUST SEND REQUESTS AS THE OWNER THAT ROOM

# Product List Register:
To create a new product list send a http POST request containing a JSON that especifies the product list that should be create.

There is an exemple of how a product list JSON should look like on this directory: JSON_example/product_list/register.JSON

endpoint: localhost:6000/api/product_list/register

# Product Register:
To create a new product send a http POST request containing a JSON that especifies the product that will be created.

There is an example of how the product JSON should look like on the directory: JSON_examples/products/register.JSON

endpoint: localhost:6000/api/product/register
