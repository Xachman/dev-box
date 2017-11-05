FROM debian:jessie

RUN apt-get update && apt-get install -y sshfs fuse unison && mkdir /mnt/droplet /unison

COPY start.sh /usr/local/bin/start.sh

RUN chmod +x /usr/local/bin/start.sh


CMD start.sh 