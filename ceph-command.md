# Ceph

## Command

### 1. 자주 사용하는 명령

* 기본

```log
# ceph orch ls 
# ceph orsch ps 
# ceph health 
# ceph status 
# ceph mon dump
# ceph mgr stat 
# ceph osd pool ls 
# ceph pg stat 
# ceph osd status
# ceph osd tree 
```

* command

```log
# ceph health 
# ceph status 
# cephadm shell -ceph orch ls
# cephadm shell
# ceph config ls | grep  mon_max_pg
# ceph config dump
# ceph config show osd.1
# ceph config show mon.c1

# ceph config get osd public_network
# ceph config get mon public_network
# ceph config set osd.1 debug_ms 10
# ceph config set mon.* mon_cluster_log_to_file true

# ceph tell type.id config set debug_subsystem debug-level
# ceph tell osd.0 config show

# ceph -w

# ceph config set mgr mgr/cephadm/log_to_cluster_level debug
# ceph -W cephadm --watch-debug


# ceph log last cephadm
# cephadm logs --name <name-of-daemon>
# cephadm logs --fsid <fsid> --name <name-of-daemon>
# journalctl  -r -u ceph-da15da2a-cd1c-11ed-9d70-7760e74ff87e@rgw.myorg.myzone.c3.epmnxq.service

# ceph crash -l 
# ceph crash  info  <crash-id>


# cephadm shell
# ceph orch host ls 
# ceph orch host ls --detail 
# ceph orch host add c01

# ceph orch ls
# ceph orch ps 
# ceph orch ps  --daemon_type  osd

# ceph orch device ls
# ceph orch device ls —hostname=servere.lab.example.com
# ceph orch device zap <node> /dev/vda --force
# ceph orch apply osd --all-available-devices
# ceph orch apply osd --all-available-devices --unmanaged=true
# ceph orch daemon add osd c06:/dev/sdc 

# ceph orch apply mon c01,c02,c03,c04,c05,c06
# ceph orch apply rgw myorg ko-mid-1 --placement="2 c1 c2"
# ceph orch daemon add osd node:/dev/vdb
# ceph orch daemon stop osd.12
# ceph orch daemon rm osd.12

# ceph mgr stat 
# ceph mon dump 
# ceph mon stat

# ceph osd status 
# ceph osd tree
# ceph osd pool ls
# ceph osd df 
# ceph osd dump 
# ceph osd rm 12
# ceph osd crush tree 
# ceph osd crush class ls
# ceph osd set noout
# ceph osd set nobackfill
# ceph osd set norecover
# ceph osd set norebalance
# ceph osd set noscrub
# ceph osd set nodeep-scrub
# ceph osd pool create pool-name pg-num pgp-num erasure erasure-code-profile crush-rule-name 

# ceph pg stat 
# ceph pg ls 
# ceph pg dump 
# ceph pg ls-by-pool ec-4k2m-pool
# ceph pg ls-by-pool default.rgw.buckets.data
# ceph pg ls-by-osd osd.1
# ceph pg ls-by-primary osd.1
# ceph pg map 3.1c

# ceph-volume lvm create --bluestore --data /dev/vdc
# ceph-volume lvm prepare --bluestore --data /dev/vdc
# ceph-volume lvm activate <osd-fsid>
# ceph-volume lvm batch --bluestore /dev/vdc /dev/vdd /dev/nvme0n1
# ceph-volume inventory

# ceph auth ls
# ceph auth get client.admin
# ceph auth print-key client.admin
# ceph auth export client.operator1 > ~/operator1.export
# ceph auth import -i ~/operator1.export

# ceph osd getcrushmap -o ./map.bin
# crushtool -d ./map.bin -o ./map.txt



# ceph mgr module disable cephadm
# ceph fsid
# cephadm rm-cluster --force  --fsid  420b849a-cb10-11ed-91b1-5d6354a9bda7
```

### 2. osd

#### 특정 목표 프로비저닝

* osd 생성 : osd 정지 : OSD 데몬을 정지하려면 `ceph orch daemon stop` 명령을 OSD ID와 함께 사용합니다.

```
# ceph orch daemon stop osd.12
```

* osd 데몬 제거 : OSD 데몬을 제거하려면 `ceph orch daemon rm` 명령을 OSD ID와 함께 사용합니다.

```
# ceph orch daemon rm osd.12
```

* osd id 제거  : OSD ID를 해제하려면 `ceph osd rm` 명령을 사용합니다.

```
# ceph osd rm 12
```

#### osd 상태 확인

* 클러스터 상태를 보고 OSD가 실패했는지 확인합니다.

```
  # ceph health detail
```

* 실패한 OSD를 식별합니다.

```
# ceph osd tree | grep -i down
````

* OSD가 실행되는 OSD 노드를 찾습니다.

```
# ceph osd find osd.OSD_ID
```

* 실패한 OSD를 시작해 봅니다.

````
# ceph orch daemon start OSD_ID
````

#### 물리 장치를 교체 절차

* 스크럽을 일시적으로 비활성화합니다.

````

# ceph osd set noscrub ; ceph osd set nodeep-scrub

````

* 클러스터에서 OSD를 제거합니다.

````

# ceph osd out OSD_ID

````

* 클러스터 이벤트를 보고 백필 작업이 시작되었는지 확인합니다.

````

# ceph -w

````

* 백필 프로세스가 모든 PG를 OSD에서 분리하여 이제 안전하게 제거할 수 있는지 확인합니다.

````

# while ! ceph osd safe-to-destroy osd.OSD_ID ; \

  do sleep 10 ; done

````

* OSD를 안전하게 제거할 수 있으면 물리적 스토리지 장치를 교체하고 OSD를 삭제하십시오. 필요한 경우 장치에서 모든 데이터, 파일 시스템, 파티션을 제거합니다.

````

# ceph orch device zap HOST_NAME _OSD_ID --force

````

##### osd 추가

```
# ceph orch device ls
# ceph orch device ls —hostname=servere.lab.example.com
# ceph orch device zap c01 /dev/sdb  --force
# ceph orch apply osd --all-available-devices
# ceph orch apply osd --all-available-devices --unmanaged=true
# ceph orch daemon add osd c06:/dev/sdc
# ceph osd tree
# ceph orch daemon start OSD_ID

```

* osd 추가할때 아래 메세지 나오면 osd 정보가 완전히 정리되지 않아서 그런거다.

```

# ceph orch daemon add osd c03:/dev/sdb
Created no osd(s) on host c03; already created?

```

##### osd stop -> osd  destroy -> osd  rm

```
# cehp osd out osd.1
# ceph osd stop osd.1
# ceph osd rm osd.1
# ceph osd tree  <<=== 찌꺼기 있으면 Crush에서 제거
# cehp osd crush remove osd.1
# ceph osd destroy 1

```

##### osd 삭제

* device만 zap 해서는 확실하게 repository에서 빠지지 않는다.
* daemon도 죽이고 crush에서 bucket 정보도 빼주고 auth도 모두 제거해줘야 다시 붙일 수 있다  (재활용)

```
# ceph osd out osd.${num}
# ceph orch daemon stop osd.${num}
# ceph osd purge ${num} --yes-i-really-mean-it
# ceph orch daemon rm osd.${num} --force
# ceph osd crush remove osd.${num}
# ceph osd auth del osd.${num}
# ceph osd rm ${num}
```

##### osd destroyed, exists

```
# ceph osd destroy 17 --yes-i-really-mean-it
```

```

