# Stock Ticker

### Context

A simple service to get return and collate stock information for a given number of days and symbol.

Written in Go

Deployed via Kubernetes, see the Kubernetes directory for details on how to deploy it.

See `openapi.yml` for how to interact with the service.

Run tests by running `go test`, not that there are many....

### Resilience I Had Time For

Structured logging with `trace-ids` to make correlating requests simple across multiple users.

Controlled rollouts via Kubernetes deployments

Vendored dependencies for supply chain control

Technically we're pinning CAs, so that could be a security win?

### Suggestions

Could add a custom HTTP client to add retries to the outbound HTTP requests.

Secrets should be passed in via a Kubernetes secret.

Actual integration tests.

Passing variables as a configuration file vs environment variables, can be more ergonomic and help prevent other processes viewing secrets.

If actually deploying this into a cloud, load balancing with a real load balancer vs a node port.

Multiple API keys for the case where one is rate limited

Adding rate limiting middle ware to rate limit requests from users, let's say via IP to start with.

Likely need some form of authentication eventually to also prevent rate limiting on outbound connections. Given it's an API, either a JWT or OAuth potentially.

