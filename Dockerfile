FROM public.ecr.aws/docker/library/golang:1.24-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache --update make git

RUN mkdir /sp-app
WORKDIR /sp-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM public.ecr.aws/docker/library/alpine:3.21.3

#we need timezone database + certificates
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /sp-app/bin/sport /

COPY --from=builder /sp-app/driven/storage/sport-definitions.json /driven/storage/sport-definitions.json

COPY --from=builder /sp-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf /sp-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf
COPY --from=builder /sp-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf /sp-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf

ENTRYPOINT ["/sport"]