# ceph osd status
ID  HOST   USED  AVAIL  WR OPS  WR DATA  RD OPS  RD DATA  STATE
19  c01      0      0       0        0       0        0   destroyed,exists  <<---
# ceph orch  device zap  c01 /dev/sdd --force
zap successful for /dev/sdd on c01

# ceph osd status
ID  HOST   USED  AVAIL  WR OPS  WR DATA  RD OPS  RD DATA  STATE
19  c01      0      0       0        0       0        0   exists,new <<-- 다시 올라옴.

```

#### osd pool create

```
# ceph osd lspools
# osd pool default pg num = 100
# osd pool default pgp num = 100
# ceph osd pool create {pool-name} {pg-num} [{pgp-num}] [replicated] [crush-ruleset-name]
# ceph osd pool create {pool-name} {pg-num}  {pgp-num} erasure [erasure-code-profile] [crush-ruleset-name] [expected_num_objects]
# ceph osd pool rename {current-pool-name} {new-pool-name}
# ceph osd pool delete {pool-name} [{pool-name} --yes-i-really-really-mean-it]
# ceph osd pool set-quota data max_objects 10000
# rados df
# ceph osd pool set {poolname} size {num-replicas}
```

### 3. CRUSH

```
# ceph osd crush dump

# ceph osd tree

ID  CLASS  WEIGHT   TYPE NAME      STATUS  REWEIGHT  PRI-AFF
-1         0.02939  root default
-3         0.00980      host c1
 0    hdd  0.00980          osd.0      up   1.00000  1.00000
-5         0.00980      host c2
 1    hdd  0.00980          osd.1      up   1.00000  1.00000
-7         0.00980      host c3
 2    hdd  0.00980          osd.2      up   1.00000  1.00000

# ceph osd crush add-bucket DC1 datacenter
# ceph osd crush add-bucket rackA1 rack
# ceph osd crush add-bucket rackA2 rack
# ceph osd crush move rackA1 datacenter=DC1
# ceph osd crush move rackA2 datacenter=DC1
# ceph osd crush move DC1 root=default
# ceph osd crush move c1 rack=rackA1
# ceph osd crush move c2 rack=rackA1
# ceph osd crush move c3 rack=rackA2

root@c1:~# ceph osd tree
ID   CLASS  WEIGHT   TYPE NAME              STATUS  REWEIGHT  PRI-AFF
 -1         0.02939  root default
 -9         0.02939      datacenter DC1
-10         0.01959          rack rackA1
 -3         0.00980              host c1
  0    hdd  0.00980                  osd.0      up   1.00000  1.00000
 -5         0.00980              host c2
  1    hdd  0.00980                  osd.1      up   1.00000  1.00000
-11         0.00980          rack rackA2
 -7         0.00980              host c3
  2    hdd  0.00980                  osd.2      up   1.00000  1.00000

```

### 4. pool

##### pool 생성

```
# ceph osd lspools
# ceph osd pool ls detail
# ceph osd pool get ec-19k5m-pool all
# ceph osd pool get ec-19k5m-pool all
# ceph osd pool delete ec-4k2m-pool ec-4k2m-pool --yes-i-really-really-mean-it

Error EPERM: pool deletion is disabled; you must first set the mon_allow_pool_delete config option to true before you can destroy a pool

# vi /etc/ceph/ceph.conf

mon_allow_pool_delete = true  << 라인추가

# ceph tell mon.\* injectargs --mon-allow-pool-delete=true
# ceph osd pool delete ec-4k2m-pool ec-4k2m-pool --yes-i-really-really-mean-it
```

##### pool 관리

```
# ceph osd lspools
# ceph osd pool ls detail
# ceph osd pool rename
# ceph osd pool delete
# ceph osd pool set pool_name nodelete true
# ceph osd pool set
# ceph osd pool get
# ceph df
# ceph osd pool stats
# ceph osd pool application enable
# ceph osd pool set-quota
```

##### erasure-code profile, create pool

```
# ceph osd erasure-code-profile ls
# ceph osd erasure-code-profile get default
# ceph osd erasure-code-profile set ecprofile-4k2m k=4 m=2 plugin=jerasure crush-failure-domain=host
# ceph osd erasure-code-profile set ecprofile-19k5m k=19 m=5 plugin=jerasure crush-failure-domain=host
# ceph osd erasure-code-profile ls
# ceph osd erasure-code-profile get ecprofile-19k5m
# ceph osd erasure-code-profile get ecprofile-4k2m
# ceph osd erasure-code-profile get ecprofile-4k2m
# ceph osd pool create ec-4k2m-pool  64 64 erasure ecprofile-4k2m
# ceph osd pool create ec-19k5m-pool 64 64 erasure ecprofile-19k5m
# ceph osd lspools

3 ec-19k5m-pool
4 ec-4k2m-pool

# ceph osd erasure-code-profile rm  ecprofile-4k2m
```

```
# ceph osd erasure-code-profile get <name>  
# ceph osd erasure-code-profile ls
# ceph osd erasure-code-profile rm <name>
# ceph osd erasure-code-profile set <name> [<profile>...] [--force]
```

##### pool 제거

```
# ceph osd pool delete ec-4k2m-pool ec-4k2m-pool --yes-i-really-really-mean-it
Error EPERM: pool deletion is disabled; you must first set the mon_allow_pool_delete config option to true before you can destroy a pool
# vi /etc/ceph/ceph.conf
mon_allow_pool_delete = true  << 라인추가
# ceph tell mon.\* injectargs --mon-allow-pool-delete=true
# ceph osd pool delete ec-4k2m-pool ec-4k2m-pool --yes-i-really-really-mean-it
# ceph osd pool delete  mypool1 mypool1  --yes-i-really-really-mean-it
# ceph osd pool delete  mypool2 mypool2  --yes-i-really-really-mean-it
```

### 5. PG

```
# ceph osd lspools
# ceph pg ls
# ceph pg ls-by-pool ec-4k2m-pool
# ceph pg ls-by-pool default.rgw.buckets.data
# ceph pg ls-by-osd osd.1
# ceph pg ls-by-primary osd.1
# ceph pg map 3.1c

```

* PG: Placement Group의 ID를 나타냅니다.
* OBJECTS: 해당 PG에 속한 오브젝트(object)의 수를 나타냅니다.
* DEGRADED: 해당 PG에서 복제(replication)된 오브젝트 중에서 몇 개가 손상되었는지를 나타냅니다.
* MISPLACED: 해당 PG에서 복제(replication)된 오브젝트 중에서 몇 개가 다른 OSD(Object Storage Daemon)에 저장되어야 하는데 잘못 저장된 상태인지를 나타냅니다.
* UNFOUND: 해당 PG에서 복제(replication)된 오브젝트 중에서 몇 개가 찾을 수 없는 상태인지를 나타냅니다.
* BYTES: 해당 PG에 속한 모든 오브젝트의 크기 합계를 바이트 단위로 나타냅니다.
* OMAP_BYTES*: 해당 PG에 저장된 모든 오브젝트의 omap 데이터 크기를 바이트 단위로 나타냅니다.
* OMAP_KEYS*: 해당 PG에 저장된 모든 오브젝트의 omap 데이터 키(key)의 수를 나타냅니다.
* LOG: 해당 PG의 로그(log) 크기를 나타냅니다.
* STATE: 해당 PG의 상태(state)를 나타냅니다.
* SINCE: 해당 PG의 상태가 언제부터 지속되고 있는지를 나타냅니다.
* VERSION: 해당 PG의 버전(version)을 나타냅니다.
* REPORTED: 해당 PG의 상태가 언제 마지막으로 보고되었는지를 나타냅니다.
* UP: 해당 PG에서 복제(replication)된 오브젝트가 어느 OSD에서 작동하고 있는지를 나타냅니다.
* ACTING: 해당 PG에서 복제(replication)된 오브젝트를 어느 OSD에서 읽고 쓸 수 있는지를 나타냅니다.
* SCRUB_STAMP: 해당 PG의 스크럽(scrub)이 마지막으로 수행된 시간을 나타냅니다.
* DEEP_SCRUB_STAMP: 해당 PG의 딥 스크럽(deep scrub)이 마지막으로 수행된 시간을 나타냅니다.
* AST_SCRUB_DURATION: 해당 PG의 자동 스크럽(scrub) 기간을 나타냅니다. (AST는 Automatic Scrubbing and Repairing Technology의 약어입니다.)
* SCRUB_SCHEDULING: 해당 PG의 스크럽(scrub)이 예정된 시간을 나타냅니다.

#### File -> Object-> Pool -> PG

```

