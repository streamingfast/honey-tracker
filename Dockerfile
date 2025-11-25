FROM ubuntu:24.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install \
    gcc libssl-dev pkg-config protobuf-compiler \
    ca-certificates libssl3 vim strace lsof curl jq git && \

ADD /honey-tracker /app/honey-tracker
ADD /dbt/hivemapper /app/hivemapper

ENV PATH "/app:$PATH"

ENTRYPOINT ["/app/honey-tracker"]