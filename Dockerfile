FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/robot-gitee-welcome
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o robot-gitee-welcome -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'" .

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

COPY  --chown=welcome --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-welcome/robot-gitee-welcome /opt/app/robot-gitee-welcome

RUN chmod 550 /opt/app/robot-gitee-welcome && \
    echo "umask 027" >> /home/welcome/.bashrc && \
    echo 'set +o history' >> /home/welcome/.bashrc

ENTRYPOINT ["/opt/app/robot-gitee-welcome"]