# ceph osd map ec-4k2m-pool n01484850_8845.JPEG

osdmap e234 pool 'ec-4k2m-pool' (12) object 'n01484850_8845.JPEG' -> pg 12.e026cb5f (12.1f) -> up ([4,0,5,3,1,2], p4) acting ([4,0,5,3,1,2], p4)

```

* 특정 pool과 object의 OSD (Object Storage Device) 매핑 정보를 조회
* 데이터의 분산 및 복제 상태를 확인하고 문제가 발생한 OSD를 식별할 수 있다
* "up"은 복제된 데이터의 저장소 위치
* "acting"은 현재 데이터의 저장소 위치
* Crush 맵은 데이터를 저장하기 위해 사용 가능한 OSD들을 나타내는 맵

##### 10MB 파일 f1 업로드

```

# ceph  osd  map  default.rgw.buckets.data  <object-name>

osdmap e979 pool 'default.rgw.buckets.data' (24) object 'f1' -> pg 24.67b909e9 (24.9) -> up ([10,6,3], p10) acting ([10,6,3], p10)

```

* 10MB 파일은 3개 objet로 나눠지고 각각의 PG에 할당되었다.
* 10MB file -> { 24.0, 24.5, 24.6} 3개 pg에 할당되었다. 이 PG는 자신의 OSD할당 모듬을 가지고 있으니 group 이라는 이름을 사용한것 같다.
* pg 앞에 붙는 24는 pool id 이다. 그리고 거기에 하나씩 pg 일련 번호가 붙는다. 32개이다.
* 이렇게 덩어리로 관리한다.
* PG라고 하는 이유가 결국은 그냥

```
# ceph pg ls-by-pool default.rgw.buckets.data
PG     OBJECTS  DEGRADED  MISPLACED  UNFOUND  BYTES    OMAP_BYTES*OMAP_KEYS*  LOG  STATE         SINCE  VERSION  REPORTED  UP             ACTING
24.0         1         0          0        0  4194304            0           0    6  active+clean    28m    955'6    978:61    [4,15,20]p4    [4,15,20]p4  
...
24.5         1         0          0        0  1611392            0           0    6  active+clean    28m    955'6    978:42    [18,4,9]p18    [18,4,9]p18  
24.6         2         0          0        0  4194304            0           0    6  active+clean    28m    955'6    978:51    [5,11,14]p5    [5,11,14]p5  
24.9         0         0          0        0        0            0           0    6  active+clean    28m    955'6    978:36    [10,6,3]p10    [10,6,3]p10  
...
24.1f        0         0          0        0        0            0           0    6  active+clean    28m    955'6    978:36   [16,11,0]p16   [16,11,0]p16

```

##### pg 조회

* pool에서 OSD로 찾아 가는 길인데.. 덩어리로 찾아 간다는 것

```

# ceph pg map 24.0
osdmap e979 pg 24.0 (24.0) -> up [4,15,20] acting [4,15,20]
# ceph pg map 24.5
osdmap e979 pg 24.5 (24.5) -> up [18,4,9] acting [18,4,9]
# ceph pg map 24.6
osdmap e979 pg 24.6 (24.6) -> up [5,11,14] acting [5,11,14]

```

### 6. rgw

#### 1. rgw 서비스 기동

중요한것은 realm, zonegroup, zone 이 생성되고 그것 기반으로 rgw 서비스를 올리는 것이다.
이때 zone와 zonegroup에  placement를 추가하여 pool에 설정한 내용이 등록되도록 정의할 수 있다는 것이다.

* `cephadm` 을 사용하면 Ceph Object Gateway 데몬은 `ceph.conf` 파일 또는 명령줄 옵션 대신 Ceph Monitor 구성 데이터베이스를 사용하여 구성됩니다. 구성이 `client.rgw` 섹션에 없는 경우 Ceph Object Gateway 데몬은 기본 설정으로 시작하고 포트 `80` 에 바인드됩니다.

* 방법1

```
# cephadm shell
# radosgw-admin realm create --rgw-realm=test_realm --default
# radosgw-admin zonegroup create --rgw-zonegroup=default  --master --default
# radosgw-admin zone create --rgw-zonegroup=default --rgw-zone=test_zone --master --default
# radosgw-admin period update --rgw-realm=test_realm --commit
# ceph orch apply rgw myrgw --realm=test_realm --zone=test_zone --placement="2 host01 host02"
```

* 방법2

```
# ceph orch apply rgw <realm> <zone>
```

* 방법3

```
# ceph orch host label add host01 rgw  # the 'rgw' label can be anything
# ceph orch host label add host02 rgw
# ceph orch apply rgw foo '--placement=label:rgw count-per-host:2' --port=8000
# ceph orch host label add ceph-01 rgw
# ceph orch host label add ceph-02 rgw
# ceph orch host label add ceph-03 rgw
```

* realm 생성-> zonegroup-> zone 생성

```
# radosgw-admin realm create --rgw-realm=myorg --default
# radosgw-admin zonegroup create --rgw-zonegroup=myzonegroup --master --default
# radosgw-admin zone create --rgw-zonegroup=myzonegroup --rgw-zone=myzone --master --default
# radosgw-admin period update --rgw-realm=myorg --commit
# radosgw-admin zonegroup placement add  --placement-id myplacement  --rgw-zonegroup myzonegroup
# radosgw-admin zone placement add --placement-id myplacement  --rgw-zonegroup myzonegroup  --rgw-zone myzone  --data-pool  rgw.data --index-pool rgw.index  --data-extra-pool rgw.extra 
# ceph orch apply rgw myorg myzone --placement="2 c2 c3"
# radosgw-admin realm  list
# radosgw-admin zonegroup list
# radosgw-admin zone list
```

#### 2. rgw appplication 제거

```
# ceph orch ps 
# ceph orch daemon stop  rgw.myorg.myzone.c3.qpmfzx
Scheduled to stop rgw.myorg.myzone.c3.qpmfzx on host 'c3'

# ceph orch ls 
# ceph orch rm rgw.myorg.myzone
Removed service rgw.myorg.myzone
```

* 찌꺼기 들 지우기.

```
# radosgw-admin realm list
# radosgw-admin zonegroup list 
# radosgw-admin zone list
# radosgw-admin zonegroup placement list
# radosgw-admin zone placement list

# radosgw-admin realm delete --rgw-realm=myorg
# radosgw-admin zonegroup delete --rgw-zonegroup=myzonegroup
# radosgw-admin zone delete --rgw-zone=myzone
# radosgw-admin zonegroup  placement rm --placement-id="myplacement"
# radosgw-admin zone placement rm --placement-id="myplacement"

