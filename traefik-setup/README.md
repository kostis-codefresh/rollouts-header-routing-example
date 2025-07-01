# Using Traefik Gateway API with Argo Rollouts

This guide will describe how to use Traefik proxy an an implementation
for the Gateway API in order to do split traffic with Argo Rollouts.

Versions used

* Argo Rollouts [1.7.2](https://github.com/argoproj/argo-rollouts/releases)
* Argo Rollouts Gateway API plugin [0.4.0](https://github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/releases)
* [Traefik 3.1.4](https://doc.traefik.io/traefik/getting-started/install-traefik/)
* GatewayAPI 1.1 (Part of the [Traefik Helm chart](https://github.com/traefik/traefik-helm-chart))

Note that Argo Rollouts also [supports Traefik natively](https://argoproj.github.io/argo-rollouts/features/traffic-management/traefik/).

## Step 1 - Enable Gateway Provider and create Gateway entrypoint

First let's install Traefik as a Gateway provider. Follow the official [installation instructions](https://doc.traefik.io/traefik/getting-started/install-traefik/).

You should also read the documentation on how [Traefik implements the Gateway API](https://doc.traefik.io/traefik/providers/kubernetes-gateway/).

Install Traefik with Gateway support

```
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm install traefik traefik/traefik --set experimental.kubernetesGateway.enabled=true --set providers.kubernetesGateway.enabled=true --set ingressRoute.dashboard.enabled=true --version v32.0.0 --namespace=traefik --create-namespace
```

Note that using Helm automatically installs the Kubernetes Gateway API CRDs
as well as the appropriate RBAC resources so that Traefik can manage
HTTP routes inside your cluster.

Enabling the Traefik dashboard is optional, but helpful when debugging
routes.

After initial installation you can expose the dashboard with

```
kubectl port-forward -n traefik $(kubectl -n traefik get pods --selector "app.kubernetes.io/name=traefik" --output=name) 9000:9000
```

And then visit `http://127.0.0.1:9000/dashboard/` (or whatever is the IP
of your Loadbalancer)

Also notice that the Helm chart creates a Loadbalancer service by default.
If your cluster has already a Loadbalancer or you want to customize Traefik installation you need to pass your own [options](https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml).

## Step 2 - Create GatewayClass and Gateway resources

After installing Traefik with need a GatewayClass and an actual Gateway.

The Helm chart already created a GatewayClass for you.

You can verify it with

```
kubectl get gatewayclass
```

Make sure that the value returned is "True" in the "Accepted" column.

Now let's create a Gateway

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: traefik-gateway
spec:
  gatewayClassName: traefik
  listeners:
    - protocol: HTTP
      name: web
      port: 8000 # Default endpoint for Helm chart
```

Apply the file it with

```
kubectl apply -f gateway.yml
```

Notice that we installed the gateway on the default namespace which is where our application will be deployed as well. If you want the gateway to honor routes from other namespaces you need to install Traefik with a different option for `namespacePolicy` in the Helm chart.

This concludes the setup that is specific to Traefik Proxy. The rest of the steps are generic to any implementation of the Gateway API.


## Step 3 - Give access to Argo Rollouts for the Gateway/Http Route

Create Cluster Role resource with needed permissions for Gateway API provider.

```yaml title="cluster-role.yml"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: gateway-controller-role
  namespace: argo-rollouts
rules:
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - "*"
```

__Note:__ These permission are not very strict. You should lock them down according to your needs.

With the following role we allow Argo Rollouts to have write access to HTTPRoutes and Gateways.

```yaml title="cluster-role-binding.yml"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: gateway-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: gateway-controller-role
subjects:
  - namespace: argo-rollouts
    kind: ServiceAccount
    name: argo-rollouts
```

Apply both files with `kubectl`:

```shell
kubectl apply -f cluster-role.yml
kubectl apply -f cluster-role-binding.yml
```

## Step 4 - Create an example Rollout

See folder `static-routing` for 3 static URLs

See folder `dynamic-routing` for a header-based example.



Finished!







