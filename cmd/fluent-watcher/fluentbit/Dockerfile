ARG FLUENTBIT_BASE
FROM ghcr.io/verrazzano/verrazzano-base-ol9:v1.0.0-20230627070315-744eefa as buildergo
RUN microdnf update \
    && microdnf install -y \
    golang && \
    microdnf clean all

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

