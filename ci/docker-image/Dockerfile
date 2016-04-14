FROM golang:alpine

RUN apk add --update \
    bash \
    coreutils \
    curl \
    git  \
    jq \
    py-pip \
 && rm -rf /var/cache/apk/*

RUN pip install awscli --ignore-installed six
