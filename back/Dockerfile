# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.25.2-alpine AS build
RUN apk --no-cache add --upgrade ca-certificates && update-ca-certificates
RUN adduser -g "" -s /sbin/nologin -D -H -u 1001 appuser

# create a working directory inside the image
WORKDIR /usr/src/app

# copy Go modules and dependencies to image
COPY ./go.mod ./
COPY ./go.sum ./

# download Go modules and dependencies
RUN go mod download && go mod verify

# copy directory files i.e all files ending with .go
COPY ./ ./

# generate swagger docs and remove LeftDelim/RightDelim lines
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init && \
    sed -i -e '/LeftDelim/d' -e '/RightDelim/d' docs/docs.go

# compile application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/src/app

##
## STEP 2 - DEPLOY
##
FROM scratch

COPY --from=build /usr/src/app /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/

USER appuser

EXPOSE 3000

ENTRYPOINT ["/app"]
