# How to run it locally?

1. Start Docker and Kubernetes
```shell
colima start --arch aarch64 --vm-type=vz --vz-rosetta --kubernetes
```

2. Deploy ingress-nginx
```shell
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

3. Deploy echo server
```shell
kubectl apply -f echo-server.yaml
```

4. Deploy ingress
```shell
kubectl apply -f ingress.yaml
```

5. Apply nginx config
```shell
kubectl apply -f nginx-config.yaml
```

6. Modify /etc/hosts
Add `127.0.0.1 echo.test` in `/etc/hosts` 

7. Run reverse-proxy service
```shell
go run main.go
```

8. Verify
* Ingress returns compressed response
```shell
curl -H "Accept-Encoding: gzip" -H "Content-Type: application/json" -d '{"test": "data"}' --compressed  http://echo.test -v
```
The output indicates the response is compressed 
```
* Host echo.test:80 was resolved.
* IPv6: (none)
* IPv4: 127.0.0.1
*   Trying 127.0.0.1:80...
* Connected to echo.test (127.0.0.1) port 80
> POST / HTTP/1.1
> Host: echo.test
> User-Agent: curl/8.7.1
> Accept: */*
> Accept-Encoding: gzip
> Content-Type: application/json
> Content-Length: 16
>
* upload completely sent off: 16 bytes
< HTTP/1.1 200 OK
< Date: Fri, 09 Aug 2024 06:56:48 GMT
< Content-Type: application/json; charset=utf-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< Vary: Accept-Encoding
< ETag: W/"579-VWLE5muZacw1Fl5ebEvtBXDObIs"
< Content-Encoding: gzip
<
* Connection #0 to host echo.test left intact
{"host":{"hostname":"echo.test","ip":"::ffff:10.42.0.9","ips":[]},"http":{"method":"POST","baseUrl":"","originalUrl":"/","protocol":"http"},"request":{"params":{"0":"/"},"query":{},"cookies":{},"body":{"test":"data"},"headers":{"host":"echo.test","x-request-id":"aea7ecca3b13677976c3f25a762a4667","x-real-ip":"10.42.0.2","x-forwarded-for":"10.42.0.2","x-forwarded-host":"echo.test","x-forwarded-port":"80","x-forwarded-proto":"http","x-forwarded-scheme":"http","x-scheme":"http","content-length":"16","user-agent":"curl/8.7.1","accept":"*/*","accept-encoding":"gzip","content-type":"application/json"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"echo-server-5cd559894b-2mf2z","KUBERNETES_SERVICE_HOST":"10.43.0.1","KUBERNETES_PORT_443_TCP_PROTO":"tcp","KUBERNETES_PORT_443_TCP_PORT":"443","ECHO_SERVER_SERVICE_PORT":"80","ECHO_SERVER_PORT_80_TCP_PORT":"80","KUBERNETES_PORT":"tcp://10.43.0.1:443","ECHO_SERVER_PORT":"tcp://10.43.224.119:80","KUBERNETES_SERVICE_PORT":"443","KUBERNETES_SERVICE_PORT_HTTPS":"443","KUBERNETES_PORT_443_TCP":"tcp://10.43.0.1:443","KUBERNETES_PORT_443_TCP_ADDR":"10.43.0.1","ECHO_SERVER_PORT_80_TCP_PROTO":"tcp","ECHO_SERVER_PORT_80_TCP_ADDR":"10.43.224.119","ECHO_SERVER_SERVICE_HOST":"10.43.224.119","ECHO_SERVER_PORT_80_TCP":"tcp://10.43.224.119:80","NODE_VERSION":"20.11.0","YARN_VERSION":"1.22.19","HOME":"/root"}}%
```

* Send a request to reverse-proxy service, reverse-proxy service sends a compressed body to Ingress and then 
Ingress decompresses the request boy and returns compressed response
```shell
curl -X POST -H "Accept-Encoding: gzip" -H "Content-Type: application/json" -d '{"test": "data"}' --compressed http://localhost:8080 -v
```
The output indicates the response is compressed
```
Note: Unnecessary use of -X or --request, POST is already inferred.
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* connect to ::1 port 8080 from ::1 port 52655 failed: Connection refused
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080
> POST / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Accept-Encoding: gzip
> Content-Type: application/json
> Content-Length: 16
>
* upload completely sent off: 16 bytes
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Content-Type: application/json; charset=utf-8
< Date: Fri, 09 Aug 2024 07:02:46 GMT
< Etag: W/"5ba-UQYBV7vi01fN9joKH35+0Ff32mw"
< Vary: Accept-Encoding
< Transfer-Encoding: chunked
<
* Connection #0 to host localhost left intact
{"host":{"hostname":"echo.test","ip":"::ffff:10.42.0.9","ips":[]},"http":{"method":"POST","baseUrl":"","originalUrl":"/","protocol":"http"},"request":{"params":{"0":"/"},"query":{},"cookies":{},"body":{"test":"data"},"headers":{"host":"echo.test","x-request-id":"59fcc2fa32a79ca333034ebc034987b3","x-real-ip":"10.42.0.2","x-forwarded-for":"10.42.0.2","x-forwarded-host":"echo.test","x-forwarded-port":"80","x-forwarded-proto":"http","x-forwarded-scheme":"http","x-scheme":"http","x-original-forwarded-for":"127.0.0.1","content-length":"40","user-agent":"curl/8.7.1","accept":"*/*","accept-encoding":"gzip","content-encoding":"gzip","content-type":"application/json"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"echo-server-5cd559894b-2mf2z","KUBERNETES_SERVICE_HOST":"10.43.0.1","KUBERNETES_PORT_443_TCP_PROTO":"tcp","KUBERNETES_PORT_443_TCP_PORT":"443","ECHO_SERVER_SERVICE_PORT":"80","ECHO_SERVER_PORT_80_TCP_PORT":"80","KUBERNETES_PORT":"tcp://10.43.0.1:443","ECHO_SERVER_PORT":"tcp://10.43.224.119:80","KUBERNETES_SERVICE_PORT":"443","KUBERNETES_SERVICE_PORT_HTTPS":"443","KUBERNETES_PORT_443_TCP":"tcp://10.43.0.1:443","KUBERNETES_PORT_443_TCP_ADDR":"10.43.0.1","ECHO_SERVER_PORT_80_TCP_PROTO":"tcp","ECHO_SERVER_PORT_80_TCP_ADDR":"10.43.224.119","ECHO_SERVER_SERVICE_HOST":"10.43.224.119","ECHO_SERVER_PORT_80_TCP":"tcp://10.43.224.119:80","NODE_VERSION":"20.11.0","YARN_VERSION":"1.22.19","HOME":"/root"}}%
```
