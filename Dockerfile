FROM ghcr.io/battleofthebots/botb-base-image:latest AS builder
RUN apt install -y golang
WORKDIR /opt
COPY WOPR.go . 
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath WOPR.go

FROM builder

# Update this with youe challenge name if you are pushing to a docker registry
ARG NAME=wopr
LABEL org.opencontainers.image.title=$NAME org.opencontainers.image.description=$NAME org.opencontainers.image.url=https://github.com/battleofthebots/$NAME org.opencontainers.image.source=https://github.com/battleofthebots/$NAME org.opencontainers.image.version=main

COPY --from=builder /opt/WOPR .
EXPOSE 4000
USER user
CMD ./WOPR server