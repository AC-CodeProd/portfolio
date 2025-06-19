
FROM node:22.16.0-alpine AS node
FROM alpine:3.22.0 AS base
ARG TZ
ARG UID
ARG GID
ENV TZ=${TZ} UID=${UID} GID=${GID}

FROM base AS development
RUN apk --no-cache add bash python3 gcc make g++ zlib-dev openssl git curl tzdata linux-headers\
  && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone \
  && curl -fsSL "https://go.dev/dl/go1.24.4.linux-amd64.tar.gz" | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}" \
  GOCACHE=/portfolio/server/tmp/.cache \
  GOLANGCI_LINT_CACHE=/portfolio/server/tmp/.cache

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
RUN go install github.com/air-verse/air@latest 
WORKDIR /portfolio
COPY . .

WORKDIR /portfolio/server

RUN go mod download
RUN go mod tidy

WORKDIR /portfolio/ui

# Node
ENV PNPM_HOME="/portfolio/ui/.pnpm"
ENV PATH="$PNPM_HOME:$PATH"
COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules

# COPY --from=node /usr/lib /usr/lib
# COPY --from=node /usr/local/lib /usr/local/lib
# COPY --from=node /usr/local/include /usr/local/include
# COPY --from=node /usr/local/bin /usr/local/bin

RUN node -v
RUN npm install -g pnpm@latest-10
RUN pnpm config set store-dir /portfolio/ui/.pnpm-store
RUN pnpm install --store-dir /portfolio/ui/.pnpm-store


WORKDIR /portfolio


EXPOSE 8080 5173
CMD ["./start-dev.sh"]