FROM ubuntu:24.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
      apt-get -y install -y ca-certificates libssl3

ADD /honey-tracker /app/honey-tracker

ENV PATH "/app:$PATH"

ENTRYPOINT ["/app/honey-tracker"]