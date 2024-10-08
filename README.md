# Static and Dynamic routing for Canaries with Argo Rollouts

This repository contains 3 examples for using Traefik 3.x with Argo Rollouts
via the Gateway API

1. Example with just a [route](simple) 
2. Example with [static routes for preview/stable/canary](static-routing)
3. Example with [header based routing](dynamic-routing).

Read the full blog at TDB

## Traefik setup

This process is common to both examples.

Install traefik 3.x, Argo Rollouts and the Gateway API plugin as described [traefik-setup/README.md](traefik-setup/README.md) in a local Kubernetes cluster.

Expose the Traefik endpoint as Load Balancer so you can visit it via the browser.

## Simple HTTP route example

To test just the Gateway API implementation of Traefik apply the manifests
at `simple` folder. Then visit `http://localhost/example/` to see the application.

## Static routing example

Apply the manifests in the `static-routing` folder. Then start a canary
by editing the `rollout.yml` file and change the docker image to different tags.

Then visit the 3 urls

* `http://localhost/stable` This is always old version
* `http://localhost/canary` This is active canary version
* `http://localhost/preview` This is always new  version

## Dynamic routing example

Apply the manifests in the `dynamic-routing` folder. Then start a canary
by editing the `rollout.yml` file and change the APP_VERSION and APP_COLOR to start a canary

Then visit `http://localhost` and see that if you change the active header on the UI, the application can be forced to be part of the canary.




