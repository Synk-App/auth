# Auth
Service to manage sessions and authentication from system.

# Getting Started

First step is to create a `.env` file in project root and change example values to your config. You can use `example.env` file from `_setup` folder as template.

And then, run `docker compose up -d` into project root to start project.

## Tests

The easy way to run tests is just run `docker compose up -d` command to start project with variables. So, enter in `synk_auth` with `docker exec` and run `go test ./tests -v`.

## Certificates

This app must run in HTTPS to authentication works properly. So, to install it, just setup `[mkcert](https://github.com/FiloSottile/mkcert)` into your machine and then run command below into root directory of this project.

```
mkcert -key-file ./.cert/key.pem -cert-file ./.cert/cert.pem localhost synk_auth
```

## Network

You can use a custom network for this services, using then `synk_network` you must create before run it. So, to create on just run command below once during initial setup.

```
docker network create synk_network
```

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

## Register a new User

> `POST` /users/register

### Request

```json
{
	"user_name": "Usuário2",
	"user_email": "usuario2@usuario.com",
	"user_pass": "Usuario2"
}
```

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"user": {
		"user_id": 4,
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJfaWQiOjR9LCJleHAiOjE3NjI3MzA2MTgsImlhdCI6MTc2MjcyOTcxOH0.82riiYls98h3LQHfjE07sOecSr542Opzx_9WqjZaji0"
	}
}
```

## Log in a User

> `POST` /users/login

### Request

```json
{
	"user_email": "usuario2@usuario.com",
	"user_pass": "Usuario2"
}
```

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"user": {
		"user_id": 4,
		"user_name": "Usuário2",
		"user_email": "usuario2@usuario.com",
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJfaWQiOjR9LCJleHAiOjE3NjI3MzE4OTYsImlhdCI6MTc2MjczMDk5Nn0.RYKtwMEt3J9z9v_bwu3mC6Loy_YiduuGm65gxahj44w"
	}
}
```

## Refresh login with `refresh_token`

> `GET` /users/refresh

### Headers

In Headers, must be `refresh_token` header with a valid refresh token generated previously with API.

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"user": {
		"user_id": 4,
		"user_name": "Usuário2",
		"user_email": "usuario2@usuario.com",
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJfaWQiOjR9LCJleHAiOjE3NjI3MzE5MTEsImlhdCI6MTc2MjczMTAxMX0.PsEYIPBMfz4KiH_5TiO-sBZrj-kfdzp4PEai3-Srx0o"
	}
}
```

## Check if `access_token` is valid

> `GET` /users/check
### Headers

```
Authorization: Bearer eyJhbGciOiJIUzI1NiI...
```
### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"user": {
		"user_id": 6
	}
}
```

## Log out an User

> `GET` /users/logout
### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	}
}
```
