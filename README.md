# replicator

Kubernetes operator that replicates resources, with templating functionality to namespaces matching the provided label selector

## Templating

In the templated resources, you can use variables on the form `[[ .Values.<key> ]]`. 
Values can either be: 
- set directly in the `ReplicationConfig` resource in `spec.templateValues.values` (simplest)
- contained in a secret referred to by `spec.templateValues.secrets` (if it's a secret)

If the value is specific for the namespace you can pick out labels or annotations in the target namespace by enumerating them in `spec.templateValues.namespace.{labels,annotations}`
  - If keys are formatted as url, e.g. `foo.bar.acme/key`, they will be normalized into `key`

## Example

```yaml
apiVersion: nais.io/v1
kind: ReplicationConfig
metadata:
  name: team-resources
spec:
  namespaceSelector:
    matchExpressions:
      - key: team
        operator: Exists
    #matchLabels:
    #  team-namespace: "true"
  templateValues:
    values: 
      project: abc-123
    secrets:
      - name: secret-containing-value
    namespace:
      labels:
        - team
      annotations:
        - beam
  resources:
    - template: |
        kind: Secret
        apiVersion: v1
        type: kubernetes.io/Opaque
        metadata:
          name: replicator-secret
        stringData:
          apiKey: [[ .Values.apikey ]] # loaded from secret 
    - template: |
        apiVersion: core.cnrm.cloud.google.com/v1beta1
        kind: ConfigConnectorContext
        metadata:          
          name: configconnectorcontext.core.cnrm.cloud.google.com
        spec:
          googleServiceAccount: cnrm-[[ .Values.team ]]@[[ .Values.project ]].iam.gserviceaccount.com
```

## Force reconciliation of resource
If you want to trigger a reconciliation of a ReplicationConfig you can patch the `ReplicationConfig` resource, removing the `status.synchronizationHash` field using the command: `kubectl patch repconf <name> -p '[{"op": "remove", "path": "/status/synchronizationHash"}]' --type=json`. 
