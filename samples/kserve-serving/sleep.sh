#!/bin/sh

# To use with RedHat OpenShift Service Mesh, set
# SERVICE_MESH=servicemesh

if [ -z ${SERVICE_MESH+x} ]; then
  SERVICE_MESH="istio"
fi

if [ -z ${NS+x} ]; then
  NS="default"
fi

MANIFEST=/tmp/manifest.$$
cat <<EOF > $MANIFEST
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sleep
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sleep
  template:
    metadata:
      labels:
        app: sleep
EOF
if [ "${SERVICE_MESH}" = "istio" ]; then
cat <<EOF >> $MANIFEST
        sidecar.istio.io/inject: "true"
EOF
elif [ "${SERVICE_MESH}" = "servicemesh" ]; then
cat <<EOF >> $MANIFEST
      annotations:
        sidecar.istio.io/inject: "true"
EOF
fi
cat <<EOF >> $MANIFEST
    spec:
      containers:
      - name: sleep
        image: curlimages/curl
        command: ["/bin/sh", "-c", "sleep 3650d"]
        workingDir: /demo
        imagePullPolicy: IfNotPresent
        volumeMounts:
        - name: config-volume
          mountPath: /demo
        securityContext:
          runAsNonRoot: true
          runAsUser: 1001040000
          allowPrivilegeEscalation: false
      volumes:
      - name: config-volume
        configMap:
          name: demo-input
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo-input
data:
  input.json: |
    {
      "inputs": [
        {
          "name": "input-0",
          "shape": [2, 4],
          "datatype": "FP32",
          "data": [
            [6.8, 2.8, 4.8, 1.4],
            [6.0, 3.4, 4.5, 1.6]
          ]
        }
      ]
    }
  wisdom.sh: |
    curl -H 'Content-Type: application/json' http://wisdom.default -d @input.json -s -D - \
    | grep -e HTTP -e app-version
  wisdom-test.sh: |
    curl -H 'Content-Type: application/json' http://wisdom.default -d @input.json -s -D - \
    -H 'traffic: test' \
    | grep -e HTTP -e app-version
EOF

kubectl -n $NS apply -f $MANIFEST
rm -f $MANIFEST
