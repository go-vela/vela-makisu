#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:binary .   ##
#########################################################################

FROM gcr.io/uber-container-tools/makisu:v0.3.1 as makisu

#########################################################################
##    docker build --no-cache --target certs -t vela-makisu:certs .    ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-makisu:local .    ##
##########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=certs /makisu-internal/makisu /bin/makisu


COPY release/vela-makisu /bin/vela-makisu

ENTRYPOINT [ "/bin/vela-makisu" ] 