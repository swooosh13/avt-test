
## Запуск
2. Создать файл `.env` на основе `.env.example`
3. экспортировать перемнные среды в сеансе терминала из файла .env
4. Поднять контейнер бд `make up-db`
5. инициализировать бд из `db/init.sql`
6. Запустить приложение `docker-compose up -d --no-deps --build server`

## Запросы
1. Экспортировать коллекции в postman из папки `postman`
2. [Документация](./swagger/api.yaml)


## Дополнительно  

- Для получения ссылки на месячный отчет было использовано S3 хранилище (в .env открыто только на время теста)
- [Сценарий разрезервирования денег (в Revenue)](./internal/controller/http/transactions.go)
- [Метод получения отчета пользователя с сортировками и пагинацией (в UserReport)](./internal/controller/http/transactions.go)
- [Swagger](./swagger/api.yaml)

### Пример отчета

#### Запрос
`GET /api/v1/transactions?date=2022-11`

#### Ответ
```json
{
    "data": {
        "link": "http://csv.avito.report.hb.bizmrg.com/report_202211_20221024.csv"
    }
}
```











