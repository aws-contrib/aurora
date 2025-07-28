# syntax=docker/dockerfile:1.4

FROM --platform=$BUILDPLATFORM golang:1.24 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /root/workspace

RUN \
  --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS 
ENV GOARCH=$TARGETARCH

# build everything
# hadolint ignore=DL4006
RUN \
  --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=.,target=/root/workspace \
  go build -o /root/bin/aurora -mod readonly -v ./cmd/aurora

# Final minimal image
FROM scratch
COPY --from=builder /root/bin/aurora /bin/aurora

ENV PATH=/bin

ENTRYPOINT ["aurora"]
