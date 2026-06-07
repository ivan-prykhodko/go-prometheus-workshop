# Go and Prometheus usage workshop

## Prometheus

### Check configs

```bash
promtool check rules /etc/prometheus/recording_rules.yml
promtool check rules /etc/prometheus/alerting_rules.yml
```

### Unit tests

```bash
promtool test rules /etc/prometheus/rules_test.yml
```

## Alertmanager

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

```bash
docker-compose run --rm k6 run /scripts/load.js
```
