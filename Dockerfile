
FROM python:alpine

LABEL maintainer="${QL_MAINTAINER}"
ARG TARGETARCH

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/local/go/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/.local/share/pnpm/global/5/node_modules \
    LANG=zh_CN.UTF-8 \
    SHELL=/bin/bash \
    PS1="\u@\h:\w \$ " \
    QL_DIR=/ql \
    QL_BRANCH=${QL_BRANCH}

WORKDIR ${QL_DIR}

RUN set -x \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update -f \
    && apk upgrade \
    && apk --no-cache add -f bash \
                             coreutils \
                             moreutils \
                             git \
                             curl \
                             wget \
                             tzdata \
                             perl \
                             openssl \
                             nodejs \
                             jq \
                             openssh \
                             npm \
    && rm -rf /var/cache/apk/* \
    && apk update \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && git config --global user.email "qinglong@@users.noreply.github.com" \
    && git config --global user.name "qinglong" \
    && git config --global http.postBuffer 524288000 \
    && npm install -g yarn \
    && rm -rf /root/.cache \
    && rm -rf /root/.npm \
    && mkdir -p ${QL_DIR}/data \
    && cd ${QL_DIR} \
    && git clone https://github.com/whyour/qinglong-static.git ./static/

COPY ./dist/docker_linux_$TARGETARCH*/qinglong-go ${QL_DIR}/ql

RUN  chmod -R 777 /ql/ql

EXPOSE 5700

VOLUME ${QL_DIR}/data


CMD cd /ql && chmod -R 777 ./ql && ./ql