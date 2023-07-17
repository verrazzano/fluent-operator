# Build the manager binary
FROM ghcr.io/verrazzano/verrazzano-base:v1.0.0-20230626111229-3d949a4 as builder
RUN microdnf install -y \
    golang && \
    microdnf clean all

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/fluent-manager/main.go main.go
COPY apis apis/
COPY controllers controllers/
COPY pkg pkg/

RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM ghcr.io/verrazzano/verrazzano-base:v1.0.0-20230327155846-4653b27
RUN groupadd -r verrazzano \
        && useradd -u 1000 -g verrazzano verrazzano \
        && chown -R 1000:verrazzano /home/verrazzano
WORKDIR /
COPY --from=builder --chown=verrazzano:verrazzano /workspace/manager .

RUN mkdir -p /license
COPY LICENSE README.md THIRD_PARTY_LICENSES.txt /license/

USER verrazzano

ENTRYPOINT ["/manager"]

