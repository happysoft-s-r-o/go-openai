# syntax=docker/dockerfile:1
ARG GOLANG_IMAGE=golang:1.22

# ========================================================================================================================
# Dev `GO` container image
# ========================================================================================================================
FROM $GOLANG_IMAGE

ENV PATH $PATH:/go/bin
ENV GOPATH /go
ENV GO111MODULE on
ENV CHOKIDAR_USEPOLLING true
ENV CGO_CFLAGS="-g -O2 -Wno-return-local-addr"

VOLUME [ "/commandhistory" ]

RUN \
    SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
    && echo $SNIPPET >> "/root/.bashrc" \
    && apt-get update \
    && apt-get install -y \
    git \
    make \
    vim \
    gcc \
    libc-dev \
    libssl-dev \
    libghc-zlib-dev \
    unzip 

RUN \
    go install github.com/go-delve/delve/cmd/dlv@v1.21.0 \
    && go install github.com/rogpeppe/godef@v1.1.2 \
    && go install golang.org/x/tools/gopls@latest \
    && go install golang.org/x/tools/cmd/goimports@v0.12.0 \
    && go install github.com/cweill/gotests/gotests@v1.6.0 \
    && go install honnef.co/go/tools/cmd/staticcheck@v0.4.5 \
    && go install github.com/oligot/go-mod-upgrade@v0.9.1 \
    && go install github.com/swaggo/swag/cmd/swag@v1.8.12 \
    && go install github.com/securego/gosec/v2/cmd/gosec@v2.17.0 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1 \
    && go install github.com/ramya-rao-a/go-outline@v0.0.0-20210608161538-9736a4bde949 \
    && echo "done"

ARG WORKDIR
WORKDIR ${WORKDIR}

CMD ["bash", "-c", "while true; do sleep 60; done"]
