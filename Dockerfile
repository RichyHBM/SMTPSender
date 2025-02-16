# Multistage Docker File

ARG BUILDER_IMAGE=golang:alpine

############################
# STEP 1 build binary
############################

FROM ${BUILDER_IMAGE} AS builder

# Install dependencies for final container.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=1000

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/mypackage/myapp/

# use modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download && go mod verify

COPY . .

# Build as static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/server .

############################
# STEP 2 build runtime image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy static executable
COPY --from=builder /go/bin/server /go/bin/server

# Use an unprivileged user.
USER appuser:appuser

# Run the server binary.
ENTRYPOINT ["/go/bin/server"]