# Setup

## K3S
this was made with k3s

install/configure k3s and make sure the firewall is disabled

## Helm Pre-Reqs


### Registry

docker registry must be setup before building scripts

```
helm install registry helm/registry
```

### Postgresql

be sure to set the password or postgress will not come up

```
helm install postgress helm/postgresql/ --set-string "postgresConfig.POSTGRES_PASSWORD=mypassword"
```

### Monitoring Stack

using opentelemetry collector to collect and distribute telemetry

#### Otel Collector

```
helm install collector helm/collector
```


#### Jaeger

```
helm install jaeger helm/jaeger
```

#### Prometheus/Grafana stack

```
helm install monitoring helm/monitoring
```

configure jaeger backend

connections -> Add new connection -> jaeger

set url to `http://jaeger-query.default.svc.cluster.local:16686`

# My-Service

simple rest widget

## Podman and Buildah

This was built with the [buildah/podman toolchain](https://developers.redhat.com/blog/2019/02/21/podman-and-buildah-for-docker-users)

it is supposed to be a 1:1 feature paridy with docker, so you can do

```
export BUILD_EXEC='docker'
```

to use docker to build the containers when running the local scripts


## Build Containers

the containers are build in go, so it should just be a matter of running

```
/build/build-containers.sh
```

the following are optional environmental variables

```
BUILD_EXEC=${BUILD_EXEC:-'buildah'}
OCI_REGISTRY=${OCI_REGISTRY:-'localhost:31000'}
BUILD_TAG=${BUILD_TAG:-'latest'}
```

## Configure Service

copy config.example.yaml to make a new values file and update password. this will be stored as a
secret

```
serviceConfig:
  database:
    password: "mypassword"
```

## Deploy

finally install service

```
helm install -f config.yaml my-service helm/my-service/
```

# Tests

the tests can be run from a docker compose file.

## Start Test Containers

```
podman-compose -f build/compose.yaml up -d
#check status is ok
podman-compose -f build/compose.yaml ps
```

## Run Tests

```
podman-compose -f build/compose.yaml exec test
go test ./...
exit
```

## Stop Test Containers

```
podman-compose -f build/compose.yaml down
```
