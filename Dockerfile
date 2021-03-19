#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:binary .   ##
#########################################################################

FROM gcr.io/uber-container-tools/makisu:v0.4.2 as makisu

#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:certs .    ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:conf .    ##
#########################################################################

FROM alpine as conf

RUN mkdir -p /makisu/registry/ && touch /makisu/registry/config.json

##########################################################
##    docker build --no-cache -t vela-makisu:local .    ##
##########################################################

FROM scratch

COPY --from=makisu /makisu-internal/makisu /bin/makisu
COPY --from=makisu /makisu-internal/certs/cacerts.pem /makisu-internal/certs/cacerts.pem
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=conf /makisu/registry/config.json /makisu/registry/config.json

COPY release/vela-makisu /bin/vela-makisu

ENTRYPOINT [ "/bin/vela-makisu" ] 
