# Kubernetes files

Can run in any Kubernetes environment `minikube` for example 

Add the API key to the configuration map where indicated.

Run in the following order

```
kubectl apply -f config-map.yml
kubectl apply -f deployment.yml
kubectl apply -f service.yml
```

Access the service either on the node port 30000, highly dependent on your local setup.

To deploy changes run `kubectl rollout restart deployment/stock-ticker`