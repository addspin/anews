## Anews - Агрегатор RSS-новостей



## Конфигурация

- **URL сервиса новостей**: http://localhost:8080
  
Приложение принимает конфигурационный файл в формате JSON со следующей структурой:

```json
{
  "server": {
    "port": 8080
  },
  "database": {
    "path": "db/news.db"
  },
  "rss": {
    "feeds": [
      "https://example.com/rss",
      "https://another-site.com/feed"
    ],
    "updatePeriod": 15
  }
}
```

где:
- `port` - порт для запуска HTTP-сервера
- `path` - путь к файлу базы данных SQLite
- `feeds` - массив URL-адресов RSS-лент для мониторинга
- `updatePeriod` - период опроса RSS-лент в минутах   

## API
Получение новостей
- `GET /api/news/` - Получение всех новостей
- `GET /api/news/{limit}` - Получение указанного количества последних новостей