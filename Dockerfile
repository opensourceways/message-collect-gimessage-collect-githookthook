FROM golang:latest as BUILDER

MAINTAINER shishupei

ARG USER
ARG PASS
RUN echo "machine github.com login $USER password $PASS" >/root/.netrc

# build binary
WORKDIR /go/src/github.com/opensourceways/message-collect-githook
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o message-collect-githook -buildmode=pie --ldflags "-s -extldflags '-Wl,-z,now'" .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 1000 welcome && \
    useradd -u 1000 -g welcome -s /sbin/nologin -m welcome && \
    echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd && \
    mkdir /home/welcome -p && \
    chmod 700 /home/welcome && \
    chown welcome:welcome /home/welcome && \
    echo 'set +o history' >> /root/.bashrc && \
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs && \
    rm -rf /tmp/*

USER welcome

WORKDIR /opt/app

COPY  --chown=welcome --from=BUILDER /go/src/github.com/opensourceways/message-collect-githook/message-collect-githook /opt/app/message-collect-githook

RUN chmod 550 /opt/app/message-collect-githook && \
    echo "umask 027" >> /home/welcome/.bashrc && \
    echo 'set +o history' >> /home/welcome/.bashrc

ENTRYPOINT ["/opt/app/message-collect-githook"]