# Point System For High Transaction

Implement point system using kafka in golang

## How to run

```bash
# clone this repo
git clone git@github.com:cpartogi/tpoint.git
# or
git clone https://github.com/cpartogi/tpoint.git

# Add kafka to your hosts
sudo nano /etc/hosts
127.0.0.1 kafka

cd tpoint/dbase
nano dbase.sql
copy all create statements to mysql client
run create statements at mysql client

cd tpoint

# running docker compose
docker-compose up

# run rest api
go run main.go

# run consumer
go run cmd/consumer/main.go
```

API Documentation
https://documenter.getpostman.com/view/12167504/T1DqfwAa