# Build the manager binary
ARG BASE_IMAGE=ghcr.io/verrazzano/golang:v1.20.12
FROM $BASE_IMAGE as builder

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

RUN groupadd -r verrazzano \
        && useradd -u 1000 -g verrazzano verrazzano \
        && chown -R 1000:verrazzano /home/verrazzano

RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
FROM ghcr.io/verrazzano/ol8-static:v0.0.1-20231102152128-e7afc807
WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=verrazzano:verrazzano /workspace/manager .

COPY LICENSE README.md THIRD_PARTY_LICENSES.txt /license/

USER verrazzano

ENTRYPOINT ["/manager"]
