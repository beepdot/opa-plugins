# opa-plugins

> OPA version - 0.34.2
#### Build opa binary
```
go mod tidy
go mod vendor
go build -o opa cmd/main.go
```


#### Build Docker Image
- To build the docker image run
```
docker build -t docker_user_name/opa:0.34.2-envoy --build-arg BASE=gcr.io/distroless/cc -f Dockerfile .
```