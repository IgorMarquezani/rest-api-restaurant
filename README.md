API Rest

# Starting: 
For starting the application run the bash script in "container-files", called "start.sh" and pass the argument "start". Any question
about the script functionalities pass the argument "help" in it. If you wanna run the front-end as well, you can find it here: 
https://github.com/IgorMarquezani/front-end-restaurant

# User Register:
To create a new user send a http POST request containing a JSON that especifies the user that should be created.

There is a example of how a user JSON should be on the directory "JSON_examples/user/register.JSON"

endpoint: localhost:6000/api/user/register

# User Login:
To make login send a http POST request with the same JSON as the one that you use it in register but,
containing only the user's password and email. This operation will create a user session if everything ocurrs correctly,
so don't forget to add the session cookie named "_SecurePS" that latter will be used for every operation on the API.

There is a example of how a user JSON should be on the directory: "JSON_examples/user/login.JSON"

endpoint: localhost:6000/api/user/login

# Product List Register:
To create a new product list send a http POST request containing a JSON that especifies the product list that should be created.

There is an exemple of how a product list JSON should look like on this directory: JSON_example/product_list/register.JSON

endpoint: localhost:6000/api/product_list/register

# Product Register:
To create a new product send a http POST request containing a JSON that especifies the product that will be created.

There is an example of how the product JSON should look like on the directory: JSON_examples/products/register.JSON

endpoint: localhost:6000/api/product/register