# ceph osd pool rm default.rgw.log default.rgw.log --yes-i-really-really-mean-it
# ceph osd pool rm default.rgw.meta default.rgw.meta --yes-i-really-really-mean-it
# ceph osd pool rm default.rgw.control default.rgw.control --yes-i-really-really-mean-it
# ceph osd pool rm default.rgw.data.root default.rgw.data.root --yes-i-really-really-mean-it
# ceph osd pool rm default.rgw.gc default.rgw.gc --yes-i-really-really-mean-it
```

#### 3. rgw 재동

```
# ceph orch ps  ==> daemon 이름 확인
# ceph orch daemon restart rgw.foo.c01.mrawxe
# ceph orch daemon restart rgw.foo.c02.udcxyw
Scheduled to restart rgw.foo.c01.mrawxe on host 'c01'
```

#### 4. placement 추가

* <주의> zone을 사용하는 경우 해당 zonename.rgw.* 이렇게 pool 이름을 맞춰 줘야 한다.
* pool 생성

 ```
# cephadm shell
# radosgw-admin realm create --rgw-realm=myorg --default
# radosgw-admin zonegroup create --rgw-zonegroup=myzonegroup --master --default
# radosgw-admin zone create --rgw-zonegroup=myzonegroup --rgw-zone=myzone --master --default
# ceph orch apply rgw myrgw --realm=myorg --zone=myzone --placement="2 host01 host02"

# ceph orch apply rgw  myorg myzone
==> 이 시점에 자동으로 zone의 default pool 생성된다. 
==> 이 pool이 마음에 안들면 placement 생성해서 교체하면된다.  
==> 먼저 pool을 만들고 rgw를 생성하면 rgw가 자신이름에 맞는 pool을 선택해서 사용한다.  
# ceph osd pool create myzone.rgw.buckets.data 4 4 replicated
# ceph osd pool create myzone.rgw.buckets.index 4 4 replicated
# ceph osd pool create myzone.rgw.buckets.non-ec 4 4 replicated

# ceph osd pool application enable myzone.rgw.buckets.data rgw
# ceph osd pool application enable myzone.rgw.buckets.index rgw
# ceph osd pool application enable myzone.rgw.buckets.non-ec rgw

==>> 여기까지 오면 default-placement가 하나 생성된다. 
# radosgw-admin zonegroup placement add --placement-id myplacement1 --rgw-zonegroup myzonegroup 
# radosgw-admin zone placement add --placement-id myplacement1 --rgw-zonegroup myzonegroup --rgw-zone myzone \
   --data-pool  myzone.rgw.data --index-pool myzone.rgw.index --data-extra-pool myzone.extra
# radosgw-admin period update --rgw-realm=<realm-name> --commit
==> rgw  재기동    
# ceph orch daemon restart rgw.myorg.myzone.c1.aerixz
# ceph orch daemon restart rgw.myorg.myzone.c3.epmnxq
```

#### 5. placemnet 등록

[placement 참고](https://docs.ceph.com/en/latest/radosgw/placement/)

* radosgw-admin zonegroup get

```
# radosgw-admin zonegroup get
{
    "id": "0842a525-093d-4dd1-aba0-07070e9c2cc8",    
    "master_zone": "baa66b66-4787-483b-8041-24e4714eed37",   <<-----
    "zones": [
        {
            "id": "baa66b66-4787-483b-8041-24e4714eed37",
            "name": "myzone",            
            "redirect_zone": ""
        }
    ],
    "placement_targets": [
        {
            "name": "default-placement",
            "tags": [],
            "storage_classes": [
                "STANDARD"
            ]
        },
        {
            "name": "myplacement",
            "tags": [],
            "storage_classes": [
                "STANDARD"
            ]
        }
    ],
    "default_placement": "default-placement",   <<-----
    "realm_id": "8f2d61c5-194b-4b7e-b8c5-c18689dabf0c",  <<---
}
```

```
# radosgw-admin zone get
{
    "id": "baa66b66-4787-483b-8041-24e4714eed37",
    "name": "myzone",   <<<----
    "domain_root": "myzone.rgw.meta:root", <--- 자동생성
    "control_pool": "myzone.rgw.control", <--- 자동생성
    "gc_pool": "myzone.rgw.log:gc",<--- 자동생성
    "lc_pool": "myzone.rgw.log:lc", <--- 자동생성
    "log_pool": "myzone.rgw.log", <--- 자동생성
    "intent_log_pool": "myzone.rgw.log:intent",<--- 자동생성
    "usage_log_pool": "myzone.rgw.log:usage",  <--- 자동생성
    "roles_pool": "myzone.rgw.meta:roles", <--- 자동생성
    "reshard_pool": "myzone.rgw.log:reshard", <--- 자동생성
    "user_keys_pool": "myzone.rgw.meta:users.keys", <--- 자동생성
    "user_email_pool": "myzone.rgw.meta:users.email", 
    "user_swift_pool": "myzone.rgw.meta:users.swift",
    "user_uid_pool": "myzone.rgw.meta:users.uid",
    "otp_pool": "myzone.rgw.otp",
    "system_key": {
        "access_key": "",
        "secret_key": ""
    },
    "placement_pools": [
        {
            "key": "default-placement",
            "val": {
                "index_pool": "myzone.rgw.buckets.index",
                "storage_classes": {
                    "STANDARD": {
                        "data_pool": "myzone.rgw.buckets.data"
                    }
                },
                "data_extra_pool": "myzone.rgw.buckets.non-ec",
                "index_type": 0
            }
        },
        {
            "key": "myplacement",
            "val": {
                "index_pool": "myzone.rgw.log",
                "storage_classes": {
                    "STANDARD": {
                        "data_pool": "myzone.rgw.meta"
                    }
                },
                "data_extra_pool": "myzone.control",
                "index_type": 0
            }
        }
    ],
    "realm_id": "8f2d61c5-194b-4b7e-b8c5-c18689dabf0c"
}
```

* 새로운 placement 추가

```
$ radosgw-admin zonegroup placement add \
      --rgw-zonegroup default \
      --placement-id temporary

$ radosgw-admin zone placement add \
      --rgw-zone default \
      --placement-id temporary \
      --data-pool default.rgw.temporary.data \
      --index-pool default.rgw.temporary.index \
      --data-extra-pool default.rgw.temporary.non-ec

```

* glacier placement 예시

```
# cephadm shell
# radosgw-admin realm create --rgw-realm=myorg --default
# radosgw-admin zonegroup create --rgw-zonegroup=myzonegroup  --master --default
# radosgw-admin zone create --rgw-zonegroup=myzonegroup --rgw-zone=myzone --master --default
# ceph orch apply rgw myrgw --realm=myorg --zone=myzone --placement="2 c1 c2"

# ceph osd pool create myzone.rgw.glacier.data 
# ceph osd pool create myzone.rgw.glacier.index
# ceph osd pool application enable myzone.rgw.glacier.data rgw
# ceph osd pool application enable myzone.rgw.glacier.index rgw

# radosgw-admin zonegroup get
# radosgw-admin zone get

# radosgw-admin zonegroup placement add --rgw-zonegroup myzonegroup --placement-id myplacement-glacier       
# radosgw-admin zone placement add --placement-id myplacement-glacier --rgw-zone myzone  --data-pool  myzone.rgw.glacier.data --index-pool myzone.rgw.glacier.index 

