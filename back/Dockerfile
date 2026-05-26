# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.26.3-alpine AS build
RUN apk --no-cache add --upgrade ca-certificates && update-ca-certificates
RUN adduser -g "" -s /sbin/nologin -D -H -u 1001 appuser

# create a working directory inside the image
WORKDIR /usr/src/app

# copy Go modules and dependencies to image
COPY ./back/go.mod ./
COPY ./back/go.sum ./

# download Go modules and dependencies
RUN go mod download && go mod verify

# copy directory files i.e all files ending with .go
COPY ./back/ ./

# generate swagger docs and remove LeftDelim/RightDelim lines
RUN go tool swag init && \
    sed -i -e '/LeftDelim/d' -e '/RightDelim/d' docs/docs.go

# compile application
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/src/app/slotfinder

# copy JWT keys (they will be available in the build context)
COPY ./back/config/jwt/ /usr/src/app/config/jwt/

##
## STEP 2 - DEPLOY
##
FROM scratch

COPY --from=build /usr/src/app/slotfinder /slotfinder
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/
COPY --from=build /usr/src/app/config/jwt/ /config/jwt/

USER appuser

EXPOSE 3000

ENTRYPOINT ["/slotfinder"]
