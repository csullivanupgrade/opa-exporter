FROM golang:1.17.11 as build

WORKDIR /app

COPY . .

RUN echo "Building" && \
  GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build . && \
  echo

## TODO ## ADD TESTS ##
# FROM compile AS test

# RUN echo "Testing" && \
#   go test .../. && \
#   echo

## TODO ## FIX LINT/ENABLE LINTER ##
# FROM compile AS linter

# RUN echo "Linting" && \
#   golangci-lint run ./... --timeout 60s && \
#   echo

FROM gcr.io/distroless/static:nonroot

COPY --from=build /app/opa-scorecard /app/opa-scorecard

CMD [ "/app/opa-scorecard", "--incluster=true" ]
