# Monitoring (Metrics + Traces)

## 1. Amaç
Sistem sağlığı gözlemlenebilir. Prometheus (metric) + Jaeger/Tempo (trace) + Grafana (visualize) stack.

## 2. Stack

| Layer | Tool |
|---|---|
| Metric collector | Prometheus (pull scrape) |
| Metric storage | Prometheus TSDB local + Thanos/Cortex (long-term) Phase 02 |
| Traces | OpenTelemetry SDK + Jaeger or Grafana Tempo |
| Log | Loki (bkz. [logging.md](logging.md)) |
| Dashboard | Grafana (unified metric+log+trace) |
| Alert | Alertmanager (bkz. [alerting.md](alerting.md)) |
| APM (Phase 02) | Grafana Application Observability |

## 3. Metrics — Scrape

Prometheus scrape targets:
- App API: `/metrics` (prom-client Go lib).
- Worker: `/metrics`.
- HAProxy: haproxy_exporter.
- Postgres: postgres_exporter.
- Redis: redis_exporter.
- Node: node_exporter (host CPU/mem/disk).
- K8s: kube-state-metrics.
- Blackbox: probe URL uptime.

Interval: 15s default; kritik 10s.

## 4. Metric Naming

Prometheus convention:
```
otter_<domain>_<subject>_<unit>_[total|sum|count|bucket]
```

Örnek:
```
otter_activity_reveal_duration_seconds{quantile}
otter_ledger_entry_total{kind,currency}
otter_ws_conn_gauge
otter_rate_limit_hit_total{endpoint,boyut}
```

Labels sparingly (cardinality budget).

## 5. Traces — OpenTelemetry

- SDK Go: `otelhttp` middleware; each request → span.
- Trace context propagation W3C traceparent header.
- Sampling: 10% baseline; 100% error path; 100% slow (> p95).
- Backend: Jaeger UI veya Grafana Tempo.

Correlation:
- Trace ID = correlation_id (audit log ile match).
- Log include trace_id in structured field.
- Grafana Tempo "trace to log" one-click navigate.

## 6. Grafana Dashboards (v1 taslak)

### Overview
- CCU (WS conn gauge).
- Request rate / error rate / p95 latency (RED metric).
- DB / Redis health.
- SLO burn rate (bkz. [performance/slo-slis.md](../performance/slo-slis.md)).

### Per-Domain
- **Activity/AFK**: contract active, reveal rate, analytical calc ms.
- **Combat**: session count, skill cast rate, damage pipeline.
- **Market**: search cache hit, buy latency, listing gauge.
- **Ekonomi**: ledger rate, wallet aggregate, sink/source ratio.
- **Auth**: login success/fail, session count.
- **WS**: conn count, fanout latency, backpressure drop.

### Infra
- **Postgres**: connections, query rate, cache hit, replication lag.
- **Redis**: memory, ops/sec, pubsub subscribers.
- **HAProxy**: backend up, queue, request rate.
- **Node**: CPU / memory / disk / network.

### Alerts
- Aktif alert list.
- Recent fired.
- Silence management.

## 7. SLO Dashboard

Per-endpoint:
- Target vs actual p95.
- Error rate.
- Error budget remaining %.
- Burn rate (1h, 24h, 7d).

Auto freeze trigger: burn > 100% → release freeze (DD-201).

## 8. Recording Rules

Prometheus precomputed:
```yaml
- record: otter:endpoint:latency:p95
  expr: histogram_quantile(0.95, sum by (endpoint, le)(rate(otter_endpoint_request_duration_seconds_bucket[5m])))
```

Reduces query load.

## 9. Retention

- Prometheus local: 15 gün high-res.
- Thanos/Cortex (Phase 02): 12+ ay downsampled.
- Grafana query federated.

## 10. Cost

Managed Grafana Cloud ~$50-500/ay (data volume-based).
Self-hosted: k8s persistent storage + admin time.

## 11. Metrikler (Meta)
- `prometheus_scrape_duration_seconds{job}`
- `prometheus_target_up{job}`
- `grafana_dashboard_view_total{dashboard}`

## 12. Cross-Reference
- [logging.md](logging.md), [alerting.md](alerting.md)
- [performance/slo-slis.md](../performance/slo-slis.md)
- [security/audit.md](../security/audit.md) — correlation_id

## 13. Bekleyen Kararlar
- Prometheus long-term storage (Thanos vs Cortex vs managed) — Phase 02.
- Trace backend (Jaeger vs Tempo) — Phase 01.
- Grafana Cloud vs self-hosted — ops cost karar.
- APM full (Application Observability) Phase 02.
