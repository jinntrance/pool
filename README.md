golang-pool
====

## pool.go

pool.go is a pool implementation for golang.

Its features include:

* Network failure tolerance.
* Establishing and closing connections on demand.

# Installation

Just run `go get github.com/jinntrance/pool`

## Examples

initial a pool
```
var pool = Pool{
        New: func() (interface{}, error) {
                cli, _, err := InitSocket(GetAServer())
                return cli, err
        },
        Close: func(x interface{}) {
                x.(*Client).Close()
        },
        poolSize: 100,
        defaultToWait: true,
}

```
and then use the pool
```
cli, err := pool.Get()
if nil == err {
        client := cli.(*Client) //get a client and convert it to the specific type
        re, lastError := client.doSomething()
        err = lastError //record the last error
} 
pool.Put(cli, err) // put back the client with the error
```
