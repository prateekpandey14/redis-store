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

```Go
redisAddr := "localhost:6379"
objectDB := NewRedisObjectDB(redisAddr)
```

4. To validate the code i have written the  `main_test.go` which creates the multiple
   objects of given interfaces and try to store and get the data objects
   simultaneously,  to run the test, execute the below command

```go

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
