FROM ubuntu:bionic

COPY ./ /tmp/source

RUN set -ex ;\
    cp -rfav /tmp/source/builder/assets/* / ;\
    mkdir -p /srv/clementine-ubuntu/repo ;\
    cp -rfav /tmp/source/definition /srv/clementine-ubuntu/ ;\
    apt-get update ;\
    apt-get dist-upgrade -y ;\
    DEBIAN_FRONTEND=noninteractive apt-get install  --no-install-recommends -y \
        rsync \
        ostree \
        python3-apt \
        attr \
        grub2 \
        git \
        multistrap \
        uuid-runtime

WORKDIR /srv/clementine-ubuntu
