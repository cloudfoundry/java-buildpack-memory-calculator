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

RUN curl --location 'https://api.bintray.com/content/jfrog/jfrog-cli-go/$latest/jfrog-cli-linux-amd64/jfrog?bt_package=jfrog-cli-linux-amd64' > /usr/local/bin/jfrog && \
    chmod +x /usr/local/bin/jfrog
