# test

### Create
```
curl --location --request POST 'localhost/user' \
--header 'Content-Type: application/json' \
--data-raw '{
"id":"12",
"name": "chales13",
"password": "123123",
"status": 2
}'
```

### Update By ID
```
curl --location --request PUT 'localhost/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id":"12",
    "name": "chales13",
    "password": "123123",
    "status": 2
}'
```

### List By Page And Size
```
curl --location --request GET 'localhost/user/list/1/10'
```

### Get By ID
```
curl --location --request GET 'localhost/user/1'
```

### Delete By ID
```
curl --location --request DELETE 'localhost/user/1'
```

直接访问 http://localhost/ 可以注册登陆修改密码

### 遇到的问题

因为几乎每个函数都有 err，所以处理到什么程度