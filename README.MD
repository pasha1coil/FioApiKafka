# FioApiKafka
API по определения гендера, пола и национальности по указанному имени
<details>
  <summary>Содержание</summary>
  <ol>
    <li><a href="#установка-и-запуск">Установка и запуск</a></li>
    <li><a href="#реализовано">Реализовано</a></li>
    <li><a href="#примеры запросов">Возникшие вопросы</a></li>
  </ol>
</details>

## Установка и запуск

Клонировать проект.
Возможно понадобиться изменить конфигурации, для этого есть файл .env

Далее через `makefile`:
- `make test` - запустить тесты
- `make build` - запушить в докер образы
- `go run cmd/main.go` - запустить проект, командой в терминале

## Реализовано

Rest методы:
1. Метод добавления человека.
2. Метод получения данных о людях по возрасту.
3. Метод получения данных о людях по полу.
4. Метод получения данных о людях по имени.
5. Метод получения данных о людях по ID.
6. Метод получения данных о всех людях.
7. Метод обновления данных человека.
8. Метод удаления человека.

GraphQL методы:
1. Метод получения данных о людях по возрасту. 
2. Метод получения данных о людях по полу. 
3. Метод получения данных о людях по имени. 
4. Метод получения данных о людях по ID. 
5. Метод получения данных о всех людях.
6. Метод удаления человека.

Также реализованно:
1. Кеширование в Redis
2. Apache Kafka для реализации очереди сообщений с хендлера создания человека, два топика - FIO и FIO_FAILED. В FIO записываются все сообщения, далее сортируются, либо в канал для работы, либо в FIO_FAILED. В FIO_FAILED попадают те сообщения которые получили нулевые значения из api ответа по идентификации пола, возраста и национальности.
3. Тесты для service и repository. в планах покрыть весь код тестами.
4. В качестве веб фреймворка бфл использован gofiber.
5. БД Postgres

## Примеры запросов

Для взаимодействия с сервером есть 7 способов:
1. (POST) /handler/people
2. (GET) /handler/people
3. (GET) /handler/people/:id
4. (GET) /handler/people/name/:name
5. (GET) /handler/people/age/:age
6. (GET) /handler/people/gender/:gender
7. (PUT) /handler/people/:id
8. (DELETE) /handler/people/:id
9. (GraphQL) /graphql

В качестве ответа возвращается JSON с данными.

Примеры запросов и ответов:
1. http://localhost:8080/handler/people POST:
- INPUT
  ```json
    {
    "name":"Alex",
    "surname":"Parker",
    "patronymic":"testovich"
    }
  ```
- OUTPUT
  ```json
  {
    "OK" - ушло в работу
  }
  ```

2. http://localhost:8080/handler/people GET:
- Не предусматривает передачи данных через json, поэтому просто отправляем get запрос
- OUTPUT
  ```json
    [
    {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "male",
        "Age": 46,
        "Nationality": [
            {
                "country_id": "CZ",
                "probability": 0.101
            },
            {
                "country_id": "UA",
                "probability": 0.042
            },
            {
                "country_id": "RO",
                "probability": 0.042
            },
            {
                "country_id": "RU",
                "probability": 0.03
            },
            {
                "country_id": "IL",
                "probability": 0.026
            }
        ]
    }
  ]
  ```

3. http://localhost:8080/handler/people/1 GET:
- Последняя цифра id - человека
- OUTPUT
  ```json
  [
    {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "male",
        "Age": 46,
        "Nationality": [
            {
                "country_id": "CZ",
                "probability": 0.101
            },
            {
                "country_id": "UA",
                "probability": 0.042
            },
            {
                "country_id": "RO",
                "probability": 0.042
            },
            {
                "country_id": "RU",
                "probability": 0.03
            },
            {
                "country_id": "IL",
                "probability": 0.026
            }
        ]
    }
  ]
  ```

