## Build Instructions

The upstream tag this release is branched from is `v2.5.0`

### Create Environment Variables

```
export DOCKER_REPO=<Docker Repository>
export DOCKER_NAMESPACE=<Docker Namespace>
export DOCKER_TAG=<Docker Tag>
export FO_IMG=${DOCKER_REPO}/${DOCKER_NAMESPACE}/fluent-operator:${DOCKER_TAG}
```

### Build and Push Images

From the root of the repo, run the following command to build the operator image
```
docker build -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}
```

Once the build completes successfully, push the image:
```
docker push ${FO_IMG}
```
