FROM ubuntu:24.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install ca-certificates libssl3 \

ADD /honey-tracker /app/honey-tracker
ADD /dbt/hivemapper /app/hivemapper

ENV PATH "/app:$PATH"

ENTRYPOINT ["/app/honey-tracker"]