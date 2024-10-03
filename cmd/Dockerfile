FROM ubuntu:latest AS ubuntu
# Get latest CA certificate bundle from Ubuntu, to avoid "x509: certificate signed by unknown
# authority"
RUN apt-get update && apt-get install -y ca-certificates

RUN useradd -u 999 user

# FROM scratch

# COPY --from=ubuntu /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD ./build/ /opt/webhooks-adaptor/
ADD ./data/ /opt/webhooks-adaptor/data/
WORKDIR /opt/webhooks-adaptor/

USER 999
CMD ["./webhooks-adaptor"]
