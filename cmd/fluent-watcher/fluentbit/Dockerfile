ARG FLUENTBIT_BASE
ARG BASE_IMAGE=ghcr.io/verrazzano/golang:v1.20.12
FROM $BASE_IMAGE as buildergo
RUN mkdir -p /fluent-bit
RUN mkdir -p /code
COPY . /code/
WORKDIR /code
RUN echo $(ls -al /code)
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /fluent-bit/fluent-bit /code/cmd/fluent-watcher/fluentbit/main.go

FROM $FLUENTBIT_BASE
LABEL Description="Fluent Bit docker image" Vendor="Verrazzano" Version="2.0.11"

COPY conf/fluent-bit.conf conf/parsers.conf /fluent-bit/etc/
COPY --from=buildergo /fluent-bit/fluent-bit /fluent-bit/bin/fluent-bit-watcher

# Entry point
ENTRYPOINT ["/fluent-bit/bin/fluent-bit-watcher"]

