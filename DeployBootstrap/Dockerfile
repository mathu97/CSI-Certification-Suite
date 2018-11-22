FROM golang:latest
COPY ./e2e-bootstrap /e2e-bootstrap
ENTRYPOINT ["/e2e-bootstrap", "-endpoint", "/plugin/csi.sock", "-network", "unix"] 
