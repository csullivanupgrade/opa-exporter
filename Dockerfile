FROM golang:1.18.3 AS build

ENV TASK_VER 3.13.0
ENV GOLANGCI_LINT_VER 1.46.2

WORKDIR /build

RUN echo "Fetching dev dependencies" && \
  go install github.com/go-task/task/v3/cmd/task@v${TASK_VER} && \
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v${GOLANGCI_LINT_VER}

COPY . .

RUN echo "Building" && \
  task build-bin && \
  echo

RUN echo "Linting" && \
  task lint && \
  echo

FROM gcr.io/distroless/static:nonroot

COPY --from=build /build/tmp/opa-exporter /opa-exporter

CMD [ "/opa-exporter" ]
