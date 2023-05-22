Commands to run

`kubectl apply -f k8s/service1/`

`kubectl apply -f k8s/service2/`

`kubectl apply -f k8s/postgres/`

`minikube service ingress-nginx-controller -n ingress-nginx`

replace port

`curl http://127.0.0.1:52072/service1/read`

replace port

`curl http://127.0.0.1:52072/service2/read`