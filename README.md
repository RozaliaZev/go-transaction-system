# go-transaction-system

## Схема работы

![alt text](https://sun9-58.userapi.com/impg/RvqtuViduZ97lHbiYNd_MTdBxlD9bSd11L2VCA/zOXGWaBoHds.jpg?size=1366x250&quality=95&sign=4e2679ae872e12ad55747cc2182b7f28&type=album)

## Конфигурирование сервиса

Конфигурирование приложения происходит через docker-compose

- **db** информация для подключения к базе данных PostgreSQL

```sh
db:
  HOST: localhost
  PORT: 5432
  USER: user
  PASSWORD: password
  DBNAME: base
  SSLMODE: disable
```

## Запросы

Можно отправить два вида запросов (см файл swagger)

- **POST** загрузить новую транзакцию, указав данные о ней, а в ответ получить её номер из бд.

```sh
/transaction
```

```sh
{
  "amount": 0,
  "currency": "string",
  "sender_card_number": "stringstringstri",
  "recipient_phone_number": "stringstrin"
}
```

- **GET** запросить данные об определенной транзакции: основную информацию и время, когда она была совершена.

```sh
/transaction/{id}
```

## Замечания

- **проблемы:** не приходит ответ с get запроса (косяк с типами данных, скоро исправлю)  :(
- **доработки:** в будущем нужно прикрутить брокер


