# Bop Bag

[![Build, Test and coverage](https://github.com/balchua/bop-bag/actions/workflows/ci.yml/badge.svg)](https://github.com/balchua/bop-bag/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/balchua/bop-bag/branch/main/graph/badge.svg?token=S9UITN8L61)](https://codecov.io/gh/balchua/bop-bag)


The application is inspired by those toys used by kids as a punching bags.  These toys do not tip over.

When one tries to hold them down, the toy immediately spring back up the moment one releases their hands on the toys.

What is a bop bag?

You can find several bop bag toys in [Amazon](https://www.amazon.com/bop-bag/s?k=bop+bag)

A application that can tolerate node failures.

This is a demonstration of a very low ops and highly available application.

## Goal
* Demonstrate an application which can self heal.
* Embed raft enabled highly available data store.
* Simple enough to be understood.
* Simple to join a node to the cluster.
* Must be able to demonstrate the auto cluster realignment when a node goes down.

## Generate certificates

There is a sample openssl configuration template that you can use [here](default-certs/csr-dqlite.conf.template)

```shell
openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes -keyout ${HOME}/bopbag/certs/cluster.key -out ${HOME}/bopbag/certs/cluster.crt -subj "/CN=bopbag" -config csr-dqlite.conf -extensions v3_ext
```

## Implementation

We will be using [dqlite](https://dqlite.io/) -  _Dqlite is a fast, embedded, persistent SQL database with Raft consensus that is perfect for fault-tolerant IoT and Edge devices._ to keep our data.

Expose REST endpoint to create, update, read and delete records

### Table structure

The implementation will be a simple go based called `TaskRepository`.

`TASK_MASTER` Table structure:

| Columns | Type | Description |
|---------|------|-------------|
| ID | INT | The primary key of the task|
| TITLE | VARCHAR(50) | The title of the task to do |
| DETAILS | VARCHAR(1000) | Details of the task |
| CREATED_DATE | DATE | The date the task is created |

### REST Endpoints

- [X] GET all tasks
  
  * Endpoint: `/api/v1/tasks/`
  * Method : `GET`
 
- [X] GET a task
  * Endpoint: `/api/v1/task/{id}`
  * Method: `GET`
 
- [X] Insert a task
  * Endpoint `/api/v1/task`
  * Method: `POST`

- [ ] Update a task
  * Endpoint: `/api/v1/task/{id}`
  * Method: `PUT`

- [ ] Delete a task
  * Endpoint: `/api/v1/task/{id}`
  * Method: `DELETE`

- [X] Shows the cluster information
  * Endpoint: `/api/v1/clusterInfo`
  * Method: `GET`

Sample output of Cluster Info


### Joining the cluster

Joining or forming a cluster must be easy.

### Leaving the cluster

TODO.

## Build

Assume that you have local registry running at `localhost:32000`

### Building with docker

Build the base image

```shell
docker build -f Dockerfile.base -t localhost:32000/dqlite-base:1.9.0 .
```

Building the app 

```
docker build -t localhost:32000/bopbag .
```

### Building in your host

Pre-requisite:
* Install `dqlite`, instructions here https://github.com/canonical/dqlite#install.  Use the `master` version.  

```shell
sudo add-apt-repository ppa:dqlite/master 
sudo apt update
sudo apt-get -y install clang lcov libsqlite3-dev libraft-dev 
```

Finally build the application

```shell
export CGO_LDFLAGS_ALLOW="-Wl,-z,now"
go build .
```

Run unit test
```
go test -p=1 --coverpkg=./... -coverprofile=cover.out ./...
```

### Starting the nodes on local machine

First node
```
./bopbag serve --db /tmp/dbPath --certs default-certs --dbAddress norse:9000
```

Second node

```
./bopbag serve --db /tmp/dbPath2 --certs default-certs --port 8081 --dbAddress norse:9001 --join norse:9000
```

Third node

```
./bopbag serve --db /tmp/dbPath3/ --certs default-certs --port 8082 --dbAddress  norse:9003 --join norse:9000
```

Show cluster information `/api/v1/clusterInfo`

```yaml
[
  {
    "ID": 3297041220608546300,
    "Address": "norse:9000",
    "Role": 0
  },
  {
    "ID": 7997991560008497000,
    "Address": "norse:9001",
    "Role": 0
  },
  {
    "ID": 11590821130369819000,
    "Address": "norse:9003",
    "Role": 0
  }
]
```

### Deploying to Kubernetes

Install MicroK8s on Ubuntu

```shell

sudo snap install microk8s --channel 1.22/stable --classic
```
Once the cluster is up and running, enable the following addons

```shell
microk8s enable dns registry storage
```

Push the docker imge to the local registry

```shell
docker push localhost:32000/bopbag
```

Finally apply the manifest [here](manifests/bopbag.yaml)

```shell
# The sample manifest is deployed to the bopbag namespace
microk8s kubectl create ns bopbag
# deploy the application
microk8s kubectl apply -f manifest/bopbag.yaml
```

Check that the application is running

```shell
microk8s kubectl -n bopbag get pods -o wide

kubectl -n bopbag get pods -o wide
NAME       READY   STATUS        RESTARTS   AGE     IP            NODE    NOMINATED NODE   READINESS GATES
bopbag-0   1/1     Running       0          11m     10.1.205.14   norse   <none>           <none>
```

Scale the number of replicas to 3

```
microk8s kubectl -n bopbag scale sts/bopbag --replicas=3
```

### Simulating faults

To simulate faults, simple scale the replicas to less than a majority for example `1`

You will no longer be able to access the endpoints like `/api/v1/tasks`

## Operations

### Adding entries

```
curl -d '{ "title": "My First Task", "details": "Here you go, this is what i should do", "createdDate": "2021-10-25"}' -H "Content-Type: application/json" -X POST http://localhost:32657/api/v1/task
```

### Get cluster information

```
curl -X GET http://localhost:32657/api/v1/clusterInfo
```

### Get all tasks

```
curl -X GET http://localhost:32657/api/v1/tasks
```
