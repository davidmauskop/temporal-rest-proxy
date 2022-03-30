# Temporal REST to gRPC Proxy Server

This is a fork of [Temporal's UI
server](https://github.com/temporalio/ui-server) with a new entrypoint at
[rest-proxy/main.go](rest-proxy/main.go). This new entrypoint is a server that
exposes rpcs from the Temporal workflow service as REST API endpoints.
Documentation for these endpoints can be found at
https://temporal-rest-api-docs.onrender.com/. Example usage can be found at
[examples/client-js/index.js](examples/client-js/index.js). This [OpenAPI
schema](server/generated/openapi/service.swagger.json) can be used to generate
client bindings in your language of choice.

## Motivation

In some scenarios it's inconvenient or impossible to use gRPC. For those
scenarios, we'd like the ability to start and interact with Temporal workflows
via REST. In particular, this server is designed to run on
[Render](https://render.com) alongside a [Temporal
cluster](https://github.com/render-examples/temporal). It is meant to support
use cases where there's some application code running outside of Render that
needs access to the Temporal cluster running on Render.

## Deploying

- Create a new repository using the [render-examples/temporal](https://github.com/render-examples/temporal) template.
- Uncomment the `rest-to-grpc-proxy` service definition in the [render.yaml file](https://github.com/render-examples/temporal/blob/main/render.yaml). Commit and push the change.
- Follow the steps outlined in the [Temporal on Render deployment
  guide](https://render.com/docs/deploy-temporal) to create your Temporal
  cluster with a Render Blueprint. Enter or generate an `AUTH_TOKEN` environment
  variable when prompted.

## Authentication

The proxy server is secured with a single Bearer token, supplied as the `AUTH_TOKEN`
environment variable. To enable zero-downtime token rotation, there's an
additional `EXPIRING_AUTH_TOKEN` environment variable that, when set, acts as a
secondary Bearer token. Most of the time it should be empty.

Follow these steps to perform a zero-downtime token rotation:

- Copy the `AUTH_TOKEN` value that you're about to expire into `EXPIRING_AUTH_TOKEN`.
- Set `AUTH_TOKEN` to a new value. This value can be generated in the Render dashboard or elsewhere.
- Deploy the server. It will now accept both the about-to-expire token and the new token.
- Update all clients to use the newly generated token.
- Set `EXPIRING_AUTH_TOKEN` to the empty string or whitespace, and re-deploy the
  server. Now only the new token will be accepted.

## :warning: Production Readiness :warning:

The proxy server has not been tested at scale. Please perform your own tests before choosing it for production workloads, and file issues here as needed.
