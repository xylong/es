
```azure
PUT news
    PUT news
    {
        "mappings": {
        "properties": {
        "msgid": {"type":  "keyword"},
        "action": {"type":  "keyword"},
        "from": {"type":  "keyword"},
        "from_userid": {"type":  "keyword"},
        "from_name": {"type":  "keyword"},
        "from_avatar": {"type":  "text"},
        "tolist": {"type":  "text"},
        "roomid": {"type":  "keyword"},
        "msgtime": {"type":  "long"},
        "msgtype": {"type":  "keyword"},
        "url": {"type":  "text"},
        "detail": {"type":  "object"}
    }
  }
}
```

```azure
GET /news/_count
DELETE /news
```
```azure
DELETE /news/_doc/1
```