# radosgw-admin period update --rgw-realm=myorg --commit
# ceph orch daemon restart rgw.myorg.myzone.c1.aerixz
# ceph orch daemon restart rgw.myorg.myzone.c3.epmnxq  
```

#### 6. multi-zone, DR

* master cluster

```
# radosgw-admin realm create --rgw-realm=cl260 --default
# radosgw-admin zonegroup create --rgw-zonegroup=classroom --endpoints=http://serverc:80 --master --default
# radosgw-admin zone create  --rgw-zonegroup=classroom --rgw-zone=main --endpoints=http://serverc:80 --access-key=replication --secret=secret --master --default
# radosgw-admin user create --uid="repl.user" --display-name="Replication User" --secret=secret --system --access-key=replication
# radosgw-admin period update --commit
# ceph orch apply rgw cl260-1 --realm=cl260 --zone=main --placement="1 serverc.lab.example.com"
# ceph config set client.rgw rgw_zone main
# ceph config get client.rgw rgw_zone
# ceph config set client.rgw rgw_dynamic_resharding false
# ceph config get client.rgw rgw_dynamic_resharding
# radosgw-admin user create --display-name="S3 user" --uid="apiuser" --access="full" --access_key="review" --secret="securekey"
```

* slave cluster

```
# radosgw-admin realm pull --url=http://serverc:80 --access-key=replication --secret-key=secret
# radosgw-admin period pull --url=http://serverc:80 --access-key=replication --secret-key=secret
# radosgw-admin realm default --rgw-realm=cl260
# radosgw-admin zonegroup default  --rgw-zonegroup=classroom
# radosgw-admin zone create --rgw-zonegroup=classroom --rgw-zone=fallback --endpoints=http://serverf:80 --access-key=replication --secret-key=secret --default
# radosgw-admin period update --commit --rgw-zone=fallback
# ceph orch apply rgw cl260-2 --zone=fallback --placement="1 serverf.lab.example.com" --realm=cl260
# ceph config set client.rgw rgw_zone fallback
# ceph config get client.rgw rgw_zone
# ceph config set client.rgw rgw_dynamic_resharding false
# ceph config get client.rgw rgw_dynamic_resharding
```

### 7. S3 API

#### 1. user 생성

* 키 지정해서 만들면 지정한 키로 만들어 지고, 지정하지 않으면 자동 생성된다.
* --gen-access-key : 동일 user에게서 복수개의 access key 생성한다.
* --gen-secret :  security key 만  재 생성한다.
* key rm 명령을 이용해서 key를 제거한다.
* --purge-data : 사용자 제거와 bucket 제거 포함.

```
# radosgw-admin user create --uid=testuser --display-name="Test User" --email=test@example.com --access-key=12345 --secret=67890
# radosgw-admin key create --uid=s3user --access-key="8PI2D9ARWNGJI99K8TOS" --gen-secret
# radosgw-admin key rm --uid=s3user  --access-key=8PI2D9ARWNGJI99K8TOS
```

#### 2. quota 지정

```
# radosgw-admin quota set --quota-scope=user --uid=app1  --max-objects=1024
[ceph: root@node /]# radosgw-admin quota enable --quota-scope=user --uid=app1
```

```
# radosgw-admin user info --uid=*`uid`*` 
# radosgw-admin user stats --uid=*`uid`*`
```

#### 3.bucket

```
# radosgw-admin bucket list
# radosgw-admin bucket rm 
```

#### 4. aws 설치

```log
# curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
# unzip awscliv2.zip
# sudo ./aws/install
# cat ~/.aws/credentials
[default]
aws_access_key_id = 2VEVI5DW1GOO2QUX14BL
aws_secret_access_key = GefWYyvAlQhszFNZEwDjeDyMa8fS0vCUJ5jSvUgP

# aws s3 mb s3://demobucket

# dd if=/dev/urandom of=file1 bs=1MB count=10
# aws s3 cp ./file1 s3://mybucket/file1 --endpoint-url http://192.168.57.11:80
# aws s3 sync ./test_data/ s3://mybucket/test_data/ --endpoint-url http://192.168.57.11:80
# aws s3 rm  s3://mybucket/f1  --endpoint-url http://192.168.57.11:80
# aws s3 ls mybucket --endpoint-url http://192.168.57.11:80
# aws s3 cp s3://mybucket/file1 ~/f2 --endpoint-url http://192.168.57.11:80
```

#### 5. ceph S3 api 사용방법

[https://github.com/jhyunleehi/ceph-s3-go](https://github.com/jhyunleehi/ceph-s3-go)

### 8. rados

```
# rados lspools
# rados df
```

##### get

```
# rados --pool default.rgw.buckets.data ls
>> 조회된 결과 중에서 선택
# rados -p default.rgw.buckets.data get edb79dac-6966-41fd-8a17-84f2c06323d8.34467.1_test_data/n01484850/n01484850_8845.JPEG n01484850_8845.JPEG
```

* pool name : default.rgw.buckets.data
* 특정 object 경로 : test_data/n01484850/
* 미리 생성해 둔 local 폴더 이름 : n01484850

```
# for obj in $(rados -p default.rgw.buckets.data ls | grep test_data/n01484850/)
do
rados -p default.rgw.buckets.data get $obj ./n01484850/$(basename $obj)
done
```

##### put

```
# rados -p ec-4k2m-pool-``01` `put n01484850_8845.JPEG /home/testdata/n01484850_8845.JPEG
```

##### delete

```
# rados -p ec19k5m-test-01-pool rm n01484850_8845.JPEG
```

### 9. rbd

##### 1. ceph

```
# rbd create mypool/myimage --size 102400
# rbd create mypool/myimage --size 102400 --object-size 8M
# rbd rm mypool/myimage
# rbd map mypool/myimage --id admin --keyfile secretfile
# rbd unmap /dev/rbd0
```

##### 2. client

```
# rbd ls -l
# rbd  map vol01 --pool <pool-name>
# rbd  map vol02 --pool <pool-name>
# rbd  map vol03 --pool <pool-name>
# lsblk
```

```
[client]
# yum -y install ceph-fuse
# ceph-authtool -p /etc/ceph/ceph.client.admin.keyring > admin.key
# cat admin.key
# mount -t ceph h6:6789:/  /cephfile  -o name=admin,secretfile=/root/admin.key
```

CRUSH – Contraolled, Scalable, Decentralized Placement of Replicated Data

### 10.Ceph 명령참고

* [ceph-volume – Ceph OSD deployment and inspection tool](https://docs.ceph.com/en/quincy/man/8/ceph-volume/)
* [ceph-volume-systemd – systemd ceph-volume helper tool](https://docs.ceph.com/en/quincy/man/8/ceph-volume-systemd/)
* [ceph – ceph administration tool](https://docs.ceph.com/en/quincy/man/8/ceph/)
* [ceph-authtool – ceph keyring manipulation tool](https://docs.ceph.com/en/quincy/man/8/ceph-authtool/)
* [ceph-clsinfo – show class object information](https://docs.ceph.com/en/quincy/man/8/ceph-clsinfo/)
* [ceph-conf – ceph conf file tool](https://docs.ceph.com/en/quincy/man/8/ceph-conf/)
* [ceph-debugpack – ceph debug packer utility](https://docs.ceph.com/en/quincy/man/8/ceph-debugpack/)
* [ceph-dencoder – ceph encoder/decoder utility](https://docs.ceph.com/en/quincy/man/8/ceph-dencoder/)
* [ceph-mon – ceph monitor daemon](https://docs.ceph.com/en/quincy/man/8/ceph-mon/)
* [ceph-osd – ceph object storage daemon](https://docs.ceph.com/en/quincy/man/8/ceph-osd/)
* [ceph-kvstore-tool – ceph kvstore manipulation tool](https://docs.ceph.com/en/quincy/man/8/ceph-kvstore-tool/)
* [ceph-run – restart daemon on core dump](https://docs.ceph.com/en/quincy/man/8/ceph-run/)
* [ceph-syn – ceph synthetic workload generator](https://docs.ceph.com/en/quincy/man/8/ceph-syn/)
* [crushdiff – ceph crush map test tool](https://docs.ceph.com/en/quincy/man/8/crushdiff/)
* [crushtool – CRUSH map manipulation tool](https://docs.ceph.com/en/quincy/man/8/crushtool/)
* [librados-config – display information about librados](https://docs.ceph.com/en/quincy/man/8/librados-config/)
* [monmaptool – ceph monitor cluster map manipulation tool](https://docs.ceph.com/en/quincy/man/8/monmaptool/)
* [osdmaptool – ceph osd cluster map manipulation tool](https://docs.ceph.com/en/quincy/man/8/osdmaptool/)
* [rados – rados object storage utility](https://docs.ceph.com/en/quincy/man/8/rados/)

## Ceph install : Ceph Guide

#### 1. host setup

* /etc/host
* ip 설정
* hostname
* ntp 설정

```
# cat /etc/hosts
192.168.57.11 c1
192.168.57.12 c2
192.168.57.13 c3

