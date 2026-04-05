0s
30s
1s
4s
Run layer5labs/Kanvas-Snapshot@v0.2.43
Run echo "IS_PLAYGROUND_RUNNING=$( echo $(./action/playground-ping.sh))" >> "$GITHUB_OUTPUT"
Run echo $IS_PLAYGROUND_RUNNING

Run echo "MESHERY_URL=https://playground.meshery.io" >> $GITHUB_ENV
Run echo true
true
Run FILE_PATH=install/deployment_yamls/k8s ./action/manifest-merger.sh # creates file in root dir
install/deployment_yamls/k8s
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-app-mesh
  name: meshery-app-mesh
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-app-mesh
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-app-mesh
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-app-mesh:stable-latest
        imagePullPolicy: Always
        name: meshery-app-mesh
        ports:
        - containerPort: 10005
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-app-mesh
  name: meshery-app-mesh
spec:
  ports:
  - name: "10005"
    port: 10005
    targetPort: 10005
  selector:
    io.kompose.service: meshery-app-mesh
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-cilium
  name: meshery-cilium
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-cilium
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-cilium
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-cilium:stable-latest
        imagePullPolicy: Always
        name: meshery-cilium
        ports:
        - containerPort: 10012
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-cilium
  name: meshery-cilium
spec:
  ports:
  - name: "10012"
    port: 10012
    targetPort: 10012
  selector:
    io.kompose.service: meshery-cilium
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-consul
  name: meshery-consul
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-consul
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-consul
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-consul:stable-latest
        imagePullPolicy: Always
        name: meshery-consul
        ports:
        - containerPort: 10002
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  # creationTimestamp: null
  labels:
    io.kompose.service: meshery-consul
  name: meshery-consul
spec:
  ports:
  - name: "10002"
    port: 10002
    targetPort: 10002
  selector:
    io.kompose.service: meshery-consul
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery
  name: meshery
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery
    spec:
      serviceAccountName: meshery-server
      containers:
      - env:
        - name: EVENT
          value: mesheryLocal
        - name: PROVIDER_BASE_URLS
          value: https://cloud.layer5.io
        - name: ADAPTER_URLS
          value: meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-nsm:10004 meshery-app-mesh:10005 meshery-kuma:10007 meshery-nginx-sm:10010
        image: meshery/meshery:stable-latest
        imagePullPolicy: Always
        name: meshery
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /healthz/live
            port: 8080
          initialDelaySeconds: 80
          periodSeconds: 12
          failureThreshold: 4
        readinessProbe:
          httpGet:
            path: /healthz/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 4
          failureThreshold: 4
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-istio
  name: meshery-istio
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-istio
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-istio
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-istio:stable-latest
        imagePullPolicy: Always
        name: meshery-istio
        ports:
        - containerPort: 10000
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-istio
  name: meshery-istio
spec:
  ports:
  - name: "10000"
    port: 10000
    targetPort: 10000
  selector:
    io.kompose.service: meshery-istio
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-kuma
  name: meshery-kuma
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-kuma
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-kuma
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-kuma:stable-latest
        imagePullPolicy: Always
        name: meshery-kuma
        ports:
        - containerPort: 10007
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-kuma
  name: meshery-kuma
spec:
  ports:
    - name: "10007"
      port: 10007
      targetPort: 10007
  selector:
    io.kompose.service: meshery-kuma
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-linkerd
  name: meshery-linkerd
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-linkerd
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-linkerd
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-linkerd:stable-latest
        imagePullPolicy: Always
        name: meshery-linkerd
        ports:
        - containerPort: 10001
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-linkerd
  name: meshery-linkerd
spec:
  ports:
  - name: "10001"
    port: 10001
    targetPort: 10001
  selector:
    io.kompose.service: meshery-linkerd
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-nginx-sm
  name: meshery-nginx-sm
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-nginx-sm
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-nginx-sm
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-nginx-sm:stable-latest
        imagePullPolicy: Always
        name: meshery-nginx-sm
        ports:
        - containerPort: 10010
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-nginx-sm
  name: meshery-nginx-sm
