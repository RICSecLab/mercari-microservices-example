FROM gcr.io/oss-fuzz-base/base-builder-go:v1
ENV GOPRIVATE "github.com/RICSecLab/mercari-microservices-example"
RUN apt-get update && \
    apt-get install -y make autoconf automake libtool \
    ca-certificates curl gnupg lsb-release
COPY . $SRC/mercari-microservices-example
WORKDIR mercari-microservices-example
COPY .clusterfuzzlite/build.sh $SRC/