# hostnamectl set-hostname c1
# root@c1:~# cat /etc/netplan/00-installer-config.yaml
# This is the network config written by 'subiquity'
network:
  ethernets:
    enp0s3:
      dhcp4: true
    enp0s8:
      addresses:
      - 192.168.57.11/24
      gateway4: 192.168.57.1
      nameservers:
        addresses: []
        search: []
  version: 2
root@c1:~# netplan apply

```

* docker 설치, package 설치

```
# sudo apt-get install ca-certificates curl gnupg lsb-release
# curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
#  echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# cat /etc/apt/sources.list.d/docker.list
deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu   focal stable

# apt update 
# apt-get install -y docker-ce  docker-ce-cli containerd.io

# cp deploy/CARL.crt /usr/local/share/ca-certificates/
# update-ca-certificates
# vi /etc/apt/sources.list
# apt install -y ntp
# vi /etc/ntp.conf
# apt install -y make
```

#### 2. ceph bootstrap 설치

```
# curl --silent --remote-name --location https://github.com/ceph/ceph/raw/octopus/src/cephadm/cephadm

# chmod +x cephadm
# ./cephadm add-repo --release octopus
# ./cephadm install
# which cephadm
/usr/sbin/cephadm

# which cephadm
/usr/sbin/cephadm

# mkdir -p /etc/ceph
# cephadm bootstrap --mon-ip  192.168.57.11
Ceph Dashboard is now available at:

URL: https://c1:8443/
User: admin
Password: 9hpblkbcnu
sudo /usr/sbin/cephadm shell --fsid 420b849a-cb10-11ed-91b1-5d6354a9bda7 -c /etc/ceph/ceph.conf -k /etc/ceph/ceph.client.admin.keyring

# cephadm shell
# cephadm add-repo --release octopus
# cephadm install ceph-common

# ceph  -v
ceph version 15.2.17 (8a82819d84cf884bd39c17e3236e0632ac146dc4) octopus (stable)

```

#### 3. add host to cluster

```
# ceph cephadm get-pub-key > ~/ceph.pub
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c1
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c2
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c3
```

```
# cephadm shell
# ceph orch host add c2
# ceph orch host add c3


# ceph config set mon public_network 192.168.57./24
# ceph  orch apply mon 3
# ceph orch host label add c1 _admin
# ceph orch apply mon *<host1,host2,host3,...>*
# ceph orch host ls
# ceph orch apply mon label:mon
```

* yaml로 설치하기.

```
# ceph orch apply -i file.yaml

service_type: mon
placement:
  hosts:
   - host1
   - host2
   - host3
```

* osd 정리

```
# ceph orch device zap c1 /dev/sdb --force
# ceph orch device zap c2 /dev/sdb --force
# ceph orch device zap c3 /dev/sdb --force
# ceph orch apply osd --all-available-devices
# ceph orch daemon add osd host1:/dev/sdb
```

* MDS

```
# ceph orch apply mds *<fs-name>* --placement="*<num-daemons>* [*<host1>* ...]"
```

* RGW : realm, zone 설치

```
# ceph orch apply rgw *<realm-name>* *<zone-name>* --placement="*<num-daemons>* [*<host1>* ...]"
# ceph orch apply rgw myorg us-east-1 --placement="2 myhost1 myhost2"
```

```
# radosgw-admin realm create --rgw-realm=myorg --default
# radosgw-admin zonegroup create --rgw-zonegroup=myzonegroup --master --default
# radosgw-admin zone create --rgw-zonegroup=myzonegroup --rgw-zone=myzone --master --default
# radosgw-admin period update --rgw-realm=myorg --commit

# ceph orch apply rgw myorg myzone --placement="2 c1 c2"

# radosgw-admin zone list 
```

* user add : dashboard

```
# radosgw-admin user create --uid=dashboard --display-name=dashboard --system
# radosgw-admin user info --uid=dashboard
```

* RGW 화면이 안 보일 때 : 뭔가 user 를 생성해줘야 하는 듯

```
# radosgw-admin user create --uid=dashboard --display-name=dashboard --system
root@c1:/# radosgw-admin user info --uid=dashboard
{
    "user_id": "dashboard",
    "display_name": "dashboard",
     "keys": [
        {
            "user": "dashboard",
            "access_key": "B1S1DYWLG52YNPXFJ7F0",
            "secret_key": "noi7877VTggsLmNDqQj0uCpfeaUcH4UTW93j6t7u"
        }
    ],
}

# echo  "B1S1DYWLG52YNPXFJ7F0" > akey
# echo  "noi7877VTggsLmNDqQj0uCpfeaUcH4UTW93j6t7u" > skey
# ceph dashboard set-rgw-api-access-key -i akey
# ceph dashboard set-rgw-api-secret-key -i skey
```

#### 4. fsid 여러개 일때  rm-cluster

```
c# cephadm shell  -- ceph -v
ERROR: Cannot infer an fsid, one must be specified (using --fsid): ['0824744e-a65e-7c16-8bca-5317bd09676f', '1373ff90-46ab-678d-2467-8030180fe5c0']

# cephadm ls

# cephadm rm-cluster --fsid 1373ff90-46ab-678d-2467-8030180fe5c0 --force
```

#### 5. 설치 결과

```
# ceph orch ls
# ceph -s
# ceph osd unset noout
# osd unset norecover
# ceph osd unset norebalance
# ceph osd unset nobackfill
# ceph osd unset nodown
# ceph osd unset pause
```

#### 6. 모니터링, log 보기

* docker 보자

```
# podman exec -it ceph-mon-MONITOR_NAME /bin/bash
# podman exec -it ceph-499829b4-832f-11eb-8d6d-001a4a000635-mon.host01 /bin/bash
# docker inspect  ceph-0824744e-a65e-7c16-8bca-5317bd09676f-mon-c01
# docker logs -f -tail 1 ceph-0824744e-a65e-7c16-8bca-5317bd09676f-mon-c01
```

* cephadmin

```
# cephadm shell
# ceph health
# ceph status 
# ceph config set mgr mgr/cephadm/log_to_cluster_level info
# ceph log last cephadm
# ceph -w