spec:
  ports:
  - name: "10010"
    port: 10010
    targetPort: 10010
  selector:
    io.kompose.service: meshery-nginx-sm
status:
  loadBalancer: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-nsm
  name: meshery-nsm
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-nsm
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-nsm
    spec:
      serviceAccount: meshery-server
      containers:
        - image: meshery/meshery-nsm:stable-latest
          imagePullPolicy: Always
          name: meshery-nsm
          ports:
            - containerPort: 10004
          resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-nsm
  name: meshery-nsm
spec:
  ports:
  - name: "10004"
    port: 10004
    targetPort: 10004
  selector:
    io.kompose.service: meshery-nsm
status:
  loadBalancer: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery
  name: meshery
spec:
  ports:
  - name: "http"
    port: 9081
    targetPort: 8080
  selector:
    io.kompose.service: meshery
  type: LoadBalancer


---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-traefik-mesh
  name: meshery-traefik-mesh
spec:
  selector:
    matchLabels:
      io.kompose.service: meshery-traefik-mesh
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        io.kompose.service: meshery-traefik-mesh
    spec:
      serviceAccount: meshery-server
      containers:
      - image: meshery/meshery-traefik-mesh:stable-latest
        imagePullPolicy: Always
        name: meshery-traefik-mesh
        ports:
        - containerPort: 10006
        resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    kompose.cmd: kompose convert -f ../docker-compose.yaml
    kompose.version: 1.32.0 ()
  creationTimestamp: null
  labels:
    io.kompose.service: meshery-traefik-mesh
  name: meshery-traefik-mesh
spec:
  ports:
  - name: "10000"
    port: 10006
    targetPort: 10006
  selector:
    io.kompose.service: meshery-traefik-mesh
status:
  loadBalancer: {}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: meshery-server
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: meshery-server
  labels:
    app: meshery
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs: ["/metrics", "/health", "/ping"]
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: meshery-server
  labels:
    app: meshery
roleRef:
  kind: ClusterRole
  name: meshery-server
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: ServiceAccount
  name: meshery-server
  namespace: meshery
---
---
Run echo "Uploading using application file"
Uploading using application file
jq: parse error: Invalid numeric literal at line 1, column 7
Error: Process completed with exit code 5.
Run actions/upload-artifact@v4
Warning: No files were found with the provided path: action/log.txt. No artifacts will be uploaded.
Run actions/upload-artifact@v4
With the provided path, there will be 1 file uploaded
Artifact name is valid!
Root directory input is valid!
Beginning upload of artifact content to blob storage
Uploaded bytes 5949
Finished uploading artifact content to blob storage!
SHA256 digest of uploaded artifact zip is 91d2d8efcedef588e5124bb52be3db086fb85d6c0095ab2a70877885f6739ec4
Finalizing artifact upload
Artifact meshery-log1.zip successfully finalized. Artifact ID 6275428597
Artifact meshery-log1 has been successfully uploaded! Final size is 5949 bytes. Artifact ID is 6275428597
Artifact download URL: https://github.com/meshery/meshery/actions/runs/23995269920/artifacts/6275428597
Run actions/upload-artifact@v4
Warning: No files were found with the provided path: action/cypress-action/cypress/videos. No artifacts will be uploaded.
Run actions/upload-artifact@v4
  with:
    name: cypress-screenshots
    path: action/cypress-action/cypress/screenshots
    if-no-files-found: warn
    compression-level: 6
    overwrite: false
    include-hidden-files: false
  env:
    PULL_NO: 18362
    MESHERY_URL: https://playground.meshery.io
Warning: No files were found with the provided path: action/cypress-action/cypress/screenshots. No artifacts will be uploaded