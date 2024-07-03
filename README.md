# Setup

## Helm Pre-Reqs

### Prometheus/Grafana stack

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install monitoring prometheus-community/kube-prometheus-stack
```

## Registry

docker registry must be setup before building scripts

```
helm install registry helm/registry
```

## Postgresql

todo

## Podman and Buildah

This was built with the [buildah/podman toolchain](https://developers.redhat.com/blog/2019/02/21/podman-and-buildah-for-docker-users)

it is supposed to be a 1:1 feature paridy with docker, so you can do

```
export BUILD_EXEC='docker'
```

to use docker to build the containers when running the local scripts


# Build

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
