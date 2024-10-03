# webhooks-adaptor

## Try it out
```
docker run --net=host -v $(pwd)/cmd/data:/opt/webhooks-adaptor/data gcr.io/unique-caldron-775/antony/webhooks-adaptor:latest
```

Edit the `data.json` and `data.template` on your local disk to update the adaptor logic. They are read at every request, so no need to restart the webhooks-adaptor, just keep it running in the background.


## Test

```
curl -XPOST http://localhost:8090/webhooks/data --data '{
  "id": "",
  "time": "2024-10-02T10:40:35+00:00",
  "description": "[TEST] Traffic inside your cluster triggered Web Application Firewall rules.",
  "origin": "Web Application Firewall",
  "severity": 80,
  "type": "",
  "attack_vector": "Network",
  "mitre_tactic": "Initial Access",
  "mitre_ids": [
    "T1190"
  ],
  "mitigations": [
    "This Web Application Firewall event is generated for the purpose of webhook testing, no action is required.",
    "Payload of this event is consistent with actual expected payload when a similar event happens in your cluster."
  ],
  "record": {
    "@timestamp": "2024-01-01T12:00:00.000000000Z",
    "destination": {
      "hostname": "",
      "ip": "10.244.151.190",
      "name": "frontend-7d56967868-drpjs",
      "namespace": "online-boutique",
      "port_num": "8080"
    },
    "host": "aks-agentpool-22979750-vmss000000",
    "level": "",
    "method": "GET",
    "msg": "WAF detected 2 violations [deny]",
    "path": "/test/artists.php?artist=0+div+1+union%23foo*%2F*bar%0D%0Aselect%23foo%0D%0A1%2C2%2Ccurrent_user",
    "protocol": "HTTP/1.1",
    "request_id": "460182972949411176",
    "rules": [
      {
        "disruptive": "true",
        "file": "/etc/modsecurity-ruleset/@owasp_crs/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
        "id": "942100",
        "line": "5195",
        "message": "SQL Injection Attack Detected via libinjection",
        "severity": "critical"
      },
      {
        "disruptive": "true",
        "file": "/etc/modsecurity-ruleset/@owasp_crs/REQUEST-949-BLOCKING-EVALUATION.conf",
        "id": "949110",
        "line": "6946",
        "message": "Inbound Anomaly Score Exceeded (Total Score: 5)",
        "severity": "emergency"
      }
    ],
    "source": {
      "hostname": "",
      "ip": "10.244.214.122",
      "name": "busybox",
      "namespace": "online-boutique",
      "port_num": "33387"
    }
  },
  "geo_info": {},
  "labels": {
    "Cluster": "cluster"
  }
}'
```


## Build go binary (required for docker build)
```
go build -o cmd/build/webhooks-adaptor cmd/main.go
```

## Build and run

```
docker build -t webhooks-adaptor cmd
docker run --net=host webhooks-adaptor:latest
```

## Build and run (while testing changes to the data)
The `data.json` and `data.template` is read from disk at every request. Not optimised for performance,
but is allows to update the data between requests for quick testing.

```
docker build -t webhooks-adaptor cmd
docker run --net=host -v $(pwd)/cmd/data:/opt/webhooks-adaptor/data webhooks-adaptor:latest
```
