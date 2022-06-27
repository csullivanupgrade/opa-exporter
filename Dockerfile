FROM golang:1.17.11 as build

WORKDIR /build

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN echo "Building" && \
  task build && \
  echo

## TODO ## ADD TESTS ##
# FROM compile AS test

# RUN echo "Testing" && \
#   task test && \
#   echo

## TODO ## FIX LINT/ENABLE LINTER ##
# FROM compile AS linter

# RUN echo "Linting" && \
#   task lint --timeout 60s && \
#   echo

FROM gcr.io/distroless/static:nonroot

COPY --from=build /build/tmp/opa-scorecard /opa-scorecard

CMD [ "/opa-scorecard", "--incluster=true" ]