# systemctl | grep ceph
# journalctl -f -u ceph-5c5a50ae-272a-455d-99e9-32c6a013e694@mon.c01
# journalctl -r -u ceph-5c5a50ae-272a-455d-99e9-32c6a013e694@mon.c01

```

* log 설정

```
# ceph config set global log_to_file true
# ceph config set global mon_cluster_log_to_file true
```

## Cluster 제거 후 재설치

#### 1. cephadm 정지

```
# ceph mgr module disable cephadm
# ceph fsid
# cephadm rm-cluster --force  --fsid  420b849a-cb10-11ed-91b1-5d6354a9bda7
```

#### 2. 노드별  정리 reset

```sh
#!/bin/bash
systemctl | grep 'ceph-' | awk '{print $2}' | grep -v loaded | xargs systemctl stop
systemctl | grep 'ceph-' | awk '{print $2}' | grep -v loaded | xargs systemctl disable 
systemctl | grep 'system-ceph' | awk '{print $1}' | xargs systemctl stop
systemctl | grep 'system-ceph' | awk '{print $1}' | xargs systemctl disable 
systemctl | grep 'ceph' | awk '{print $2}' | grep -v loaded | xargs systemctl stop
systemctl | grep 'ceph' | awk '{print $2}' | grep -v loaded | xargs systemctl disable 
rm -rf  /etc/systemd/system/ceph-*
rm -rf /etc/systemd/system/ceph.*
systemctl daemon-reload
systemctl reset-failed
rm -rf  /etc/systemd/system/ceph-*
rm -rf /etc/systemd/system/ceph.*
systemctl daemon-reload
systemctl reset-failed
```

#### 3. config-cluster.yaml

```yaml
---
service_type: host 
addr: 192.168.57.11
hostname: c1
---
service_type: host
addr: 192.168.57.12
hostname: c2
---
service_type: host
addr: 192.168.57.13
hostname: c4
---
service_type: mon 
placement:
  hosts:
    - c1
    - c2
    - c3
---
service_type: rgw 
service_id: realm.zone
placement:
  hosts:
    - c2
    - c3
---
service_type: mgr 
placement:
  hosts:
    - c1
    - c2
---
service_type: osd 
service_id: default_drive_group
placement: 
  host_pattern: 'c*'
data_devices:
  paths:
    - /dev/sdb
    - /dev/sdc
```

#### 4. boot strap

```sh
# cephadm bootstrap --mon-ip=192.168.57.11
--apply-spec=config-cluster.yaml \
--initial-dashboard-password=1234qwer \
--dashboard-password-noupdate 
```

* 생성 결과

```sh
Ceph Dashboard is now available at:
             URL: https://c1:8443/
            User: admin
        Password: 3uoxtbks79

You can access the Ceph CLI with:
        sudo /usr/sbin/cephadm shell --fsid da15da2a-cd1c-11ed-9d70-7760e74ff87e -c /etc/ceph/ceph.conf -k /etc/ceph/ceph.client.admin.keyring

Please consider enabling telemetry to help improve Ceph:
        ceph telemetry on

For more information see:
        https://docs.ceph.com/docs/master/mgr/telemetry/

Bootstrap complete.

```

* telemetry on

```
#  ceph telemetry on --license sharing-1-0
```

#### 5.  add host to cluster

```
# ceph cephadm get-pub-key > ~/ceph.pub
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c1
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c2
# ssh-copy-id -f -i /etc/ceph/ceph.pub root@c3
```

```
# cephadm shell
# ceph orch host add c2
# ceph orch host add c3

# ceph orch apply osd --all-available-devices
# for i in {1..3}; do ceph orch device zap c${i} /dev/sdb  --force ; done
# for i in {1..3}; do ceph orch device zap c${i} /dev/sdc  --force ; done

# ceph config set mon public_network 192.168.57./24
# ceph orch apply mon 3
# ceph orch host label add c1 _admin
# ceph orch apply mon *<host1,host2,host3,...>*
# ceph orch host ls
# ceph orch apply mon label:mon
```

#### 6. rgw 설치

* rgw 찌꺼기 정리

```
# radosgw-admin realm list
# radosgw-admin zonegroup list
# radosgw-admin zone list
# radosgw-admin zonegroup placement list
# radosgw-admin zone placement list
==> 청소한다
# radosgw-admin realm delete --rgw-realm="fooradosgw-admin"
# radosgw-admin realm delete --rgw-realm="defult"
# radosgw-admin zone delete --rgw-zone="defult"
# radosgw-admin zone delete --rgw-zone=myzone2
# radosgw-admin zone delete --rgw-zone=myzone3
# radosgw-admin zonegroup placement  rm --placement-id="myplacement"
# radosgw-admin zone placement  rm --placement-id="myplacement"
```

* realm 에서 "--default" 옵션을 주면 전체 클러스터에서 realm 1개만 만든다는 의미임. 만약 realm 을 2개 이상 만든마면    --rgw-real= 이름을 명시적으로 사용해야 한다. 한개만 만든다면 그런것을 지정할  필요는 없다.  
* 중요한것은 realm, zonegroup, zone 이 생성되고 그것 기반으로 rgw 서비스를 올리는 것이다.
* 이때 zone와 zonegroup에  placement를 추가하여 pool에 설정한 내용이 등록되도록 정의할 수 있다는 것이다.
* zone 생성할 때 --master 누락하면 안된다.  

```
# radosgw-admin realm create --rgw-realm=myorg --default
# radosgw-admin zonegroup create --rgw-zonegroup=myzonegroup --master --default
# radosgw-admin zone create --rgw-zonegroup=myzonegroup --rgw-zone=myzone  --master --default
# radosgw-admin period update --rgw-realm=myorg --commit

==>생성한 realm, zone 기반으로 서비스를 기동한다.  
# ceph orch apply rgw  myorg myzone


==>POOL들이 자동생성됨  myzone.rgw.log myzone.rgw.meta, myzone.conrol
==>아마도 placement가 없어서 자동으로 생성해주는 것인가?
```

* dashboard user 생성

```
# radosgw-admin user create --uid=dashboard --display-name=dashboard --system
# radosgw-admin user info --uid=dashboard
{
    "user_id": "dashboard",
    "display_name": "dashboard",
     "keys": [
        {
            "user": "dashboard",
            "access_key": "B1S1DYWLG52YNPXFJ7F0",
            "secret_key": "noi7877VTggsLmNDqQj0uCpfeaUcH4UTW93j6t7u"
        }
    ],
}

# echo  "B1S1DYWLG52YNPXFJ7F0" > akey
# echo  "noi7877VTggsLmNDqQj0uCpfeaUcH4UTW93j6t7u" > skey
# ceph dashboard set-rgw-api-access-key -i akey
# ceph dashboard set-rgw-api-secret-key -i skey
```

* user 생성

```
# radosgw-admin user create  --uid myuser1  --display-name myuser1
# radosgw-admin user create  --uid myuser2  --display-name myuser2
```

#### 7.placemnet 등록

[placement 참고](https://docs.ceph.com/en/latest/radosgw/placement/)

```
# ceph osd pool create myzone.rgw.glacier.data 
# ceph osd pool create myzone.rgw.glacier.index
# ceph osd pool application enable myzone.rgw.glacier.data rgw
# ceph osd pool application enable myzone.rgw.glacier.index rgw

