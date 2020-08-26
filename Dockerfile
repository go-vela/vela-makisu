#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:certs .    ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-makisu:local .    ##
##########################################################

FROM gcr.io/uber-container-tools/makisu:v0.3.1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-makisu /bin/vela-makisu

ENTRYPOINT [ "/bin/vela-makisu" ] 