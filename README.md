# Uncapsizable

A service that can sustain failures with low ops.


## Goal
* Demonstrate an application which can self heal.
* Embed raft enabled high availability data store.
* Simple enough to be understood.
* Simple to join a node to the cluster.
* Must be able to demonstrate the auto cluster realignment when a node goes down.


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
  
  * Endpoint: `/tasks/`
  * Method : `GET`
 
- [X] GET a task
  * Endpoint: `/task/{id}`
  * Method: `GET`
 
- [X] Insert a task
  * Endpoint `/task`
  * Method: `PUT`

- [ ] Update a task
  * Endpoint: `/task/{id}`
  * Method: `PUT`

- [ ] Delete a task
  * Endpoint: `/task/{id}`
  * Method: `DELETE`

### Joining the cluster

The joining process must be automatic.

## Build

Assume that you have local registry running at `localhost:32000`

### Building with docker

Build the base image

```shell
docker build -t localhost:32000/dqlite-base:1.0.0
```

Building the app 

```
docker build -t localhost:32000/uncapsizable:1.0.0 .
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

### Starting the nodes on local machine

First node
```
./uncapsizable serve --db /tmp/dbPath
```

Second node

```
./uncapsizable serve --db /tmp/dbPath2 --port 8081 --dbAddress localhost:9001 --join 0.0.0.0:9000
```

Third node

```
./uncapsizable serve --db /tmp/dbPath3/ --port 8082 --dbAddress  localhost:9003 --join 0.0.0.0:9000
```
