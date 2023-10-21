# Container image that runs your code
FROM registry.access.redhat.com/ubi9/ubi:latest

# install OS packages
RUN yum install -y procps git make sudo && \
    yum clean all

# add builder user
RUN useradd -m builder --uid 1001
RUN echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# install homebrew
RUN curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh > install.sh
RUN chmod +x install.sh && \
    chown builder:builder install.sh
RUN sudo -u builder ./install.sh

# setup homebrew
USER builder
RUN echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.bashrc

# install build dependencies
RUN eval $(/home/linuxbrew/.linuxbrew/bin/brew shellenv) && \
    brew install go protobuf && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
