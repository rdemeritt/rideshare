# Container image that runs your code
FROM registry.access.redhat.com/ubi9/ubi:latest

RUN yum install -y procps git make && \
    yum clean all

# install brew
RUN useradd -m builder
# RUN git clone https://github.com/Homebrew/brew /home/linuxbrew/.linuxbrew
# RUN eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" && \
#     brew update --force --quiet

RUN curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh > install.sh
RUN chmod +x install.sh && \
    chown builder:builder install.sh
USER builder
RUN ./install.sh
RUN echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.bashrc

RUN eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" && \
    brew install go emscripten protobuf && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
