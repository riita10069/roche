package tmpl

const DeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}-deployment
  namespace: {{ .Name }}-ns
  labels:
    app: {{ .Name }}-label
  annotations:
    reloader.stakater.com/auto: "true"
spec:
  replicas: 4
  selector:
    matchLabels:
    app: {{ .Name }}-label
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 0
    type: RollingUpdate
  tmpl:
    metadata:
      labels:
        app: {{ .Name }}-label
    spec:
      containers:
        - name: {{ .Name }}-container
          image: {{ .Registry }}/{{ .Name }}:latest
          imagePullPolicy: Always
          env:
            - name: TZ
              value: "Asia/Tokyo"
            - name: SECRET
              valueFrom:
                secretKeyRef:
                  name: {{ .Name }}-secret
                  key: secret
          ports:
            - containerPort: 8080
          resources:
            limits:
                cpu: 2
                memory: 2048m
            requests:
                cpu: 2
            memory: 2048m
          readinessProbe:
            httpGet:
              port: 8080
              path: /health
`

const HpaTemplate  = `apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{ .Name }}-hpa
  namespace: {{ .Name }}-ns
spec:
  maxReplicas: 300
  metrics:
  - resource:
      name: cpu
      targetAverageUtilization: 50
    type: Resource
  minReplicas: 1
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ .Name }}-deployment
`

const LoadBalancerTemplate = `apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-access-log-emit-interval: "5"
    service.beta.kubernetes.io/aws-load-balancer-access-log-enabled: "true"
    service.beta.kubernetes.io/aws-load-balancer-access-log-s3-bucket-name: s3-bucket
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: http
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: arn:aws:acm:ap-northeast-1:xxxxxxxxx:certificate/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "443"
  labels:
    app: {{ .Name }}-label
  name: {{ .Name }}-svc
  namespace: {{ .Name }}-ns
spec:
  ports:
  - name: http
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app: {{ .Name }}-label
  type: LoadBalancer
`

const NodePortTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}-svc
  namespace: {{ .Name }}-ns
  labels:
    app: {{ .Name }}-label
spec:
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
  selector:
    app: {{ .Name }}-label
  type: NodePort
`

const SecretTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}-secret
  namespace: {{ .Name }}-ns
type: Opaque
data:
  {{ .Key }}: {{ .Value }}
`
