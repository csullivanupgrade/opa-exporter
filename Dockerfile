# syntax = docker/dockerfile:1-experimental

FROM golang:1.19.3 AS build

ENV TASK_VER 3.13.0
ENV GOLANGCI_LINT_VER 1.49.0

WORKDIR /build

RUN echo "Fetching dev dependencies" && \
  go install github.com/go-task/task/v3/cmd/task@v${TASK_VER} && \
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v${GOLANGCI_LINT_VER}

COPY . .

RUN \
  --mount=type=cache,target=/root/.cache/go-build \
  echo "Building" && \
  task build-bin && \
  echo

RUN \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/.cache/golangci-lint \
  echo "Linting" && \
  task lint TIMEOUT=120s && \
  echo

FROM gcr.io/distroless/static:nonroot

COPY --from=build /build/tmp/opa-exporter /opa-exporter

CMD [ "/opa-exporter" ]
