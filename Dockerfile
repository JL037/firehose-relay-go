# SCAFFOLD: skeleton multi-stage build. TODO: revisit base image pin and add a
# non-root user once there's a binary worth shipping.
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod ./
# COPY go.sum ./   # once there are dependencies
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/relay ./cmd/relay

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/relay /relay
EXPOSE 8080
ENTRYPOINT ["/relay"]
