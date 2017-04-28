# The Procurement Listener Service #


The Procurement Listener Service is a reference implementation of the 
Procurement Listener API, which is used to listen to procurement related
events coming from an online store that sells Cloud-enabled software. 
Service providers can implement the API, similar to the implementation
here to participate in the procurement flow of the upstream store.

<p>This project also contains a basic conformance suite for testing the
compliance of implementations.


## Getting Started

This section outlines the steps needed to get started with this code.

- [Install and Setup Go (Prerequisite)](#install-and-setup-go)
- [Clone the Git Repository](#clone-the-git-repository)
- [Build the Service](#build-the-service)
- [Run the Listener Service Locally](#run-the-listener-service-locally)
- [Run Conformance Tests](#run-conformance-tests)

### Install and Setup Go

This project is written in the [Go](http://golang.org) programming language.
To build it, you'll need a Go development environment. If you haven't set up a Go development
environment, please follow [these instructions](http://golang.org/doc/code.html)
to install the Go tools.

Set up your GOPATH and add a path entry for Go binaries to your PATH. Typically
added to your ~/.profile:

```shell
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
```

### Clone the Git Repository

The commands below require that you have $GOPATH set ([$GOPATH
docs](https://golang.org/doc/code.html#GOPATH)).

```shell
cd $GOPATH/src
git clone https://github.com/google/procurement-listener-service.git
mv procurement-listener-service/ procurementlistenerservice/
```

### Build the Service

The service has a few dependencies that must also be cloned from GitHub.

``` shell
cd $GOPATH/src
mkdir -p github.com/gorilla
mkdir -p github.com/xeipuuv

cd $GOPATH/src/github.com/gorilla
git clone https://github.com/gorilla/mux.git
git clone https://github.com/gorilla/handlers.git

cd $GOPATH/src/github.com/xeipuuv
git clone https://github.com/xeipuuv/gojsonpointer.git
git clone https://github.com/xeipuuv/gojsonreference.git
git clone https://github.com/xeipuuv/gojsonschema.git
```

You should now be able to build the Procurement Listener service.

```shell
go build procurementlistenerservice
```

### Run the Listener Service Locally
```shell
cd $GOPATH/src/procurementlistenerservice
. scripts/runlocal.sh
```

### Run Conformance Tests
```shell
cd $GOPATH/src/procurementlistenerservice/conformance
go test
```

