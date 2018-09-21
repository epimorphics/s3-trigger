FROM bitnami/minideb:jessie

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
ADD s3-controller /s3-controller

ENTRYPOINT ["/s3-controller"]
