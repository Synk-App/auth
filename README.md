# Auth
Service to manage sessions and authentication from system.

# Getting Started

First step to run project, is to run docker command to build database. This database will need config from SQL file from `_setup` folder. This folder can be found in [project Gateway](https://github.com/Synk-App/gateway). Within README from that project also has more instructions about database setup will help.

So next step is to create a `.env` file in project root and change example values to your config. You can use `example.env` file from `_setup` folder as template.

And then, run `docker compose up -d` into project root to start project.

## Tests

The easy way to run tests is just run `docker compose up -d` command to start project with variables. So, enter in `synk_auth` with `docker exec` and run `go test ./tests -v`.

# Routes

## Get info about app

> `GET` /about

### Response

```json
{
	"ok": true,
	"error": "",
	"info": {
		"server_port": "8080",
		"app_port": "8083",
		"db_working": true
	},
	"list": null
}
```

## Get info about an User

> `GET` /users

### GET Params

```
user_id=1
```

* `user_id`: ID do User desejado, para realizar uma consulta direta

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"user": [
		{
			"user_id": 1,
			"user_name": "Alice Johnson",
			"user_email": "alice.j@example.com",
			"created_at": "25/09/2025 21:19:06"
		}
	]
}
```