# llm-d-routing-sidecar: RBAC

ServiceAccount bindings, roles, and resource permissions.

## RBAC Hierarchy

```mermaid
graph TD
    %% RBAC hierarchy for llm-d-routing-sidecar
    classDef sa fill:#4a90d9,stroke:#2c5f8a,color:#fff
    classDef role fill:#e8a838,stroke:#b07828,color:#fff
    classDef resource fill:#5cb85c,stroke:#3d8b3d,color:#fff

    sa_2["Group: system:authenticated"] -->|bound via exec-rolebinding| rb_1["exec-rolebinding"]
    class sa_2 sa
    rb_1 -->|grants| r_exec_role["Role: exec-role"]
    class r_exec_role role
    sa_4["Group: system:authenticated"] -->|bound via exec-rolebinding| rb_3["exec-rolebinding"]
    class sa_4 sa
    rb_3 -->|grants| r___PROJECT_NAME__exec_role["Role: ${PROJECT_NAME}-exec-role"]
    class r___PROJECT_NAME__exec_role role
    sa_6["ServiceAccount: placeholder"] -->|bound via ssrf-allowlist-rolebinding| rb_5["ssrf-allowlist-rolebinding"]
    class sa_6 sa
    rb_5 -->|grants| r_ssrf_allowlist_role["Role: ssrf-allowlist-role"]
    class r_ssrf_allowlist_role role
    r_exec_role -->|create| res_7["core: pods/exec"]
    class res_7 resource
    r_exec_role -->|create| res_8["core: pods/exec"]
    class res_8 resource
    r_ssrf_allowlist_role -->|get, watch, list| res_9["core: pods"]
    class res_9 resource
    r_ssrf_allowlist_role -->|get, watch, list| res_10["inference.networking.x-k8s.io: inferencepools"]
    class res_10 resource
```

## Bindings

Subject-to-role mappings defining who has access to what.

| Binding | Type | Role | Subject |
|---------|------|------|---------|
| exec-rolebinding | RoleBinding | exec-role | Group/system:authenticated |
| exec-rolebinding | RoleBinding | ${PROJECT_NAME}-exec-role | Group/system:authenticated |
| ssrf-allowlist-rolebinding | RoleBinding | ssrf-allowlist-role | ServiceAccount/placeholder |

## Role Details

Per-rule breakdown of API groups, resources, and verbs for each role.

| Role | Kind | API Groups | Resources | Verbs |
|------|------|------------|-----------|-------|
| exec-role | Role |  | pods/exec | create |
| exec-role | Role |  | pods/exec | create |
| ssrf-allowlist-role | Role |  | pods | get, watch, list |
| ssrf-allowlist-role | Role |  | inferencepools | get, watch, list |

