





```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: nginx
  name: nginx-11111111
spec:
  nodeName: k8s-dev
  terminationGracePeriodSeconds: 30
  containers:
  - image: nginx
    imagePullPolicy: IfNotPresent
    name: nginx
    resources:
      limits:
        cpu: "5"
        memory: 500Mi
      requests:
        cpu: "5"
        memory: 200Mi
```







```shell
winkyi@k8s-dev:~$ kubectl get pod
NAME                                         READY   STATUS             RESTARTS           AGE
... ...
nginx-11111111                               0/1     OutOfcpu           0                  2s
... ...
winkyi@k8s-dev:~$ kubectl describe pod nginx-11111111
Name:         nginx-11111111
Namespace:    default
Priority:     0
Node:         k8s-dev/
Start Time:   Tue, 30 Nov 2021 17:43:10 +0800
Labels:       run=nginx
Annotations:  <none>
Status:       Failed
Reason:       OutOfcpu
Message:      Pod Node didn't have enough resource: cpu, requested: 5000, used: 1450, capacity: 4000
IP:           
IPs:          <none>
Containers:
  nginx:
    Image:      nginx
    Port:       <none>
    Host Port:  <none>
    Limits:
      cpu:     5
      memory:  500Mi
    Requests:
      cpu:        5
      memory:     200Mi
    Environment:  <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-dtctg (ro)
Volumes:
  kube-api-access-dtctg:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Burstable
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type     Reason    Age   From     Message
  ----     ------    ----  ----     -------
  Warning  OutOfcpu  21s   kubelet  Node didn't have enough resource: cpu, requested: 5000, used: 1450, capacity: 4000
```

