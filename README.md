# GoTin

## Dependencies
* Go >= 1.16
* make
* sqlite3

## Build (Linux/Mac)
```bash
make build_unx
```
## Build (Windows)
```bash
make build_win || go build -o gotin.exe main.go
```

## Usage
* Получите токен в приложении тинькофф инвестиции (Sandbox)
```bash
$ gotin init --token <ваш токен> [--dbpath <путь до бд>]
$ gotin add --ticker YNDX
$ gotin add --ticker SPCE
$ gotin delete --ticker YNDX
$ gotin get --ticker SCPE
$ gotin get --verbose
$ gotin get --dbpath <путь до бд>
```

v0.1.0