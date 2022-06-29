# opa-exporter

A simple exporter for monitoring OPA Constraint Violations in realtime.

Inspired by this
[blog
post](https://itnext.io/expose-open-policy-agent-gatekeeper-constraint-violations-with-prometheus-and-grafana-6b7ac92ea07f),
from which the exporter code itself originates.

## Why fork?

The source was written as a proof of concept. My needs dictate a production-ready product with timely updates. This fork
intends to be that.

## Contributing

Pull requests are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) before beginning.

This repo makes heavy use of [Taskfiles](https://taskfile.dev). Install this first.

Please check that you have the other necessary prerequisites installed by running:

```
task req
```

### Setting up an environment

First create your cluster:

```
task kind:create
```

Deploy gatekeeper and observability stack:

```
task gk:deploy
# OR IF USING EXTERNAL DATA IN YOUR POLICIES:
task gk:deploy ED_ENABLED=true
```

Deploy the prometheus stack:

```
task prom:deploy
```

Build the container image:

```
task build
```

Deploy the container image:

```
task deploy
```

Run port-forwarding for the local grafana:

```
task prom:port-forward
```

At this point you may open grafana at localhost:3000 and view the OPA dashboard. You can test your policies and see
the violations appear in this dashboard.
