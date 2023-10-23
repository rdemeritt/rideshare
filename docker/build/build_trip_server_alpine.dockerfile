# Container image that runs your code
FROM alpine:3.18

# install OS packages
RUN apk add procps curl git make sudo bash

# add builder user
RUN adduser -D builder -u 1001
RUN echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# install homebrew
RUN curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh > homebrew_installer.sh
RUN chmod +x homebrew_installer.sh
RUN chown builder:builder homebrew_installer.sh
RUN sudo -u builder ./homebrew_installer.sh

# setup homebrew
USER builder
ENV PATH $PATH:/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin
RUN echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.bashrc

RUN brew install go
# RUN brew install go protobuf
# RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
