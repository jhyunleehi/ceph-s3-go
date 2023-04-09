# Ceph log

### cephadm 로그

####  실시간 로그 *

```
# ceph -W cephadm
# ceph config set mgr mgr/cephadm/log_to_cluster_level debug
# ceph -W cephadm --watch-debug
# ceph -W cephadm --verbose

# ceph config set mgr mgr/cephadm/log_to_cluster_level info
# ceph log last cephadm

# ceph -w   <<==== 이것으로 모니터링 하면 ERR WRN 이 명시적으로 잘 표시된다. 
```



### Daemon log

https://access.redhat.com/documentation/ko-kr/red_hat_ceph_storage/5/html/administration_guide/ceph-daemons-logs_admin

#### 1. Docker log 설정 상태

```
# docker ps 
# docker inspect  ceph-0824744e-a65e-7c16-8bca-5317bd09676f-mon-c01
                "--default-log-to-file=false",                 # 로그 파일 사용 안함
                "--default-log-to-stderr=true",                # 로그를 표준에러로 출력
                "--default-log-stderr-prefix=debug ",          # log prefix를 debug로 설정
                "--default-mon-cluster-log-to-file=false",     # Cluster 로그도 파일이 아닌 표준에러로 출력
```

#### 2. container 별 log path 

```
root@c01:~# docker inspect  ceph-0824744e-a65e-7c16-8bca-5317bd09676f-mon-c01 --format "{{.LogPath}}"
/var/lib/docker/containers/a611c9bee20256db5a29bdd6c64faaafd1a588de48c17ac62211e5e77842115d/a611c9bee20256db5a29bdd6c64faaafd1a588de48c17ac62211e5e77842115d-json.log
# tail -f  /var/lib/docker/containers/a611c9bee20256db5a29bdd6c64faaafd1a588de48c17ac62211e5e77842115d/a611c9bee20256db5a29bdd6c64faaafd1a588de48c17ac62211e5e77842115d-json.log

```

#### 3.  docker logs 명령을 이용한 log 

```
# docker logs -f --tail 1 ceph-0824744e-a65e-7c16-8bca-5317bd09676f-mon-c01
```



#### 4. journalctl *

* -u 옵션을 사용하여  서비스 유형별 조회
* -f 옵션을 사용하여  로그 모니터링 

```
# journalctl로 데몬별 로그 확인
# journalctl -f -u ceph-5c5a50ae-272a-455d-99e9-32c6a013e694@mon.c01
# journalctl    -u ceph-5c5a50ae-272a-455d-99e9-32c6a013e694@mon.c01

# systemctl  | grep ceph
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@alertmanager.c01.service 
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@crash.c01.service        
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@grafana.c01.service      
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@loki.c01.service         
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@mgr.c01.uuuowu.service   
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@mon.c01.service          
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@node-exporter.c01.service
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@osd.0.service            
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@osd.13.service           
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@osd.19.service           
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@osd.6.service            
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@promtail.c01.service     
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f@rgw.foo.c01.mrawxe.service
  system-ceph\x2d0824744e\x2da65e\x2d7c16\x2d8bca\x2d5317bd09676f.slice                      
  ceph-0824744e-a65e-7c16-8bca-5317bd09676f.target                                         
  ceph.target 
```

####  5. 파일로 보관

```
# ceph config set global log_to_file true
# ceph config set global mon_cluster_log_to_file true
```



### log level 설정

#### 1. Runtime 설정

```
# /var/log/ceph
# ceph daemon {daemon-name} config show | less
# ceph daemon mon.c01  config show | grep log

# ceph tell {daemon-type}.{daemon id or *} config set {name} {value}
# ceph tell osd.0    config set debug_osd 0/5


# ceph tell mon.c01  config set mon_cluster_log_to_file  true

```



* 로그를 여기로 설정한다. 
* "mon_cluster_log_file": "default=/var/log/ceph/ceph.$channel.log cluster=/var/log/ceph/ceph.log",
* "mon_cluster_log_file_level": "debug",
* "mon_cluster_log_to_file": "true",   <<=== 그런데 이것이 docker 안에만 기록이 되는 문제가 있네....

