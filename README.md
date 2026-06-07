# Go and Prometheus usage workshop

This workshop demonstrates how to integrate Prometheus and Alertmanager into a Go application. It includes a sample Go server, Prometheus configuration with alerting and recording rules, Alertmanager for alert handling, and k6 for simulating load.

## Project Structure

- `cmd/product/httpserver/main.go`: The Go application with Prometheus instrumentation.
- `prometheus/`: Prometheus configuration, rules, and unit tests.
- `alertmanager/`: Alertmanager configuration.
- `k6/`: Load testing script.
- `docker-compose.yml`: Orchestrates all services.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the workshop

1. Start all services:
   ```bash
   docker-compose up -d
   ```

2. Access the services:
   - **Go Application**: [http://localhost:8080](http://localhost:8080)
     - `/metrics`: Prometheus metrics
     - `/blog/posts/:id`: Sample endpoint (custom metrics)
     - `/products/:id`: Sample endpoint (default metrics)
   - **Prometheus**: [http://localhost:9090](http://localhost:9090)
   - **Alertmanager**: [http://localhost:9093](http://localhost:9093)
   - **Mailcatcher**: [http://localhost:8081](http://localhost:8081) (to see outgoing alert emails)

## Prometheus

### Check configs

You can verify that your rule files are syntactically correct:

```bash
promtool check rules /etc/prometheus/recording_rules.yml
promtool check rules /etc/prometheus/alerting_rules.yml
```

### Unit tests

Prometheus supports unit testing for your rules. Run the tests using:

```bash
promtool test rules /etc/prometheus/rules_test.yml
```

## Alertmanager

Alertmanager handles alerts sent by Prometheus. You can also manually push alerts to it for testing:

To push alerts to Alertmanager:

```bash
curl -X POST http://localhost:9093/api/v2/alerts \
  -H "Content-Type: application/json" \
  -d '[
    {
      "labels": {
        "alertname": "TestAlert",
        "severity": "critical"
      },
      "annotations": {
        "summary": "Test alert fo bar bla bla bla"
      },
      "startsAt": "'$(date -Iseconds)'"
    }
  ]'
```

## Load tests

To simulate traffic and see metrics changing in Prometheus, use the provided k6 script:

```bash
docker-compose run --rm k6 run /scripts/load.js
```
