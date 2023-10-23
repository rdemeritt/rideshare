# Container image that runs your code
FROM registry.access.redhat.com/ubi9/ubi:latest

# install OS packages
RUN yum install -y procps git make sudo && \
    yum clean all

# add builder user
RUN useradd -m builder --uid 1001
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
