# syntax=docker/dockerfile:1
FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm

# We need to be root for a while
USER root

# Do not delete the downloaded packages
RUN \
  rm -f /etc/apt/apt.conf.d/docker-clean && \
  echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache

# Install OS applications
RUN \
  --mount=target=/var/cache/apt,type=cache,sharing=locked \
  --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
  apt-get update && apt-get -y install --no-install-recommends ca-certificates postgresql-client
