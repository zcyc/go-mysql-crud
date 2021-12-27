# test

# TODO

## 统一的打印方法

```go
package logger

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
```

# Question

## Go 自带的 database/sql 是什么？

提供了一个围绕 SQL（或类似 SQL）数据库的通用接口，必须与数据库驱动程序结合使用。

同时也是一个数据库连接池，可以非常方便的申请和归还资源。

## SQL 连接的释放

有 Close 方法的变量，在使用后要及时调用该方法释放资源。

记不住最后 Close 就在获取之后马上 defer Close。

如果这个函数有返回值，但是你忽略了，返回值依然是存在的。

如果一个返回值需要 Close 才会释放资源，直接忽略就会导致资源泄漏。

## 是否有 SQL 注入风险

正常情况下没有，使用占位符并传入参数会自动屏蔽危险语句。

SQL 千万不要自己用 fmt 或者 "+" 进行拼接。

## 几乎每个函数都有 err，所以处理到什么程度才能结束？

最开始我想的是发给客户端就结束，但是发给客户端也有 err，所以目前是处理到返回给客户端的时候发生错误也要打印出来。

# Usage

## Web

启动服务，直接访问 localhost 端口即可，Web 服务提供注册、登陆、修改

[http://localhost/](http://localhost/)

## CURL

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