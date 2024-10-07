# Header based example

Install Argo Rollouts on Cluster
with Traefik 3.x and Gateway plugin 0.4.0 as described in the
README of the `traefik-setup` folder.


Create main route

```
kubectl apply -f smart-route.yml
```

Create rollout 

```
kubectl apply -f rollout.yml
```

Test that nothing is happening and all calls see version 1

```
curl http://localhost/callme
curl -H "X-Canary: yes" http://localhost/callme
curl -H "X-Canary: maybe" http://localhost/callme
```

Then edit Rollout and change line 51 APP_VERSION to 2.0 and line APP_COLOR to something else.

Apply the manifest again to start a canary process

Then try again 

```
curl -H "X-Canary: yes" http://localhost/callme
curl -H "X-Canary: maybe" http://localhost/callme
```

You will see that the first call always goes to the new version
while the second one will go to either version.
