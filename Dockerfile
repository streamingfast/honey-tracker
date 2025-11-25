FROM ubuntu:24.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install \
    gcc libssl-dev pkg-config protobuf-compiler \
    ca-certificates libssl1.1 vim strace lsof curl jq git python3-pip && \
    DEBIAN_FRONTEND=noninteractive apt-get remove python-cffi && \
    DEBIAN_FRONTEND=noninteractive pip install --upgrade cffi && \
    DEBIAN_FRONTEND=noninteractive pip install cryptography~=3.4 && \
    rm -rf /var/cache/apt /var/lib/apt/lists/*

RUN pip install dbt-core dbt-postgres

ADD /honey-tracker /app/honey-tracker
ADD /dbt/hivemapper /app/hivemapper

ENV PATH "/app:$PATH"

ENTRYPOINT ["/app/honey-tracker"]