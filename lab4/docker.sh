docker build -t zhenyanovikov/service1:0.3 -f services/service1/Dockerfile .
docker build -t zhenyanovikov/service2:0.3 -f services/service2/Dockerfile .

docker push zhenyanovikov/service1:0.3
docker push zhenyanovikov/service2:0.3
