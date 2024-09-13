FROM ubuntu:20.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install \
    gcc libssl-dev pkg-config protobuf-compiler \
    ca-certificates libssl1.1 vim strace lsof curl jq && \
    rm -rf /var/cache/apt /var/lib/apt/lists/*

RUN DEBIAN_FRONTEND=noninteractive apt-get install git libpq-dev python-dev python3-pip
RUN python -m pip install dbt-core dbt-postgres

ADD /honey-tracker /app/honey-tracker

ENV PATH "/app:$PATH"

ENTRYPOINT ["/app/honey-tracker"]