4. http://localhost:8080/handler/people/name/Alex GET:
- За последним / - имя пользователя в том формате в котором его создавали, в данном случае первая буква заглавная
- OUTPUT
  ```json
  [
    {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "male",
        "Age": 46,
        "Nationality": [
            {
                "country_id": "CZ",
                "probability": 0.101
            },
            {
                "country_id": "UA",
                "probability": 0.042
            },
            {
                "country_id": "RO",
                "probability": 0.042
            },
            {
                "country_id": "RU",
                "probability": 0.03
            },
            {
                "country_id": "IL",
                "probability": 0.026
            }
        ]
    }
  ]
  ```

5. http://localhost:8080/handler/people/age/46 GET:
- За последнем / - пишем возраст, мне интересен возраст 46 и я получаю:
- OUTPUT
  ```json
    [
    {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "male",
        "Age": 46,
        "Nationality": [
            {
                "country_id": "CZ",
                "probability": 0.101
            },
            {
                "country_id": "UA",
                "probability": 0.042
            },
            {
                "country_id": "RO",
                "probability": 0.042
            },
            {
                "country_id": "RU",
                "probability": 0.03
            },
            {
                "country_id": "IL",
                "probability": 0.026
            }
        ]
    }
  ]
  ```

6. http://localhost:8080/handler/people/gender/male GET:
- За последнем / - пишем пол, мне интересен пол male и я получаю:
- OUTPUT:
  ```json
  [
    {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "male",
        "Age": 46,
        "Nationality": [
            {
                "country_id": "CZ",
                "probability": 0.101
            },
            {
                "country_id": "UA",
                "probability": 0.042
            },
            {
                "country_id": "RO",
                "probability": 0.042
            },
            {
                "country_id": "RU",
                "probability": 0.03
            },
            {
                "country_id": "IL",
                "probability": 0.026
            }
        ]
    }
  ]
  ```

7. http://localhost:8080/handler/people/1 PUT, меняем json ALEX:
- INPUT
  ```json
   {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "female",
        "Age": 23,
        "Nationality": [
            {
                "country_id": "gg",
                "probability": 0.101
            },
            {
                "country_id": "wp",
                "probability": 0.042
            },
            {
                "country_id": "ng",
                "probability": 0.042
            },
            {
                "country_id": "tt",
                "probability": 0.03
            },
            {
                "country_id": "vv",
                "probability": 0.026
            }
        ]
    }
  ```
- OUTPUT
  ```json
  {
        "Name": "Alex",
        "Surname": "Parker",
        "Patronymic": "testovich",
        "Gender": "female",
        "Age": 23,
        "Nationality": [
            {
                "country_id": "gg",
                "probability": 0.101
            },
            {
                "country_id": "wp",
                "probability": 0.042
            },
            {
                "country_id": "ng",
                "probability": 0.042
            },
            {
                "country_id": "tt",
                "probability": 0.03
            },
            {
                "country_id": "vv",
                "probability": 0.026
            }
        ]
    }
  ```
8. http://localhost:8080/handler/people/1 DELETE, меняем json ALEX:
 - делаем запрос - http://localhost:8080/handler/people/1, последний / - id того кого удаляем

9. http://localhost:8080/graphql ALL, тут приведу примеры query и mutation запросов через GraphQL:
- peopleAll
  ```json
   query {
      peopleAll {
      name
      surname
      patronymic
      gender
      age
      nationality {
        countryId
        probability
      }
    }
  }
  ```
- personById
```json
   query {
    personById(id: "1") {
    name
    surname
    patronymic
    gender
    age
    nationality {
      countryId
      probability
      }
    }
  }
  ```
- personsByName
```json
   query {
    personsByName(name: "Alex") {
    name
    surname
    patronymic
    gender
    age
    nationality {
      countryId
      probability
      }
    }
  }
  ```
- personsByAge
```json
   query {
    personsByAge(age: 46) {
    name
    surname
    patronymic
    gender
    age
    nationality {
      countryId
      probability
      }
    }
  }
  ```
- personsByGender
```json
   query {
    personsByGender(gender: "male") {
    name
    surname
    patronymic
    gender
    age
    nationality {
      countryId
      probability
      }
    }
  }
  ```
- deletePerson
```json
   mutation {
    deletePerson(id: "1")
  }
  ```