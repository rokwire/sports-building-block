FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache --update make git

RUN mkdir /sp-app
WORKDIR /sp-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.21.3

#we need timezone database + certificates
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /sp-app/bin/sport /

COPY --from=builder /sp-app/driven/storage/sport-definitions.json /driven/storage/sport-definitions.json

COPY --from=builder /sp-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv
COPY --from=builder /sp-app/vendor/github.com/rokwire/core-auth-library-go/v2/authorization/authorization_model_string.conf /sp-app/vendor/github.com/rokwire/core-auth-library-go/v2/authorization/authorization_model_string.conf

ENTRYPOINT ["/sport"]
