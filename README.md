# replicator

Kubernetes operator that replicates resources, with templating functionality to namespaces matching the provided label selector

## Templating

In the templated resources, you can use variables on the form `[[ .Values.<key> ]]`. 
Values can either be: 
- set directly in the `ReplicationConfig` resource in `spec.values` (simplest)
- contained in a secret referred to by `spec.valueSecrets` (if it's a secret)
- be set in a annotation in the target namespace on the form `replicator.nais.io/<key>: <value>` (if value is context/namespace specific)

## Example

```yaml
apiVersion: nais.io/v1
kind: ReplicationConfig
metadata:
  name: team-resources
spec:
  namespaceSelector:
    matchLabels:
      team-namespace: "true"
  values: 
    project: abc-123
  valueSecrets:
    - name: secret-containing-value
  resources:
    - template: |
          kind: Secret
          apiVersion: v1
          type: kubernetes.io/Opaque
          name: replicated-secret
          stringData: |
              apiKey: [[ .Values.apikey ]] # loaded from secret 
    - template: |
          apiVersion: core.cnrm.cloud.google.com/v1beta1
          kind: ConfigConnectorContext
          metadata:          
            name: configconnectorcontext.core.cnrm.cloud.google.com
          spec:
            googleServiceAccount: cnrm-[[ .Values.teamname ]]@[[ .Values.project ]].iam.gserviceaccount.com # teamname value would here be set from annotation on targeted namespace on the form: `replicator.nais.io/teamname: team`
```

## Development

Create binary (includes generating go code and new manifests):

```make build```

Generate go code: 

```make generate```

Make new manifests through kubebuilder and kustomize:

```make manifests```

Running tests (will also generate code and create new manifests):

```make test```

PRs are always welcome!

## Running in local cluster

Set up local cluster:

[Set up kind cluster](https://book.kubebuilder.io/reference/kind.html) (or equivalent)

Load image into kind cluster:

```make kind```

If you want to test locally with webhook [install cert-manager](https://book.kubebuilder.io/cronjob-tutorial/cert-manager.html) in your local cluster and set enable-webhook=true in "run" target in the Makefile.

Install:

```make install```

Run:

```make run```

