# Alerting

## 1. Amaç
Kritik + warn seviye olayların doğru kişiye zamanında ulaşması. Alert fatigue önleme.

## 2. Stack

- **Alertmanager** (Prometheus) — route + group + silence + escalate.
- **Delivery**: Slack (info/warn), PagerDuty (high/critical), email (özet).
- **Runbook**: her alert bir wiki page linkli.

## 3. Severity 4-Tier (DD-191)

| Severity | Delivery | Örnek |
|---|---|---|
| `info` | Slack #ops-info (kanal, no ping) | Backup success |
| `warn` | Slack #eng (channel ping) | Replica lag > 60s |
| `high` | Slack #eng-critical + inbox | Ledger drift < 1c, credential stuffing |
| `critical` | PagerDuty on-call + Slack | Sum violation, primary DB down, sustained downtime |

## 4. Alert Kategorileri

### Infrastructure
- Node CPU > 90% 15dk (warn).
- Node memory > 90% (warn).
- Disk > 85% (high).
- Postgres replica lag > 60s (warn); > 300s (high).
- Redis memory > 90% (warn); persistence fail (critical).
- HAProxy backend down (warn); all backends down (critical).

### Application
- 5xx rate > 1% (warn); > 5% (high).
- p95 latency > SLO target (warn).
- SLO error budget > 80% consumed (high); > 100% (critical).
- Worker queue depth > 10k (warn).
- Cron job fail (warn); 3 ardışık fail (high).

### Security
- Credential stuffing detected (high).
- Refresh token reuse (P1 high).
- Sum violation (P0 critical).
- Anomaly rule fire high severity (variable).

### Business
- Ledger drift ≥ 1c (critical).
- Backup fail (high).
- Cert expire < 7d (high).

## 5. Alert Rule Örnek (Prometheus)

```yaml
- alert: HighErrorRate
  expr: sum(rate(otter_endpoint_request_errors_total[5m])) 
        / sum(rate(otter_endpoint_request_total[5m])) > 0.05
  for: 5m
  labels:
    severity: high
  annotations:
    summary: "5xx rate > 5% for 5m"
    runbook: "https://wiki/runbook/high-error-rate"

- alert: PostgresPrimaryDown
  expr: up{job="postgres_exporter", role="primary"} == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Postgres primary down"
    runbook: "https://wiki/runbook/pg-primary-down"
```

## 6. Alertmanager Routing

```yaml
route:
  receiver: slack-info
  group_by: [alertname, severity]
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  routes:
    - match: { severity: critical }
      receiver: pagerduty-oncall
      repeat_interval: 1h
      continue: true       # also send to Slack
    - match: { severity: high }
      receiver: slack-critical
      continue: true
    - match: { severity: warn }
      receiver: slack-eng
receivers:
  - name: pagerduty-oncall
    pagerduty_configs: [ { service_key: '...' } ]
  - name: slack-critical
    slack_configs: [ { channel: '#eng-critical', ... } ]
```

## 7. Silencing

- Planned maintenance → silence specific alert 1-4h.
- Silence sadece ops team; audit trail.

## 8. Escalation Policy

PagerDuty:
- Primary on-call (30 dk).
- Secondary (10 dk).
- Team lead (5 dk).
- CTO (5 dk).

Rotation: haftalık primary on-call (weekly handoff meeting).

## 9. Alert Fatigue Prevention

- **Actionable**: alert = someone must act. Non-action → info log, not alert.
- **Grouping**: aynı root cause aggregate.
- **Deduplication**: 5-15 dk window aynı alert bir kez.
- **Tuning**: false positive her yerde review; 2+ false → rule update.

## 10. Runbook Format

Her alert için wiki page:
- **Summary**: bir cümlelik durum.
- **Impact**: kullanıcıya etki.
- **Diagnose**: hangi dashboard/log komut ile check.
- **Mitigate**: geçici çözüm.
- **Resolve**: root cause fix.
- **Escalate**: kime + ne zaman.

## 11. Post-Mortem

Kritik/high fired → post-mortem doc zorunlu:
- Timeline.
- Impact metrics.
- Root cause.
- Fix.
- Prevention.
- Action items with owner + deadline.

Blameless culture.

## 12. Metrikler
- `alertmanager_alerts_firing_total{severity}`
- `alertmanager_notification_success_total{receiver}`
- `alertmanager_notification_fail_total{receiver}`
- `alert_mttr_seconds{severity}` — mean time to resolve (post-mortem input)

## 13. Cross-Reference
- [monitoring.md](monitoring.md), [logging.md](logging.md)
- [security/audit.md](../security/audit.md)
- [performance/slo-slis.md](../performance/slo-slis.md)

## 14. Bekleyen Kararlar
- PagerDuty vs OpsGenie vs FireHydrant — ops karar.
- On-call rotation size (2-4 SRE) Phase 02.
- Blameless post-mortem template Phase 01.
- Alert tuning review sıklık (haftalık) — ops runbook.
