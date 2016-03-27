golang-pool
====

## pool.go

pool.go is a pool implementation for golang. You may need it if you use golang with version <= 1.2 

Its features include:

* Network failure tolerance.
* Establishing and closing connections on demand depending on the network state.

# Installation

Just run `go get github.com/jinntrance/pool`

## Examples

Initial a pool
```go
var pool = Pool{
        //how to create a client and then put it into the pool
        New: func() (interface{}, error) {
                cli, err := CreateAClient(GetAServer())
                return cli, err
        },
        //do something when closing the client
        Close: func(x interface{}) {
                x.(*Client).Close()
        },
        PoolSize: 100, //MaxNum of clients to retain
}

```
Where `Client` is the type for your sepecific client.
Then use the pool just created
```go
cli, err := pool.Get()
if nil == err {
        client := cli.(*Client) //get a client and convert it to a specific type
        re, lastError := client.doSomething()
        err = lastError //record the last error
} 
pool.Put(cli, err) // put back the client with the error. If error occurs and 'err' is not nil, 'cli' would be closed.
```

<script type="text/javascript" src="http://www.josephjctang.com/assets/js/analytics.js" async="async"></script>
