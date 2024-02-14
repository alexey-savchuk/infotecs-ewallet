# Текущее состояние решения
Доступна рабочая версия решения (см. [релизы](https://github.com/alexey-savchuk/infotecs-ewallet/releases)). Формулировка задания может быть найдена в [вики](https://github.com/alexey-savchuk/infotecs-ewallet/wiki/%D0%A4%D0%BE%D1%80%D0%BC%D1%83%D0%BB%D0%B8%D1%80%D0%BE%D0%B2%D0%BA%D0%B0-%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D1%8F) проекта.

## Как запустить проект
```sh
docker-compose up --build
```

## Как запустить тесты
```sh
go mod download
go test ./...
```
