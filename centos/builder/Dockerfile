FROM centos:7

COPY ./ /tmp/source

RUN set -ex ;\
    cp -rfav /tmp/source/builder/assets/* / ;\
    mkdir -p /srv/clementine-centos/{cache,build-repo,repo} ;\
    cp -rfav /tmp/source/definition /srv/clementine-centos/ ;\
    yum update -y ;\
    yum install -y \
      git \
      rpm-ostree

WORKDIR /srv/clementine-centos
