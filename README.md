# roar
Roar CI

```bash
$ go run main.go -port 4000
$ curl -i 'http://localhost:4000/api?cmd=./do_backend.sh&cmd=./do_frontend.sh'
```

You can pass any number of `cmd`.
All of them will be running at same time, so if you need something sync, just use one `cmd`.
