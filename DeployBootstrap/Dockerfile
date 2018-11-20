FROM golang:latest
COPY e2e-bootstrap /e2e-bootstrap
ENTRYPOINT ["e2e-boostrap"] 
CMD ["--endpoint /plugin/csi.sock --network unix"]
