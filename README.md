## roche

It is still under development, so please do not use it. We plan to release v.1.0.0 in the summer.

roche is a web framework optimized for microservice architecture using go.
It prioritizes schema-first and quick development over lightweight.
It is currently in the development concept stage and should not be used in production.

It is inspired by Ruby on Rails, grapi and protoeasy.
And roche relies on grapi heavily.

## Roche is named after the Roche limit and Roche radius.

The Roche limit is the critical satellite orbital radius at which the satellite will break apart due to tidal destruction if the satellite is approximated to be a self-gravitating fluid, ignoring its material strength.
The Roche radius is the size of the region where the mutual gravity of two celestial bodies orbiting around the central body exceeds the gravity of the central body.

In a microservices architecture, each service pulls together as if it has universal gravitation. However, if they are too close together (i.e., too dense), they will break, hence the name.

## Getting Started

### Framework

`go get github.com/riita10069/roche`

### CLI tool

`go get github.com/riita10069/roche/cmd/roche`

## List of Command

`roche init .`

`roche g scaffold-service service_name`

`roche server`

`roche protoc`

## core maintainers

- <a href="https://github.com/riita10069">@riita10069</a>
- <a href="https://github.com/hourglasshoro">@hourglasshoro</a>
- <a href="https://github.com/ffjlabo">@ffjlabo</a>
