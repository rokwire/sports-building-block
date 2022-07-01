FROM golang:1.16-buster as builder

ENV CGO_ENABLED=0

RUN mkdir /sp-app
WORKDIR /sp-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.13

#we need timezone database
RUN apk --no-cache add tzdata

COPY --from=builder /sp-app/bin/sport /

COPY --from=builder /sp-app/driven/storage/sport-definitions.json /driven/storage/sport-definitions.json

COPY --from=builder /sp-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv
COPY --from=builder /sp-app/vendor/github.com/rokwire/core-auth-library-go/v2/authorization/authorization_model_string.conf /sp-app/vendor/github.com/rokwire/core-auth-library-go/v2/authorization/authorization_model_string.conf

ENTRYPOINT ["/sport"]
