### Stack :

-- Golang (Tested on v1.13.4) with go modules
-- MySQL v5.7
-- Redis (Tested on v4.0.9)

### Installation : 
-- Clone the app
-- Install dependencies mentioned above
-- Add the enviroment variables below
-- run 
```sh
go mod tidy
``````sh
go run main.go
```

# Enviroment 

> Set environment variable GO_BOILER_ENV to "dev", for linux : 

```sh
sudo -H nano /etc/environment
```
> Put this at the end : GO_BOILER_ENV="dev"

> add the variable values in env/.env.dev