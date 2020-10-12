##Issues

### 1. when build docker image `docker build -t server .`

The error message

```
go/src/github.com/coreos/etcd/clientv3/auth.go:125:72: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.AuthEnable
go/src/github.com/coreos/etcd/clientv3/auth.go:130:74: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.AuthDisable
go/src/github.com/coreos/etcd/clientv3/auth.go:135:72: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.AuthStatus
go/src/github.com/coreos/etcd/clientv3/auth.go:140:152: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserAdd
go/src/github.com/coreos/etcd/clientv3/auth.go:145:144: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserAdd
go/src/github.com/coreos/etcd/clientv3/auth.go:150:86: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserDelete
go/src/github.com/coreos/etcd/clientv3/auth.go:155:122: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserChangePassword
go/src/github.com/coreos/etcd/clientv3/auth.go:160:104: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserGrantRole
go/src/github.com/coreos/etcd/clientv3/auth.go:165:80: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserGet
go/src/github.com/coreos/etcd/clientv3/auth.go:170:72: cannot use auth.callOpts (type []"github.com/coreos/etcd/vendor/google.golang.org/grpc".CallOption) as type []"go.etcd.io/etcd/vendor/google.golang.org/grpc".CallOption in argument to auth.remote.UserList
go/src/github.com/coreos/etcd/clientv3/auth.go:170:72: too many errors
```

solved: https://github.com/etcd-io/etcd/pull/10044#issuecomment-417125341


Possible fixes:

- Depend on a released etcd version (git tags like v3.3.9) instead of master since all released versions reference coreos/etcd in imports and will work just fine.
- If one needs to depend on master, explicitly update all import statements to go.etcd.io/etcd and verify there are no direct or transitive dependencies remain for github.com/coreos/etcd.


### 2. Docker For Mac 没有 docker0 网桥

在使用 Docker 时，要注意平台之间实现的差异性，如 Docker For Mac 的实现和标准 Docker 规范有区别，Docker For Mac 的 Docker Daemon 是运行于虚拟机 (xhyve) 中的，而不是像 Linux 上那样作为进程运行于宿主机，因此 Docker For Mac 没有 docker0 网桥，不能实现 host 网络模式，host 模式会使 Container 复用 Daemon 的网络栈 (在 xhyve 虚拟机中)，而不是与 Host 主机网络栈，这样虽然其它容器仍然可通过 xhyve 网络栈进行交互，但却不是用的 Host 上的端口 (在 Host 上无法访问)。bridge 网络模式 -p 参数不受此影响，它能正常打开 Host 上的端口并映射到 Container 的对应 Port。文档在这一点上并没有充分说明，容易踩坑。

docker 18.03 加入了一个 feature，在容器中可以通过`host.docker.internal`来访问主机
> Use your internal IP address or connect to the special DNS name host.docker.internal which will resolve to the internal IP address used by the host.


## Have already implemented features

-[x] grpc
-[x] docker
-[x] service discovery
-[x] service registry
-[ ] load balance
-[ ] service health check
-[x] graceful shutdown
-[ ] grpc connection pool
-[ ] logger
-[ ] swim lane
-[ ] build a framework
-[ ] rage limit
