# Contributing to this Project
**Contributions are welcome :)**
**If you don't now how to code, simply open an issue.**

Checkout the [TODO list](https://github.com/OpenTransports/Helsinki/projects/1) to see the roadmap.

- Pick a task
- Open an issue so I now you're working on it
- Make a pull request on the develop branch


# Development environment setup
Install [golang](https://golang.org/doc/install) then:
Install govendor: `go get https://github.com/kardianos/govendor`


```shell
git clone https://github.com/OpenTransports/Helsinki  # Clone the repo
go get -u github.com/kardianos/govendor               # Install govendor
govendor sync                                         # Install dependencies
go build                                              # Build the server
./Paris                                               # Launch the server on port 8080
```
WARNING: First launch will be long because it's fetching all needed resources

# Architecture
The system fetch the data from [digitransit.fi](https://digitransit.fi/en)
See: the [GraphQL API](http://dev.hsl.fi/graphql/console)