# radosgw-admin zonegroup get
# radosgw-admin zone get
$ radosgw-admin zonegroup placement add \
      --rgw-zonegroup myzonegroup \
      --placement-id myplacement-glacier       

$ radosgw-admin zone placement add \
      --rgw-zone myzone \
      --placement-id myplacement-glacier \
      --data-pool  myzone.rgw.glacier.data \
      --index-pool myzone.rgw.glacier.index 
      

# ceph orch daemon restart rgw.myorg.myzone.c1.aerixz
# ceph orch daemon restart rgw.myorg.myzone.c3.epmnxq      
```

## AWS Client 설치

### Clinet 설치

##### 1. S3 client 설치

```
# curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
# unzip awscliv2.zip
# sudo ./aws/install
```

##### 2.  credentials

```
# cat ~/.aws/credentials
[default]
aws_access_key_id = 2VEVI5DW1GOO2QUX14BL
aws_secret_access_key = GefWYyvAlQhszFNZEwDjeDyMa8fS0vCUJ5jSvUgP
```

##### 3. upload file

```
# dd if=/dev/urandom of=file1 bs=1MB count=10
# aws s3 cp ./file1 s3://mybucket/file1 --endpoint-url http://192.168.105.51:80
upload: ./file1 to s3://mybucket/file1
```

* 하위 폴더 포함해서 many file

```
# aws s3 sync ./test_data/ s3://mybucket/test_data/ --endpoint-url http://192.168.105.11:80
```

* 삭제

```
# aws s3  rm    s3://mybucket/f1  --endpoint-url http://192.168.105.51:80
delete: s3://mybucket/f1
```

##### 4. get list

```
# aws s3  ls mybucket --endpoint-url http://192.168.105.51:80
2023-03-19 05:32:00   10000000 file1
```

##### 5. download

```
# aws s3 cp s3://mybucket/file1 ~/f2 --endpoint-url http://192.168.105.51:80
```

##### 6. pg map

```
# ceph  osd  map  default.rgw.buckets.data  f1
osdmap e979 pool 'default.rgw.buckets.data' (24) object 'f1' -> pg 24.67b909e9 (24.9) -> up ([10,6,3], p10) acting ([10,6,3], p10)
```

* pg 목록

```
# ceph pg ls-by-pool default.rgw.buckets.data
PG    OBJECTS  DEGRADED  MISPLACED  UNFOUND  BYTES   LOG  STATE         SINCE  VERSION  REPORTED  UP         
24.0  1  0   0  0  4194304  6  active+clean    28m    955'6    978:61    [4,15,20]p4    [4,15,20]p4  
...
24.5  1  0   0  0  1611392   6  active+clean    28m    955'6    978:42    [18,4,9]p18    [18,4,9]p18  
24.6  2  0   0  0  4194304   6  active+clean    28m    955'6    978:51    [5,11,14]p5    [5,11,14]p5  
24.9  0  0   0  0        0   6  active+clean    28m    955'6    978:36    [10,6,3]p10    [10,6,3]p10  
...
24.1f 0  0   0  0        0   6  active+clean    28m    955'6    978:36   [16,11,0]p16   [16,11,0]p16 
```

##### 7. 저장 용량 분석

* 10MB 파일 1개가 4개 object 되고 이것을 PG에서 14MB되고 이것이 pool 기준으로는 29MiB로 표시되는 군

```
# rados df
POOL_NAME                      USED  OBJECTS  CLONES  COPIES  RD_OPS       RD  WR_OPS       WR
.mgr                        3.1 MiB        2       0       6     116  100 KiB     141  1.9 MiB
.rgw.root                    72 KiB        6       0      18       0      0 B       0      0 B
default.rgw.buckets.data     29 MiB        4       0      12       9  9.5 MiB       0      0 B
default.rgw.buckets.index       0 B       11       0      33     115  115 KiB      17    9 KiB
default.rgw.buckets.non-ec      0 B        0       0       0       0      0 B       0      0 B
default.rgw.control             0 B        8       0      24       0      0 B       0      0 B
default.rgw.log             408 KiB      209       0     627    4350  4.2 MiB    2892      0 B
default.rgw.meta             72 KiB        7       0      21      20   16 KiB      12    6 KiB
ec-4k2m-pool                    0 B        0       0       0       0      0 B       0      0 B
```

##### 8. 파일 쪼개진것 object 확인

```
# rados --pool default.rgw.buckets.data ls
e52d5020-4ee3-448f-b912-ac6356d4d5f4.105479.1__multipart_file1.2~gPUErBGcA0XFYXVXl-dKxuziKX-wVwx.1
e52d5020-4ee3-448f-b912-ac6356d4d5f4.105479.1__shadow_file1.2~gPUErBGcA0XFYXVXl-dKxuziKX-wVwx.1_1
e52d5020-4ee3-448f-b912-ac6356d4d5f4.105479.1_file1
e52d5020-4ee3-448f-b912-ac6356d4d5f4.105479.1__multipart_file1.2~gPUErBGcA0XFYXVXl-dKxuziKX-wVwx.2
```

##### 9.  osd 10 제거하면 pg는 어떻게 할당되는가?

* 기존 사용하던 PG 24.9가  ([10,6,3], p10) ==> up ([6,9,3], p6) 로 변경 되었다.  6과3은 유지되고 10이 9로 변경됨.

```
# ceph osd rm 10
# ceph osd  map  default.rgw.buckets.data  f1
osdmap e993 pool 'default.rgw.buckets.data' (24) object 'f1' -> pg 24.67b909e9 (24.9) -> up ([6,9,3], p6) acting ([6,9,3], p6)

```

* osd 디스크를 다시 붙이기 위해서 zap 실행해서 cleanzing : lv,vg,pv 정보 제거되면서 다시 osd 디스크가 붙는다. (다행인지는 모르겠지만 osd 번호는 그대로 가져간다.  )

```
# ceph orch device zap c06 /dev/sdd --force
zap successful for /dev/sdd on c06
```

##### 10. osd disk 를 dd 해버리면

```
# ceph  osd  map  default.rgw.buckets.data  f1
osdmap e1005 pool 'default.rgw.buckets.data' (24) object 'f1' -> pg 24.67b909e9 (24.9) -> up ([10,6,3], p10) acting ([10,6,3], p10)
```

* 여기서 나온 osd 3개를 모두 dd로 밀어 버리고 hexedit를 통해서 정확히 밀어 버렸는지 확인한다.  
* osd 10, 6, 3 번 모두 밀어 버림.

```
# dd if=/dev/zero of=/dev/ceph-dd50d221-fdad-4049-a1f0-9591fc5878a1/osd-block-c441467b-5937-4b00-a44b-d2b37698af0e bs=1M
10237+0 records in
10236+0 records out
10733223936 bytes (11 GB, 10 GiB) copied, 21.7701 s, 493 MB/s
# apt install hexedit
# hexedit /dev/ceph-9d62c648-791f-4204-b117-fa074ecc77a5/osd-block-f4001470-3fc9-449d-8cbd-b42f7bdabea1
# disk 손상 
```

* 목록도 잘 보이고, 다운로드도 잘 됨 ?? 이상하네...

```
root@c01:~# aws s3  ls mybucket --endpoint-url http://192.168.105.51:80
2023-03-19 05:32:00   10000000 file1
root@c01:~# aws s3 cp s3://mybucket/file1 ~/f2 --endpoint-url http://192.168.105.51:80
download: s3://mybucket/file1 to ./f2
```

* 일단 여기 까지만 뭔가...  다른 osd에 저장되는 듯
