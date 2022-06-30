FROM golang:1.18.3 as build

WORKDIR /build

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN echo "Building" && \
  task build-bin && \
  echo

## TODO ## ADD TESTS ##
# FROM build AS test

# RUN echo "Testing" && \
#   task test && \
#   echo

FROM build AS linter

RUN echo "Linting" && \
  task lint && \
  echo

FROM gcr.io/distroless/static:nonroot

COPY --from=build /build/tmp/opa-exporter /opa-exporter

CMD [ "/opa-exporter" ]
