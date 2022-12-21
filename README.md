# tracert
## Build
```
$ go build
```
## Usage example
```
$ ./tracert ya.ru
tracing route to ya.ru (87.250.250.242), 64 hops max
1 192.168.88.1 [] 1.457139ms
2 217.197.2.1 [] 7.342249ms
3 172.24.31.5 [] 1.467232ms
4 172.24.25.32 [vunk-punk.rtr.pu.ru.] 2.224909ms
5 172.24.25.38 [magma-vunk.rtr.pu.ru.] 11.051899ms
6 195.70.196.3 [vlan3.kronos.pu.ru.] 2.366009ms
7 195.70.206.129 [] 2.408212ms
8 185.1.152.57 [yandex.spb.piter-ix.net.] 3.685523ms
9 87.250.239.183 [sas-32z3-ae1.yndx.net.] 24.901664ms
10 * * *
11 87.250.250.242 [ya.ru.] 16.844785ms
```
