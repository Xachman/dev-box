FROM debian:jessie

RUN apt-get update && apt-get install -y curl && \
    curl -fsSL get.docker.com -o /opt/get-docker.sh && sh /opt/get-docker.sh && \
    rm -rf /var/lib/apt/lists/*

COPY ./src/main /app/main

CMD /app/main