# Build stage - Create static binary
FROM golang:1.8.1
WORKDIR /go/src/github.com/phoreproject/openbazaar-go
COPY . .
RUN go build --ldflags '-extldflags "-static"' -o /opt/openbazaard .

# Run stage - Import static binary, expose ports, set up volume, and run server
FROM scratch
EXPOSE 5001 5002 10005
VOLUME /var/lib/openbazaar
COPY --from=0 /opt/openbazaard /opt/openbazaard
COPY --from=0 /etc/ssl/certs/ /etc/ssl/certs/
ENTRYPOINT ["/opt/openbazaard"]
CMD ["start", "-d", "/var/lib/openbazaar"]
