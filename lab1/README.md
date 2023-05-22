Commands to run

`kubectl apply -f k8s/service1/`

`kubectl apply -f k8s/service2/`

`kubectl proxy`

`curl http://localhost:8001/api/v1/namespaces/default/services/service1-service/proxy/service1/read`

`curl http://localhost:8001/api/v1/namespaces/default/services/service1-service/proxy/service2/read`