FROM golang:1.23-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	go build -trimpath -ldflags="-s -w" -o /out/auto-workflow .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /data/app
COPY --from=builder /out/auto-workflow /data/app/auto-workflow
COPY config.yaml /data/app/config.yaml

ENV AUTO_WORKFLOW_CONFIG=/data/app/config.yaml
EXPOSE 8080
ENTRYPOINT ["/data/app/auto-workflow"]
