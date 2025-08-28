############################################
# Node stage
############################################
FROM node:22.16.0-alpine AS node

############################################
# Base Alpine + toolchain + Go
############################################
FROM alpine:3.22.0 AS base
ARG TZ=Europe/Paris
ARG UID=1000
ARG GID=1000
ARG USER=app
ARG GROUP=app

ENV TZ=${TZ} \
  UID=${UID} \
  GID=${GID} \
  USER=${USER} \
  GROUP=${GROUP}

RUN apk --no-cache add \
  bash python3 gcc g++ make zlib-dev openssl git curl tzdata linux-headers shadow ca-certificates \
  && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone \
  && curl -fsSL "https://go.dev/dl/go1.25.0.linux-amd64.tar.gz" | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:/usr/local/bin:${PATH}"

############################################
# Development stage
############################################
FROM base AS development
RUN set -eux; \
  addgroup -S -g "${GID}" "${GROUP}" 2>/dev/null || true; \
  if ! getent passwd "${UID}" >/dev/null; then \
    adduser -D -H -s /bin/sh -G "${GROUP}" -u "${UID}" "${USER}"; \
  fi; \
  mkdir -p "/home/${USER}/.cache/go-build" \
            /portfolio/server/tmp/.cache \
            /portfolio/ui/.pnpm-store; \
  chown -R "${UID}:${GID}" "/home/${USER}" /portfolio

ENV GOCACHE=/portfolio/server/tmp/.cache \
  GOLANGCI_LINT_CACHE=/portfolio/server/tmp/.cache \
  GOBIN=/usr/local/bin

RUN --mount=type=cache,target=/home/${USER}/.cache/go-build \
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0 \
  && go install github.com/air-verse/air@latest \
  && go install github.com/swaggo/swag/cmd/swag@latest

COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
ARG NPM_VERSION=11.5.2
RUN npm install -g npm@${NPM_VERSION} --no-audit --no-fund --quiet && npm -v
ARG PNPM_VERSION=10.15.0
RUN npm install -g pnpm@${PNPM_VERSION} --no-audit --no-fund --quiet \
  && pnpm --version

WORKDIR /portfolio
ENV UI_DIR=/portfolio/ui
ENV PNPM_HOME=${UI_DIR}/.pnpm \
  PNPM_STORE_PATH=${UI_DIR}/.pnpm-store \
  NPM_CONFIG_USERCONFIG=${UI_DIR}/.npmrc \
  CI=1
ENV PATH="${PNPM_HOME}:${PATH}"
RUN mkdir -p "$PNPM_HOME" && chown -R "${USER}:${GROUP}" "$PNPM_HOME"

COPY --chown=${USER}:${GROUP} . .

WORKDIR /portfolio/server
RUN --mount=type=cache,target=/home/${USER}/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod tidy

WORKDIR /portfolio/ui
RUN --mount=type=cache,target=/portfolio/ui/.pnpm-store \
  pnpm fetch --reporter=append-only || true
RUN --mount=type=cache,target=/portfolio/ui/.pnpm-store \
  pnpm install --offline --reporter=append-only --prefer-offline \
  || pnpm install --frozen-lockfile --reporter=append-only --prefer-offline

WORKDIR /portfolio

RUN chown -R "${USER}:${GROUP}" /portfolio
USER ${UID}:${GID}

EXPOSE 8080 5173
CMD ["./start-dev.sh"]