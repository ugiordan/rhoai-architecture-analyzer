# Secrets Inventory

10 secrets referenced across the platform. No secret values are extracted, only names, types, and which component references them.

## Secret Distribution

<div markdown class="bar-chart-container" style="margin: 1em 0; padding: 1em; border: 1px solid var(--md-default-fg-color--lightest); border-radius: 8px;">

**Secrets per Component**

<div style="display: flex; flex-direction: column; gap: 6px;">
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">data-science-pipelines</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 50%; background: #27ae60; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">2</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">data-science-pipelines-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 100%; background: #27ae60; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">4</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">kserve</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 75%; background: #27ae60; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">3</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">modelmesh-serving</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 25%; background: #27ae60; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">1</span>
</div>
</div>
</div>

## Secrets by Component

| Component | TLS | Opaque | Total |
|-----------|-----|--------|-------|
| data-science-pipelines | 0 | 2 | 2 |
| data-science-pipelines-operator | 0 | 4 | 4 |
| kserve | 0 | 3 | 3 |
| modelmesh-serving | 0 | 1 | 1 |

## Secret Detail

Per-component secret breakdown by name and type.

### data-science-pipelines (2 secrets)

| Secret | Type |
|--------|------|
| kfp-api-webhook-cert | Opaque |
| mlpipeline-minio-artifact | Opaque |

### data-science-pipelines-operator (4 secrets)

| Secret | Type |
|--------|------|
| ds-pipeline-db-test | Opaque |
| mariadb-certs | Opaque |
| minio | Opaque |
| minio-certs | Opaque |

### kserve (3 secrets)

| Secret | Type |
|--------|------|
| kserve-webhook-server-cert | Opaque |
| llmisvc-webhook-server-cert | Opaque |
| localmodel-webhook-server-cert | Opaque |

### modelmesh-serving (1 secrets)

| Secret | Type |
|--------|------|
| modelmesh-webhook-server-cert | Opaque |

## Patterns

- **Webhook certs** are the dominant secret type (7 of 10 secrets).

