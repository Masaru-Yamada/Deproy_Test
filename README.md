# play-security

## Launch the applicaiton
```
$ buffalo dev
```

Please access http://localhost:3000/ 

## Install Buffalo
### Windows

Install scoop 
```
C:\> Set-ExecutionPolicy RemoteSigned -scope CurrentUser
C:\> invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://get.scoop.sh')
```

Buffalo can be installed using the [Scoop](https://scoop.sh/) package manager:

```
C:\> scoop install buffalo
C:\> scoop install gcc
```

### Homebrew (macOS)
Buffalo can be installed using the [Homebrew](https://brew.sh/) package manager:
```
$ brew install gobuffalo/tap/buffalo
```

## Database settings

### Start Postgres
```
$ docker-compose up -d
```

### Install Pop 
Windows:
```
C:\> go install github.com/gobuffalo/pop/v6/soda@latest
```

Homebrew (macOS):
```
$ brew install gobuffalo/tap/pop
```

### Create databases

```
$ buffalo pop create -a
```

### Create tables

```
$ buffalo pop migrate up
```

### How to access the local database 

```
$ docker exec -it play-security-database bin/bash

Container...
$ psql -U postgres

List of databases...
postgres=# \l
```

## Install frontend tools
### Windows
```
C:\> scoop search node
...
C:\> scoop bucket add main
...
C:\> scoop install main/nodejs
...
C:\> npm -v
8.11.0
C:\> npm install --global yarn
...
```

### Homebrew (macOS)
```
$ brew install node
...
$ npm -v
8.11.0
$ npm install --global yarn 
```

-----

ここから先は実行をする必要はありません !!!

## Create the application
### Project

```
$ buffalo new play-security --module github.com/aiit2022-pbl-okuhara/play-security --db-type postgres
```
