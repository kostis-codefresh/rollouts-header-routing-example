Instruqtions

Install Argo Rollouts on Cluster
with Traefik 3.x and Gateway plugin 0.4.0


Apply Roles in static-routing/cluster-role and cluster-role-binding

Create main route

```
kubectl apply -f all-routes.yml
```

Create rollout 

```
kubectl apply -f rollout.yml
```

Test that nothing is happening and all calls see version 1

```
curl http://localhost/smart/callme
curl -H "X-Canary: yes" http://localhost/smart/callme
curl -H "X-Canary: no" http://localhost/smart/callme
```

Then edit Rollout and change line 51 APP_VERSION to 2.0

Apply the manifest again to start a canary process

Then try again 

```
curl http://localhost/smart/callme
curl -H "X-Canary: yes" http://localhost/smart/callme
curl -H "X-Canary: no" http://localhost/smart/callme
```

The first two call will succeed. 
The last one fails with 404 (which is the one match by the header route)
