# replicator
Replicates resources to namespaces

```yaml
apiVersion: nais.io/v1
kind: ReplicatorConfiguration
metadata:
  name: alerting
spec:
  namespaceSelector:
    matchLabels:
      team-namespace: true
  values: 
    secret:
      - name: slack-webhook # loads data from secret (data: foo: bar)
  resources:
    - kind: Secret
      apiVersion: v1
      type: kubernetes.io/Opaque
      name: slack-webhook
      data:
        stringData: |
          apiKey: {{ .Values.foo }} # loaded from secret
    - kind: AlertmanagerConfig
      apiVersion: monitoring.coreos.com
      name: nais-alerts
      data:
        receivers:
          - name: {{ .Values.team }}-slack # expected to be set as annotation on target namespace (replicator.nais.io/team: something)
            slackConfigs:
              - apiURL:
                  key: apiUrl
                  name: slack-webhook
                sendResolved: true
                channel: {{ .Values.team }}-dev-alerts # dev here could be interpolated, but probably simplest to handle this in templating when deploying the replicatorconfig itself (typically with helm)
                color: {{ "'{{ template \"slack.color\" . }}'" }}
                text: {{ "'{{ template \"slack.text\" . }}'" }}
                title: {{ "'{{ template \"slack.title\" . }}'" }}
        route:
          receiver: {{ .Values.team }}-slack
          groupBy:
            - alertname
          groupInterval: 5m
          groupWait: 10s
          repeatInterval: 1h      
```
