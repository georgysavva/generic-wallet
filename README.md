# Project description  
This is a generic wallet service. It provides following features:
- Send payment from one account to another.  
- See all payments.  
- See all accounts.  
# Implementation details  
- It's a domain driven microservice written in golang with [go-kit](https://github.com/go-kit/kit) library.  
- Service functional available as a RESTful API. See [API docs](https://documenter.getpostman.com/view/865221/S1ETRGPW).  
- Authentication not supported for simplicity sake.  
- Service can be easily auto-scaled, since it's stateless.  
- Postgres is used as persistence layer.  
- Core business logic covered with tests.
- It uses [dep](https://github.com/golang/dep) as dependency management tool.
# Install and run  
1. Download the repo: `git clone https://github.com/gazoon/generic_wallet.git`  
2. Create config.json file in the project root, use config_exemple.json file as an example  
3. Prepare postgres database and run `postgres/db_schema.sql` script to create db schema and add initial accounts  
4. (Optional) Run test `go test generic_wallet/wallet`    
5. Run the server `go run main.go`  
