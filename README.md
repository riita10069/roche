# ![roche](./roche_2nd.png)

It is still under development, so please do not use it. We plan to release v.1.0.0 in the summer.

roche is a web framework optimized for microservice architecture using go.
It prioritizes schema-first and quick development over lightweight.
It is currently in the development concept stage and should not be used in production.

It is inspired by Ruby on Rails, [https://github.com/izumin5210/grapi][grapi] and protoeasy.
And roche relies on [https://github.com/izumin5210/grapi][grapi] heavily.

## Roche is named after the Roche limit and Roche radius.

The Roche limit is the critical satellite orbital radius at which the satellite will break apart due to tidal destruction if the satellite is approximated to be a self-gravitating fluid, ignoring its material strength.
The Roche radius is the size of the region where the mutual gravity of two celestial bodies orbiting around the central body exceeds the gravity of the central body.

In a microservices architecture, each service pulls together as if it has universal gravitation. However, if they are too close together (i.e., too dense), they will break, hence the name.

## Getting Started
### protobuf
You have to install [https://github.com/protocolbuffers/protobuf][protobuf]

### Framework

`go get github.com/riita10069/roche`

### CLI tool

`go install github.com/riita10069/roche/cmd/roche@latest`

## Getting Started
Requires go version 1.16 or higher

`roche init .`
To Initialize project like create-react-app

`roche toml`
Generate roche.toml which is needed roche command.

`roche g scaffold-service NAME`
Generate .proto and server

`roche scaffold all NAME`
Generate CRUD template for NAME service
You can also generate a single file.
By using (domain, model, repo, migrate) instead of all.

`roche manifest app APPNAME`
Generate k8s manifest template, APPNAME is app name.

`roche server`
Run api server

## if you want to use roche server
Currently, this feature is not supported. Sorry.
```go
func Run() error {
+     defer roche.Close()
+
      s := grapiserver.New(
+               rocheserver.WithDefaultLogger(),
-               grapiserver.WithDefaultLogger(),
                grapiserver.WithServers(
                // TODO: Add server implementations
  	),
        )
        return s.Serve()
}
```


## core maintainers

- <a href="https://github.com/riita10069">@riita10069</a>
- <a href="https://github.com/hourglasshoro">@hourglasshoro</a>
- <a href="https://github.com/ffjlabo">@ffjlabo</a>

[protobuf]: https://github.com/protocolbuffers/protobuf

[grapi]: https://github.com/izumin5210/grapi
