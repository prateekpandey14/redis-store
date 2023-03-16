### Build And Run locally

To create a test Redis server to validate the above test code, you can use a 
Redis Docker container. Here are the steps to do so:

1. Pull the Redis Docker image:

```sh
docker pull redis
```

2. Run a Redis container using the Redis Docker image:
```sh
docker run -p 6379:6379 --name redis-test -d redis
```

3. Build and run locally
```sh

$ go build -o redis-store


$ ./redis-store 
got person by ID: "1": 
&{ID:1 Name:Alice Age:30}
got animal by name: "Fluffy":  
&{ID:3 Name:Fluffy Type:cat}
list of person type of objects
&{ID:1 Name:Alice Age:30}
&{ID:2 Name:Bob Age:25}
Successfully delete animal type of object 4
```

4. Verify the redis keys entries and info


```sh
127.0.0.1:6379> KEYS *
1) "person:1"
2) "person:2"
3) "animal:3"
3) "animal:4"


127.0.0.1:6379> GET person:1
"{\"ID\":\"1\",\"Name\":\"Alice\",\"Age\":30}"

```




### Test implementation

To create a test Redis server to validate the above test code, you can use a 
Redis Docker container. Here are the steps to do so:

1. Pull the Redis Docker image:

```sh
docker pull redis
```

2. Run a Redis container using the Redis Docker image:
```sh
docker run -p 6379:6379 --name redis-test -d redis
```


3. In your Go code, use the Redis connection string to connect to the Redis server running in the container:
for example:

```go
redisAddr := "localhost:6379"
objectDB := NewRedisObjectDB(redisAddr)
```

4. To validate the code i have written the  `main_test.go` which perform the tests in all scenarios,
creates the multiple objects like person and animal which implements the interfaces and perform the below test 
simultaneously.

- Store will store the object in the redis data store 
- GetObjectByID will retrieve the object with the provided ID
- GetObjectByName will retrieve the object with the given name.
- ListObjects will return a list of all objects of the given kind.
- DeleteObject will delete the object with provided ID

 to run the test, execute the below command

```sh

go test -v ./...
=== RUN   TestRedisObjectDB
--- PASS: TestRedisObjectDB (0.01s)
PASS
ok  	github.com/prateekpandey14/redis-stores  0.018s
```


4. In this example, we're assuming that the Redis container is running on the same machine
as the Go code and using the default Redis port of 6379. If your Redis container is running
on a different machine or using a different port, update the `redisAddr` variable accordingly.

When you're finished testing, stop and remove the Redis container:

```sh
docker stop redis-test
docker rm redis-test


```
