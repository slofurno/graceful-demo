FROM scratch

COPY deploy-demo /deploy-demo

ENTRYPOINT ["/deploy-demo"]