```
root@c01:/# ceph daemon mon.c01 config show | grep mon | grep log
    "clog_to_monitors": "default=true",
    "clog_to_syslog_facility": "default=daemon audit=local0",
    "log_file": "/var/log/ceph/ceph-mon.c01.log",
    "mon_client_log_interval": "1.000000",
    "mon_client_max_log_entries_per_message": "1000",
    "mon_cluster_log_file": "default=/var/log/ceph/ceph.$channel.log cluster=/var/log/ceph/ceph.log",
    "mon_cluster_log_file_level": "debug",
    "mon_cluster_log_to_file": "false",
    "mon_cluster_log_to_graylog": "false",
    "mon_cluster_log_to_graylog_host": "127.0.0.1",
    "mon_cluster_log_to_graylog_port": "12201",
    "mon_cluster_log_to_journald": "false",
    "mon_cluster_log_to_stderr": "true",
    "mon_cluster_log_to_syslog": "default=false",
    "mon_cluster_log_to_syslog_facility": "daemon",
    "mon_cluster_log_to_syslog_level": "info",
    "mon_debug_dump_location": "/var/log/ceph/ceph-mon.c01.tdump",
    "mon_health_detail_to_clog": "true",
    "mon_health_log_update_period": "5",
    "mon_health_to_clog": "true",
    "mon_health_to_clog_interval": "600",
    "mon_health_to_clog_tick_interval": "60.000000",
    "mon_log_full_interval": "50",
    "mon_log_max": "10000",
    "mon_log_max_summary": "50",
    "mon_max_log_entries_per_event": "4096",
    "mon_max_log_epochs": "500",
    "mon_op_log_threshold": "5",
```







#### 2. boot time config

```
  [global]
          debug ms = 1/5

  [mon]
          debug mon = 20
          debug paxos = 1/5
          debug auth = 2

  [osd]
          debug osd = 1/5
          debug filestore = 1/5
          debug journal = 1
          debug monc = 5/20

  [mds]
          debug mds = 1
          debug mds balancer = 1
```











### ceph 배포 환경

```
# Ceph Client 패키지 설치
$ apt install ceph-common
```

#### ceph.conf

```
# cat /etc/ceph/ceph.conf
# minimal ceph.conf for 0824744e-a65e-7c16-8bca-5317bd09676f
[global]
        fsid = 0824744e-a65e-7c16-8bca-5317bd09676f
        mon_host = [v2:192.168.105.51:3300/0,v1:192.168.105.51:6789/0] [v2:192.168.105.52:3300/0,v1:192.168.105.52:6789/0] [v2:192.168.105.53:3300/0,v1:192.168.105.53:6789/0] [v2:192.168.105.54:3300/0,v1:192.168.105.54:6789/0] [v2:192.168.105.55:3300/0,v1:192.168.105.55:6789/0]
[mon.c01]
public network = 192.168.105.0/24
mon_allow_pool_delete = true
```





```
root@c01:/# ceph df
--- RAW STORAGE ---
CLASS     SIZE    AVAIL     USED  RAW USED  %RAW USED
hdd    240 GiB  238 GiB  1.8 GiB   1.8 GiB       0.75
TOTAL  240 GiB  238 GiB  1.8 GiB   1.8 GiB       0.75

--- POOLS ---
POOL                          ID  PGS   STORED  OBJECTS     USED  %USED  MAX AVAIL
.mgr                           1    1  449 KiB        2  1.3 MiB      0     75 GiB
.rgw.root                      2   32  1.5 KiB        4   48 KiB      0     75 GiB
default.rgw.log                6   32  3.6 KiB      209  408 KiB      0     75 GiB
default.rgw.control            7   32      0 B        8      0 B      0     75 GiB
default.rgw.meta               8   32  2.3 KiB       12  120 KiB      0     75 GiB
ec-4k2m-pool                   9   64  201 MiB       38  301 MiB   0.13    150 GiB
default.rgw.buckets.index     10   32      0 B       22      0 B      0     75 GiB
default.rgw.temporary.non-ec  11   32      0 B        0      0 B      0     75 GiB
default.rgw.buckets.non-ec    12   32      0 B        0      0 B      0     75 GiB
default.rgw.buckets.data      13   32   95 MiB       25  286 MiB   0.12     75 GiB
root@c01:/#
root@c01:/#
root@c01:/# ceph osd df
```



```
# ceph config show mon.c01 | grep mon_max
mon_max_pg_per_osd   500   mon

# ceph osd pool set ecpool pg_num 128
```


max total pg > 현재pg 개수 + 추가되는 pg 개수 가 되도록 parameter 설정 

mon_max_pg_per_osd 파라미터를 수정하여 max total pg 수를 늘림