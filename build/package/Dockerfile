ARG ARCH=
FROM ${ARCH}openjdk:16-slim-buster

VOLUME /data
EXPOSE 8080 25565

COPY minecraft-server /usr/local/bin/qumine-server

RUN apt update \
&& apt install -y unzip wget \
&& groupadd -g 4242 -r qumine \
&& useradd -g 4242 -m -N -r -u 4242 qumine \
# Tempfix for fabric installer
# && useradd -g 4242 -M -N -r -u 4242 qumine \
&& mkdir /data \
&& chown -R qumine:qumine /data \
&& chmod -R 777 /data \
&& echo "qumine-server c" > /usr/local/bin/console \
&& chmod -R 755 /usr/local/bin/qumine-server /usr/local/bin/console

USER qumine
WORKDIR /data/
ENTRYPOINT [ "qumine-server", "server"]