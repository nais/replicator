# replicator

Kubernetes operator that replicates templated resources to namespaces matching the label selector.

## Templating

In the templated resources, you can use variables on the form `[[ .Values.<key> ]]`. 
Values can either be: 
- set directly in the `ReplicationConfig` resource in `spec.templateValues.values` (simplest)
- contained in a secret referred to by `spec.templateValues.secrets` (if it's a secret)

Optionally you can base64 encode the value inserted in the template by:
`[[ index .Values "key" | b64enc ]]`

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
      - name: secret-containing-tls-cert
      - name: secret-that-doesnt-exist-yet
        validate: false
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
          apiKey: [[ .Values.apikey ]] # loaded from secret-containing-value
    - template: |
        kind: Secret
        apiVersion: v1
        type: kubernetes.io/tls
        metadata:
          name: replicator-tls-secret
        data:
          tls.key: [[ index .Values "tls.key" | b64enc ]] # loaded from secret-containing-tls-cert
          tls.crt: [[ index .Values "tls.crt" | b64enc ]] # loaded from secret-containing-tls-cert
    - template: |
        apiVersion: core.cnrm.cloud.google.com/v1beta1
        kind: ConfigConnectorContext
        metadata:          
          name: configconnectorcontext.core.cnrm.cloud.google.com
        spec:
          googleServiceAccount: cnrm-[[ .Values.team ]]@[[ .Values.project ]].iam.gserviceaccount.com
```

## Force reconciliation of resource

If you want to trigger a reconciliation of a ReplicationConfig, patch the `ReplicationConfig` resource and remove the `status.synchronizationHash` field using this command:

```shell
kubectl patch repconf <name> \
  --type=json \
  --subresource=status \
  -p '[{"op": "remove", "path": "/status/synchronizationHash"}]'
```
