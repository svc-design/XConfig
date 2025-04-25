kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/standard-install.yaml

helm repo add kong https://charts.konghq.com
helm repo update
helm upgrade --install kong kong/ingress -n kong --create-namespace

kubectl patch svc kong-gateway-proxy -n kong \
  --type='merge' \
  -p '{
    "spec": {
      "type": "NodePort",
      "ports": [
        {
          "port": 80,
          "targetPort": 8000,
          "protocol": "TCP",
          "name": "http",
          "nodePort": 80
        },
        {
          "port": 443,
          "targetPort": 8443,
          "protocol": "TCP",
          "name": "https",
          "nodePort": 443
        }
      ]
    }
  }'


 kubectl patch svc kong-gateway-proxy -n kong \
  --type='merge' \
  -p '{
    "spec": {
      "externalIPs": [
        "1.15.155.245"
      ]
    }
  }'


echo "
---
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
 name: kong
 annotations:
   konghq.com/gatewayclass-unmanaged: 'true'

spec:
 controllerName: konghq.com/kic-gateway-controller
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
 name: kong
spec:
 gatewayClassName: kong
 listeners:
 - name: proxy
   port: 80
   protocol: HTTP
   allowedRoutes:
     namespaces:
        from: All
" | kubectl apply -f -


echo "
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
 name: echo
 namespace: ai
 annotations:
   konghq.com/strip-path: 'true'
spec:
 parentRefs:
 - name: kong
 hostnames:
 - 'open-webui.onwalk.net'
 rules:
 - matches:
   - path:
       type: PathPrefix
       value: /
   backendRefs:
   - name: open-webui
     kind: Service
     port: 80
" | kubectl apply -f -

kubectl create secret tls onwalk-tls --cert=/etc/ssl/onwalk.net.pem --key=/etc/ssl/onwalk.net.key

kubectl patch --type=json gateway kong -p='[{
   "op":"add",
   "path":"/spec/listeners/-",
   "value":{
       "name":"proxy-ssl",
       "port":443,
       "protocol":"HTTPS",
       "tls":{
           "certificateRefs":[{
               "group":"",
               "kind":"Secret",
               "name":"onwalk-tls"
           }]
       }
   }
}]'

curl -ksv https://onwalk.net/echo

