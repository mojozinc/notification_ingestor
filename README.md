## Atlan System Design Assignment


### Requirements
#### Setup K8s services
 prerequisite to install minikube - https://minikube.sigs.k8s.io/docs/start/
1. minikube start
2. add kafka to your k8s cluster
    
```
        helm repo add bitnami https://charts.bitnami.com/bitnami
        helm repo update
        helm install kafka bitnami/kafka --set persistence.enabled=false
```
3. build the kafka producer/consumer images
```
        minikube image build -t notification-consumer:latest . -f Dockerfile-consumer
        minikube image build -t notification-ingestor:latest . -f Dockerfile-producer
```
4. apply all k8s yaml files `find k8s -type f | xargs -n1 kubectl apply -f `


## Monte Carlo Ingestor
Once all the pods and services are up, you can test the ingestor service with a sample POST request

1. expose notification-ingestor
`kubectl port-forward svc/notification-ingestor 8080:80`

2. use Postman or curl
```bash
curl --location 'http://127.0.0.1:8080/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "alert_id": "12345",
    "timestamp": "2024-08-21T14:30:00Z",
    "alert_type": "data_quality",
    "severity": "high",
    "details": {
        "alert_category": "column_quality",
        "dataset": "sales_data",
        "table": "transactions",
        "column": "phone_number",
        "condition": "null_value_above_threshold",
        "threshold": 5,
        "current_null_count": 10,
        "total_rows": 1000
    },
    "metadata": {
        "source": "monte_carlo",
        "data_source": "database_name",
        "environment": "production"
    },
    "action_required": true,
    "description": "The column '\''phone_number'\'' in the '\''transactions'\'' table has 10 null values, exceeding the threshold of 5 null values."
}'
```

3. tail consumer logs to see the alert data being persisted
```
kubectl logs -f svc/notification-consumer
```