## Оперативни систем

Коришћен је **Lubuntu 24.04 LTS (Noble Numbat)**.
TODO: Лозинка на `starcraftrules`

## а. Покренути пројекат

`web-kangaroo` је фиктивна веб апликација написана за потребе овог задатка.
API сервер је писан у `Python 3.10` користећи `FastAPI` библиотеку. Серверу се
приступа преко `nginx`-a. Сервер комуницира са `MariaDB 11.3.2` базом података
користећи `SqlAlchemy` библиотеку.

## б. Комуникација преко мреже

Гост у VirtualBox инстанци користи NAT. Потребно је отворити порт:

![image](docs/portForwarding.png)

```
HOST                     GUEST
127.0.0.1:8300   ---->   10.0.2.15:8300
```

## ц. SSH

Треба инсталирати пакет `openssh-server`.

```bash
sudo apt-get install openssh-server
```

Следеће две команде редом започињу и заустављају SSH сервер:

```bash
sudo systemctl start sshd
sudo systemctl stop sshd
```

Уместо `start` и `stop` се може користити `enable` и `disable` како би се
покренуло чим се покрене и OS.

Да би се споља могло комуницирати са SSH сервером, треба отворити порт:

```bash
sudo ufw allow ssh
```

Пре подешавања кључа, SSH конекција се може проверити на следећи начин:

1. Пошто Linux трчи на виртуелној машини у приватној мрежи, треба отворити порт:

```
HOST                     GUEST
127.0.0.1:2222   ---->   10.0.2.15:22
```

2. Windows нативно не подржава SSH, али се може користити [Putty](https://www.putty.org/).
3. SSH Конекција се успоставља на `127.0.0.1:2222`

![image](docs/sshPutty.png)

Пошто се овде користи Putty, кључ се генерише користећи _PuttyGen_ (иначе би се користио ssh-keygen). Passphrase би требало да буде јак, овде је коришћен `bluesky123`. Приватни кључ остаје на host рачунару (тј. клијенту) а јавни кључ треба пребацити ка guest рачунару (тј. серверу). Конкретно, јавни кључ треба да се прекопира на `~/.ssh/authorized_keys`. Потенцијално се мењају и [пермисије](https://stackoverflow.com/a/49176668). Потом се на Putty алату можемо повезати ка серверу проследивши приватни кључ и passphrase.

![image](docs/sshPuttyKey.png)

## д. Secure Deployment Review

### д.1 System review

#### д.1.1 Operating System

```shell
admin@admin-virtualbox:~$ lsb_release -a
No LSB modules are available.
Distributor ID: Ubuntu
Description:    Ubuntu 22.04.4 LTS
Release:        22.04
Codename:       jammy
```

`Ubuntu 22.04.4 LTS` је подржан до јуна 2027. године, [извор](https://wiki.ubuntu.com/Releases).

```shell
admin@admin-virtualbox:~$ uname -a
Linux admin-virtualbox 6.5.0-27-generic #28~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Fri Mar 15 10:51:06 UTC 2 x86_64 x86_64 x86_64 GNU/Linux
```

Листа свих рањивости (укључујући и оних који нису нужно апликабилни) за ову верзију кернела, односно одговарајуће верзије дистрибуције, налази се [овде](https://ubuntu.com/security/cves?q=&package=&priority=critical&version=jammy&status=). Једина critical рањивост (backdoor у xz) није никада отишла у продукцију, а многе рањивости ниже опасности су митиговане.

TODO: Мало прецизније...

```shell
admin@admin-virtualbox:~$ uptime
 17:18:29 up  1:04,  6 users,  load average: 0.87, 1.06, 1.06
```

#### д.1.2 Time management

```shell
admin@admin-virtualbox:~$ cat /etc/timezone
Europe/Belgrade
```

Временска зона користи рачунање времена са померањем часовника тако да може доћи до проблема са синхронизацијом логова.

NTP није подешен:

```shell
admin@admin-virtualbox:~$ ps -edf | grep ntp
admin       1995    1847  0 18:20 ?        00:00:00 /snap/firefox/3779/usr/lib/firefox/firefox -contentproc -parentBuildID 20240206002040 -prefsLen 29341 -prefMapSize 235787 -appDir /snap/firefox/3779/usr/lib/firefox/browser {5a0d1359-463b-4c88-b3de-18776a4b1b13} 1847 true socket
admin       2029    1847  1 18:20 ?        00:01:20 /snap/firefox/3779/usr/lib/firefox/firefox -contentproc -childID 1 -isForBrowser -prefsLen 29482 -prefMapSize 235787 -jsInitLen 235124 -parentBuildID 20240206002040 -greomni /snap/firefox/3779/usr/lib/firefox/omni.ja -appomni /snap/firefox/3779/usr/lib/firefox/browser/omni.ja -appDir /snap/firefox/3779/usr/lib/firefox/browser {28f3481e-efcd-4ae1-8584-4b9009af37eb} 1847 true tab
admin       2169    1847  0 18:20 ?        00:00:26 /snap/firefox/3779/usr/lib/firefox/firefox -contentproc -childID 2 -isForBrowser -prefsLen 34983 -prefMapSize 235787 -jsInitLen 235124 -parentBuildID 20240206002040 -greomni /snap/firefox/3779/usr/lib/firefox/omni.ja -appomni /snap/firefox/3779/usr/lib/firefox/browser/omni.ja -appDir /snap/firefox/3779/usr/lib/firefox/browser {be8d8f49-3d7a-4fdd-be2f-96dc63097581} 1847 true tab
...
```

Мора се инсталирати (`sudo apt install ntp`) и онда:

```shell
admin@admin-virtualbox:~$ ps -edf | grep ntp
ntp         4410       1  0 19:47 ?        00:00:00 /usr/sbin/ntpd -p /var/run/ntpd.pid -g -u 125:133
admin       4561    4551  0 19:49 pts/1    00:00:00 grep --color=auto ntp


admin@admin-virtualbox:~$ ntpq -p -n
     remote           refid      st t when poll reach   delay   offset  jitter
==============================================================================
 0.ubuntu.pool.n .POOL.          16 p    -   64    0    0.000   +0.000   0.000
 1.ubuntu.pool.n .POOL.          16 p    -   64    0    0.000   +0.000   0.000
 2.ubuntu.pool.n .POOL.          16 p    -   64    0    0.000   +0.000   0.000
 3.ubuntu.pool.n .POOL.          16 p    -   64    0    0.000   +0.000   0.000
 ntp.ubuntu.com  .POOL.          16 p    -   64    0    0.000   +0.000   0.000
+44.190.40.123   66.220.9.122     2 u   17   64    3  170.448   +2.644  14.312
-195.178.51.145  131.188.3.221    2 u   19   64    3    5.668   -1.540  13.307
-216.177.181.129 129.7.1.66       2 u   17   64    3  155.225   +1.447  15.168
*195.178.58.245  192.168.106.2    2 u   19   64    3    3.692   -1.067  12.947
-45.79.214.107   130.207.244.240  2 u   17   64    3  130.568   +0.851  15.423
-155.248.196.28  128.138.140.44   2 u   18   64    3  190.361   -6.092  11.099
-147.91.8.1      91.187.128.199   2 u   16   64    3   49.757   -5.737  13.283
-171.66.97.126   171.64.7.105     2 u   16   64    3  169.593   -0.622  11.546
#44.31.46.123    216.218.192.202  3 u   16   64    3  168.950   +1.166   7.774
+217.24.20.5     84.16.73.33      2 u   15   64    3   10.866   -3.156   9.715
 185.125.190.56  201.68.88.106    2 u   16   64    3   44.274   -3.460  18.155
-74.208.117.38   216.239.35.4     2 u   17   64    3  152.124   -1.526  14.869
 91.189.91.157   132.163.96.1     2 u   17   64    3  102.488   -0.528  15.191
#64.142.54.12    206.55.64.77     3 u   14   64    3  185.448   +1.467  10.726
 185.125.190.57  183.160.133.132  2 u   16   64    3   34.788   -1.222  17.977
 185.125.190.58  79.243.60.50     2 u   13   64    3   44.257   -5.738  13.357

```

#### д.1.3 Packages installed

Иако `Lubuntu` има мање пакета преинсталираних од обичне `Ubuntu` инсталације, ручно скенирање `dpkg -l | less` би трајало предуго.

`sudo apt autoremove` није обрисао ништа.

За скенирање пакета за рањивости коришћен је `lynis`. Након покретања `sudo lynis audit system`, у report фајлу је пронађено следеће:

```bash
vulnerable_package[]=cpio
vulnerable_package[]=distro-info-data
vulnerable_package[]=less
vulnerable_package[]=libc-bin
vulnerable_package[]=libc6
vulnerable_package[]=libnghttp2-14
vulnerable_package[]=locales
vulnerable_package[]=python3-pil
```

Потребно је ажурирати пакете на новију верзију.

#### д.1.4 Logging

```bash
admin@admin-virtualbox:~$ ps -edf | grep syslog
message+     414       1  0 17:45 ?        00:00:01 @dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only
syslog       426       1  0 17:45 ?        00:00:00 /usr/sbin/rsyslogd -n -iNONE
admin        943     912  0 17:45 ?        00:00:00 /usr/bin/dbus-daemon --session --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only
admin      55498    4551  0 20:17 pts/1    00:00:00 grep --color=auto syslog
```

Користи се `rsyslogd`. Садржај конфигурације:

```conf
# /etc/rsyslog.conf configuration file for rsyslog
#
# For more information install rsyslog-doc and see
# /usr/share/doc/rsyslog-doc/html/configuration/index.html
#
# Default logging rules can be found in /etc/rsyslog.d/50-default.conf


#################
#### MODULES ####
#################

module(load="imuxsock") # provides support for local system logging
#module(load="immark")  # provides --MARK-- message capability

# provides UDP syslog reception
#module(load="imudp")
#input(type="imudp" port="514")

# provides TCP syslog reception
#module(load="imtcp")
#input(type="imtcp" port="514")

...
```

Видимо да је логовање искључено. Исто је јавио и `lynis`:

```bash
remote_syslog_configured=0
suggestion[]=LOGG-2154|Enable logging to an external logging host for archiving purposes and additional protection|-|-|
```

### д.2 Network Review

#### д.2.1 General Information

```bash
admin@admin-virtualbox:~$ ifconfig -a
docker0: flags=4099<UP,BROADCAST,MULTICAST>  mtu 1500
        inet 172.17.0.1  netmask 255.255.0.0  broadcast 172.17.255.255
        ether 02:42:ed:aa:0e:61  txqueuelen 0  (Ethernet)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

enp0s3: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.0.2.15  netmask 255.255.255.0  broadcast 10.0.2.255
        inet6 fe80::b45f:a1c2:7dd0:12c6  prefixlen 64  scopeid 0x20<link>
        ether 08:00:27:2e:cc:d3  txqueuelen 1000  (Ethernet)
        RX packets 360  bytes 35361 (35.3 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 364  bytes 39803 (39.8 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1000  (Local Loopback)
        RX packets 83  bytes 8380 (8.3 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 83  bytes 8380 (8.3 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

Присутни су `Docker 0`, `Ethernet Network Port 0 Seial 3` и `Loopback`

```bash
admin@admin-virtualbox:~$ route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         10.0.2.2        0.0.0.0         UG    100    0        0 enp0s3
10.0.2.0        0.0.0.0         255.255.255.0   U     100    0        0 enp0s3
172.17.0.0      0.0.0.0         255.255.0.0     U     0      0        0 docker0
```

```bash
admin@admin-virtualbox:~$ cat /etc/resolv.conf
# This is /run/systemd/resolve/stub-resolv.conf managed by man:systemd-resolved(8).
# Do not edit.
#
# This file might be symlinked as /etc/resolv.conf. If you're looking at
# /etc/resolv.conf and seeing this text, you have followed the symlink.
#
# This is a dynamic resolv.conf file for connecting local clients to the
# internal DNS stub resolver of systemd-resolved. This file lists all
# configured search domains.
#
# Run "resolvectl status" to see details about the uplink DNS servers
# currently in use.
#
# Third party programs should typically not access this file directly, but only
# through the symlink at /etc/resolv.conf. To manage man:resolv.conf(5) in a
# different way, replace this symlink by a static file or a different symlink.
#
# See man:systemd-resolved.service(8) for details about the supported modes of
# operation for /etc/resolv.conf.

nameserver 127.0.0.53
options edns0 trust-ad
search .
```

```bash
admin@admin-virtualbox:~$ cat /etc/hosts
# Standard host addresses
127.0.0.1  localhost
::1        localhost ip6-localhost ip6-loopback
ff02::1    ip6-allnodes
ff02::2    ip6-allrouters
# This host address
127.0.1.1  admin-virtualbox
```

```bash
admin@admin-virtualbox:~$ cat /etc/nsswitch.conf
# /etc/nsswitch.conf
#
# Example configuration of GNU Name Service Switch functionality.
# If you have the `glibc-doc-reference' and `info' packages installed, try:
# `info libc "Name Service Switch"' for information about this file.

passwd:         files systemd
group:          files systemd
shadow:         files
gshadow:        files

hosts:          files mdns4_minimal [NOTFOUND=return] dns
networks:       files

protocols:      db files
services:       db files
ethers:         db files
rpc:            db files

netgroup:       nis
```

#### д.2.2 Firewall Rules

Прво, firewall правила:

```bash
admin@admin-virtualbox:~$ sudo iptables -L -v
[sudo] password for admin:
Chain INPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain FORWARD (policy DROP 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DOCKER-USER  all  --  any    any     anywhere             anywhere
    0     0 DOCKER-ISOLATION-STAGE-1  all  --  any    any     anywhere             anywhere
    0     0 ACCEPT     all  --  any    docker0  anywhere             anywhere             ctstate RELATED,ESTABLISHED
    0     0 DOCKER     all  --  any    docker0  anywhere             anywhere
    0     0 ACCEPT     all  --  docker0 !docker0  anywhere             anywhere
    0     0 ACCEPT     all  --  docker0 docker0  anywhere             anywhere

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain DOCKER (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain DOCKER-ISOLATION-STAGE-1 (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DOCKER-ISOLATION-STAGE-2  all  --  docker0 !docker0  anywhere             anywhere
    0     0 RETURN     all  --  any    any     anywhere             anywhere

Chain DOCKER-ISOLATION-STAGE-2 (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DROP       all  --  any    docker0  anywhere             anywhere
    0     0 RETURN     all  --  any    any     anywhere             anywhere

Chain DOCKER-USER (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 RETURN     all  --  any    any     anywhere             anywhere
```

Пар ствари се може закључити:

- сва правила су направљена за потребе докер контејнера (који се више не користи на OS-у, те се не разматра овде али напомињемо да је лоша пракса оставити отворен саобраћај овако, поготово ако се не користи)
- Input ланац нема правила, то значи да свако може да се обрати серверу
- Forward ланац нема правила (ван докерових), што значи да се саобраћај између процеса на серверу никад не филтрира
- Output ланац нема правила, тако да сервер може свима да се обрати

За Input би требало направити правило које дозвољава HTTP саобраћај према nginx-овом порту и SSH за одређен списак IP адреса. По потреби би се додавала нова правила.

За Output би требало дозволити пролаз HTTP конекцији кроз nginx-ов порт, DNS, NTP, SSH као и конекције према серверима кроз које се ажурирају пакети (нпр. за pip и apt).

Firewall табелу можемо сачувати помоћу команде:

```bash
sudo iptables-save | sudo tee /etc/iptables.up.rules
```

(`sudo iptables-save > /etc/iptables.up.rules` није радио).

```bash
# vim etc/network/if-pre-up.d/iptables
#!/bin/bash

/sbin/iptables-restore < /etc/iptables.up.rules
```

Ажурна верзија IP табеле је сачувана на диск:

```bash
admin@admin-virtualbox:~$ cat /etc/iptables.up.rules
# Generated by iptables-save v1.8.7 on Thu May  9 16:32:28 2024
*filter
:INPUT ACCEPT [0:0]
:FORWARD DROP [0:0]
:OUTPUT ACCEPT [0:0]
:DOCKER - [0:0]
:DOCKER-ISOLATION-STAGE-1 - [0:0]
:DOCKER-ISOLATION-STAGE-2 - [0:0]
:DOCKER-USER - [0:0]
-A FORWARD -j DOCKER-USER
-A FORWARD -j DOCKER-ISOLATION-STAGE-1
-A FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A FORWARD -o docker0 -j DOCKER
-A FORWARD -i docker0 ! -o docker0 -j ACCEPT
-A FORWARD -i docker0 -o docker0 -j ACCEPT
-A DOCKER-ISOLATION-STAGE-1 -i docker0 ! -o docker0 -j DOCKER-ISOLATION-STAGE-2
-A DOCKER-ISOLATION-STAGE-1 -j RETURN
-A DOCKER-ISOLATION-STAGE-2 -o docker0 -j DROP
-A DOCKER-ISOLATION-STAGE-2 -j RETURN
-A DOCKER-USER -j RETURN
COMMIT
# Completed on Thu May  9 16:32:28 2024
# Generated by iptables-save v1.8.7 on Thu May  9 16:32:28 2024
*nat
:PREROUTING ACCEPT [0:0]
:INPUT ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
:POSTROUTING ACCEPT [0:0]
:DOCKER - [0:0]
-A PREROUTING -m addrtype --dst-type LOCAL -j DOCKER
-A OUTPUT ! -d 127.0.0.0/8 -m addrtype --dst-type LOCAL -j DOCKER
-A POSTROUTING -s 172.17.0.0/16 ! -o docker0 -j MASQUERADE
-A DOCKER -i docker0 -j RETURN
COMMIT
# Completed on Thu May  9 16:32:28 2024
```

#### д.2.3 IPv6

```bash
admin@admin-virtualbox:~$ sudo ip6tables -L -v
Chain INPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
```

Наш сервер не користи IPv6 тако да се он може искључити:

```bash
# sudo vim /etc/sysctl.conf

# IPV6

net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
net.ipv6.conf.tun0.disable_ipv6 = 1
```

### д.3 Filesystem review

#### д.3.1 Mounted partitions

```bash
admin@admin-virtualbox:~$ cat /etc/fstab
# /etc/fstab: static file system information.
#
# Use 'blkid' to print the universally unique identifier for a device; this may
# be used with UUID= as a more robust way to name devices that works even if
# disks are added and removed. See fstab(5).
#
# <file system>             <mount point>  <type>  <options>  <dump>  <pass>
UUID=3c4d532c-4a9a-4d49-9a66-2c280cebda89 /              ext4    defaults   0 1
/swapfile                                 swap           swap    defaults   0 0
```

Имамо само главни фајлсистем и swapfile. `defaults` опција подразумева `rw,
suid, dev, exec, auto, nouser, async`. То значи да `noatime` није коришћен,
али `exec` и `suid` јесу.

#### д.3.2 Sensitive Files

Видимо да је `shadow` доступан корисницима са root привилегијом:

```bash
sudo cat /etc/shadow
root:$y$j9T$LzAM7Bku1KPtjXT45Vgx70$ww3.7lmFS8pfRZUOW9NvJNnK.k1X35j96w7Cr6oD3OA:19850:0:99999:7:::
daemon:*:19769:0:99999:7:::
bin:*:19769:0:99999:7:::
sys:*:19769:0:99999:7:::
sync:*:19769:0:99999:7:::
games:*:19769:0:99999:7:::
man:*:19769:0:99999:7:::
lp:*:19769:0:99999:7:::
mail:*:19769:0:99999:7:::
news:*:19769:0:99999:7:::
uucp:*:19769:0:99999:7:::
proxy:*:19769:0:99999:7:::
www-data:*:19769:0:99999:7:::
backup:*:19769:0:99999:7:::
list:*:19769:0:99999:7:::
irc:*:19769:0:99999:7:::
gnats:*:19769:0:99999:7:::
nobody:*:19769:0:99999:7:::
systemd-network:*:19769:0:99999:7:::
systemd-resolve:*:19769:0:99999:7:::
messagebus:*:19769:0:99999:7:::
systemd-timesync:*:19769:0:99999:7:::
syslog:*:19769:0:99999:7:::
_apt:*:19769:0:99999:7:::
tss:*:19769:0:99999:7:::
uuidd:*:19769:0:99999:7:::
tcpdump:*:19769:0:99999:7:::
usbmux:*:19769:0:99999:7:::
dnsmasq:*:19769:0:99999:7:::
kernoops:*:19769:0:99999:7:::
avahi:*:19769:0:99999:7:::
cups-pk-helper:*:19769:0:99999:7:::
rtkit:*:19769:0:99999:7:::
whoopsie:*:19769:0:99999:7:::
fwupd-refresh:*:19769:0:99999:7:::
saned:*:19769:0:99999:7:::
colord:*:19769:0:99999:7:::
sddm:*:19769:0:99999:7:::
geoclue:*:19769:0:99999:7:::
pulse:*:19769:0:99999:7:::
hplip:*:19769:0:99999:7:::
admin:$6$ReFz5oqz6qJw13ly$HcS9/qOfBt88ZANfSd0vw6pG8jYMSGATSknIn8X4X4X/IpqQ0K.cWW3kJtbxSg5g3wtTCRpps4EWoDhbybpCG.:19830:0:99999:7:::
mysql:!:19850:0:99999:7:::
sshd:*:19850:0:99999:7:::
ntp:*:19851:0:99999:7:::
```

Док је MariaDB фајл доступан свима.

```bash
admin@admin-virtualbox:~$ cat /etc/mysql/my.cnf
# The MariaDB configuration file
#
# The MariaDB/MySQL tools read configuration files in the following order:
# 0. "/etc/mysql/my.cnf" symlinks to this file, reason why all the rest is read.
# 1. "/etc/mysql/mariadb.cnf" (this file) to set global defaults,
# 2. "/etc/mysql/conf.d/*.cnf" to set global options.
# 3. "/etc/mysql/mariadb.conf.d/*.cnf" to set MariaDB-only options.
# 4. "~/.my.cnf" to set user-specific options.
#
# If the same option is defined multiple times, the last one will apply.
#
# One can use all long options that the program supports.
# Run program with --help to get a list of available options and with
# --print-defaults to see which it would actually understand and use.
#
# If you are new to MariaDB, check out https://mariadb.com/kb/en/basic-mariadb-articles/

#
# This group is read both by the client and the server
# use it for options that affect everything
#
[client-server]
# Port or socket location where to connect
# port = 3306
socket = /run/mysqld/mysqld.sock

# Import all .cnf files from configuration directory
!includedir /etc/mysql/conf.d/
!includedir /etc/mysql/mariadb.conf.d/
```

Фајл `shadow.backup` тренутно не постоји:

```bash
admin@admin-virtualbox:~$ cat /etc/shadow.backup
cat: /etc/shadow.backup: No such file or directory
```

#### д.3.3 Setuid

```bash
admin@admin-virtualbox:~$ sudo find / -perm -4000 -ls
[sudo] password for admin:
      297    129 -rwsr-xr-x   1 root     root       131832 Nov 29 15:54 /snap/snapd/20671/usr/lib/snapd/snap-confine
      882     72 -rwsr-xr-x   1 root     root        72712 Feb  6 13:54 /snap/core22/1380/usr/bin/chfn
      888     44 -rwsr-xr-x   1 root     root        44808 Feb  6 13:54 /snap/core22/1380/usr/bin/chsh
      954     71 -rwsr-xr-x   1 root     root        72072 Feb  6 13:54 /snap/core22/1380/usr/bin/gpasswd
     1038     47 -rwsr-xr-x   1 root     root        47488 Mar 22 13:25 /snap/core22/1380/usr/bin/mount
     1047     40 -rwsr-xr-x   1 root     root        40496 Feb  6 13:54 /snap/core22/1380/usr/bin/newgrp
     1062     59 -rwsr-xr-x   1 root     root        59976 Feb  6 13:54 /snap/core22/1380/usr/bin/passwd
     1180     55 -rwsr-xr-x   1 root     root        55680 Mar 22 13:25 /snap/core22/1380/usr/bin/su
     1181    227 -rwsr-xr-x   1 root     root       232416 Apr  3  2023 /snap/core22/1380/usr/bin/sudo
     1241     35 -rwsr-xr-x   1 root     root        35200 Mar 22 13:25 /snap/core22/1380/usr/bin/umount
     1333     35 -rwsr-xr--   1 root     systemd-resolve    35112 Oct 25  2022 /snap/core22/1380/usr/lib/dbus-1.0/dbus-daemon-launch-helper
     2602    331 -rwsr-xr-x   1 root     root              338536 Jan  2 17:54 /snap/core22/1380/usr/lib/openssh/ssh-keysign
     8632     19 -rwsr-xr-x   1 root     root               18736 Feb 26  2022 /snap/core22/1380/usr/libexec/polkit-agent-helper-1
      879     72 -rwsr-xr-x   1 root     root               72712 Nov 24  2022 /snap/core22/1122/usr/bin/chfn
      885     44 -rwsr-xr-x   1 root     root               44808 Nov 24  2022 /snap/core22/1122/usr/bin/chsh
      951     71 -rwsr-xr-x   1 root     root               72072 Nov 24  2022 /snap/core22/1122/usr/bin/gpasswd
     1035     47 -rwsr-xr-x   1 root     root               47480 Feb 21  2022 /snap/core22/1122/usr/bin/mount
     1044     40 -rwsr-xr-x   1 root     root               40496 Nov 24  2022 /snap/core22/1122/usr/bin/newgrp
     1059     59 -rwsr-xr-x   1 root     root               59976 Nov 24  2022 /snap/core22/1122/usr/bin/passwd
     1177     55 -rwsr-xr-x   1 root     root               55672 Feb 21  2022 /snap/core22/1122/usr/bin/su
     1178    227 -rwsr-xr-x   1 root     root              232416 Apr  3  2023 /snap/core22/1122/usr/bin/sudo
     1238     35 -rwsr-xr-x   1 root     root               35192 Feb 21  2022 /snap/core22/1122/usr/bin/umount
     1330     35 -rwsr-xr--   1 root     systemd-resolve    35112 Oct 25  2022 /snap/core22/1122/usr/lib/dbus-1.0/dbus-daemon-launch-helper
     2599    331 -rwsr-xr-x   1 root     root              338536 Jan  2 17:54 /snap/core22/1122/usr/lib/openssh/ssh-keysign
     8618     19 -rwsr-xr-x   1 root     root               18736 Feb 26  2022 /snap/core22/1122/usr/libexec/polkit-agent-helper-1
  1057822     40 -rwsr-xr-x   1 root     root               37616 Sep 12  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/newgrp
  1057819     64 -rwsr-xr-x   1 root     root               64152 Sep 12  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/chage
  1058580     36 -rwsr-xr-x   1 root     root               35952 Jan 23  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/mount
  1057820     80 -rwsr-xr-x   1 root     root               78120 Sep 12  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/gpasswd
  1058595     32 -rwsr-xr-x   1 root     root               32032 Jan 23  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/su
  1058598     28 -rwsr-xr-x   1 root     root               27776 Jan 23  2023 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/bin/umount
  1059815     60 -rwsr-x---   1 root     81                 57856 Jan 10 22:38 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/libexec/dbus-1/dbus-daemon-launch-helper
  1057434     12 -rwsr-xr-x   1 root     root               11152 Jan 29 04:54 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/sbin/pam_timestamp_check
  1057436     36 -rwsr-xr-x   1 root     root               36176 Jan 29 04:54 /var/snap/docker/common/var-lib-docker/overlay2/4rwb62v8446vdom8l7ppr4oa2/diff/usr/sbin/unix_chkpwd
   263348     40 -rwsr-xr-x   1 root     root               40496 Feb  6 13:54 /usr/bin/newgrp
   263420     60 -rwsr-xr-x   1 root     root               59976 Feb  6 13:54 /usr/bin/passwd
   263864    228 -rwsr-xr-x   1 root     root              232416 Apr  3  2023 /usr/bin/sudo
   262818     44 -rwsr-xr-x   1 root     root               44808 Feb  6 13:54 /usr/bin/chsh
   262549     48 -rwsr-xr-x   1 root     root               47488 Apr  9 17:32 /usr/bin/mount
   262812     72 -rwsr-xr-x   1 root     root               72712 Feb  6 13:54 /usr/bin/chfn
   263022     72 -rwsr-xr-x   1 root     root               72072 Feb  6 13:54 /usr/bin/gpasswd
   262857     56 -rwsr-xr-x   1 root     root               55680 Apr  9 17:32 /usr/bin/su
   263523     32 -rwsr-xr-x   1 root     root               30872 Feb 26  2022 /usr/bin/pkexec
   262988     36 -rwsr-xr-x   1 root     root               35200 Mar 23  2022 /usr/bin/fusermount3
   262642     36 -rwsr-xr-x   1 root     root               35200 Apr  9 17:32 /usr/bin/umount
   266228     16 -rwsr-xr-x   1 root     root               14488 Dec  2 04:44 /usr/lib/mysql/plugin/auth_pam_tool_dir/auth_pam_tool
   262367     16 -rwsr-sr-x   1 root     root               14488 Apr  9 05:18 /usr/lib/xorg/Xorg.wrap
   270442    332 -rwsr-xr-x   1 root     root              338536 Mar 15 21:28 /usr/lib/openssh/ssh-keysign
   262629    140 -rwsr-xr-x   1 root     root              142536 Mar  6 22:18 /usr/lib/snapd/snap-confine
   264691     36 -rwsr-xr--   1 root     messagebus         35112 Oct 25  2022 /usr/lib/dbus-1.0/dbus-daemon-launch-helper
   271690     20 -rwsr-xr-x   1 root     root               18736 Feb 26  2022 /usr/libexec/polkit-agent-helper-1
   264595     56 -rwsr-xr-x   1 root     root               54184 Mar 22 15:00 /usr/share/code/chrome-sandbox
   271912    416 -rwsr-xr--   1 root     dip               424512 Feb 24  2022 /usr/sbin/pppd
```

Интегритет пакета се може добити на `/var/lib/dpkg/info/[NAME].md5sums`.
У том фајлу се налазе чексуме за све фајлове одговарајућег пакета. Онда бисмо покренули `md5sum ...` и поредили ручно. Други начин је користити пакет `debsums`:

```bash
admin@admin-virtualbox:~$ sudo debsums -c -s
```

Ништа није исписано, те је интегритет свих фајлова очуван.

Што се тиче права приступа, треба минимизовати програме који имају `setuid`.
Програми који треба да имају `setuid` су `su`, `sudo`, `passwd`, сервери за X.
`mount` и `umount` исто захтевају setuid, осим ако се не користи нека
алтернатива попут `udisk`. Сви остали програми не би требало да имају `setuid`
подешен, већ се покрећу са `sudo`, тако да пермисије нису адекватно подешене.

#### д.3.4 Normal Files

```bash
admin@admin-virtualbox:~$ find / -type f -perm -006 2>/dev/null | grep -v /proc
/sys/kernel/security/apparmor/.remove
/sys/kernel/security/apparmor/.replace
/sys/kernel/security/apparmor/.load
/sys/kernel/security/apparmor/.notify
/sys/kernel/security/apparmor/.access
```

```bash
admin@admin-virtualbox:~$ find / -type f -perm -002 2>/dev/null | grep -v /proc
/sys/kernel/security/apparmor/.remove
/sys/kernel/security/apparmor/.replace
/sys/kernel/security/apparmor/.load
/sys/kernel/security/apparmor/.notify
/sys/kernel/security/apparmor/.access
```

Једине датотеке које могу read-write сви корисници су фајлови везани за
AppArmor.

#### д.3.5 Backup

Lubuntu подразумевано нема `backup` алат нити `/backup` фолдер. У случају да је
прављење backup-а потребно (јесте), онда се треба инсталирати `backup` и
подесити пермисије над целим фолдером.

### д.4 User review

#### д.4.1 Reviewing the passwd file

```bash
admin@admin-virtualbox:/bin$ cat /etc/passwd | grep admin
gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
admin:x:1000:1000:admin:/home/admin:/bin/bash
```

`admin` је назив корисничког налога на овом OS-u. Видимо да је ID 1000 тј. није
0 (root). Подразумевани shell за овог корисника је `/bin/bash`.

Следи преглед целог `/etc/passwrd` фајла.

```bash
admin@admin-virtualbox:/bin$ cat /etc/passwd
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
sys:x:3:3:sys:/dev:/usr/sbin/nologin
sync:x:4:65534:sync:/bin:/bin/sync
games:x:5:60:games:/usr/games:/usr/sbin/nologin
man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
news:x:9:9:news:/var/spool/news:/usr/sbin/nologin
uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin
proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
backup:x:34:34:backup:/var/backups:/usr/sbin/nologin
list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
irc:x:39:39:ircd:/run/ircd:/usr/sbin/nologin
gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
systemd-network:x:100:102:systemd Network Management,,,:/run/systemd:/usr/sbin/nologin
systemd-resolve:x:101:103:systemd Resolver,,,:/run/systemd:/usr/sbin/nologin
messagebus:x:102:105::/nonexistent:/usr/sbin/nologin
systemd-timesync:x:103:106:systemd Time Synchronization,,,:/run/systemd:/usr/sbin/nologin
syslog:x:104:111::/home/syslog:/usr/sbin/nologin
_apt:x:105:65534::/nonexistent:/usr/sbin/nologin
tss:x:106:112:TPM software stack,,,:/var/lib/tpm:/bin/false
uuidd:x:107:115::/run/uuidd:/usr/sbin/nologin
tcpdump:x:108:116::/nonexistent:/usr/sbin/nologin
usbmux:x:109:46:usbmux daemon,,,:/var/lib/usbmux:/usr/sbin/nologin
dnsmasq:x:110:65534:dnsmasq,,,:/var/lib/misc:/usr/sbin/nologin
kernoops:x:111:65534:Kernel Oops Tracking Daemon,,,:/:/usr/sbin/nologin
avahi:x:112:119:Avahi mDNS daemon,,,:/run/avahi-daemon:/usr/sbin/nologin
cups-pk-helper:x:113:120:user for cups-pk-helper service,,,:/home/cups-pk-helper:/usr/sbin/nologin
rtkit:x:114:121:RealtimeKit,,,:/proc:/usr/sbin/nologin
whoopsie:x:115:122::/nonexistent:/bin/false
fwupd-refresh:x:116:123:fwupd-refresh user,,,:/run/systemd:/usr/sbin/nologin
saned:x:117:125::/var/lib/saned:/usr/sbin/nologin
colord:x:118:126:colord colour management daemon,,,:/var/lib/colord:/usr/sbin/nologin
sddm:x:119:127:Simple Desktop Display Manager:/var/lib/sddm:/bin/false
geoclue:x:120:128::/var/lib/geoclue:/usr/sbin/nologin
pulse:x:121:129:PulseAudio daemon,,,:/run/pulse:/usr/sbin/nologin
hplip:x:122:7:HPLIP system user,,,:/run/hplip:/bin/false
admin:x:1000:1000:admin:/home/admin:/bin/bash
mysql:x:123:132:MySQL Server,,,:/nonexistent:/bin/false
sshd:x:124:65534::/run/sshd:/usr/sbin/nologin
ntp:x:125:133::/nonexistent:/usr/sbin/nologin
```

Једини корисници који имају shell су:

- `root`, и то `/bin/bash`
- `sync`, и то `/bin/sync`
- `admin`, и то `/bin/bash`

#### д.4.2 Reviewing the shadow file

Све лозинке у `/etc/passwd` су хеширане.

Садржај `/etc/shadow` изгледа овако:

```bash
admin@admin-virtualbox:/bin$ sudo cat /etc/shadow
[sudo] password for admin:
root:$y$j9T$LzAM7Bku1KPtjXT45Vgx70$ww3.7lmFS8pfRZUOW9NvJNnK.k1X35j96w7Cr6oD3OA:19850:0:99999:7:::
daemon:*:19769:0:99999:7:::
bin:*:19769:0:99999:7:::
sys:*:19769:0:99999:7:::
sync:*:19769:0:99999:7:::
games:*:19769:0:99999:7:::
man:*:19769:0:99999:7:::
lp:*:19769:0:99999:7:::
mail:*:19769:0:99999:7:::
news:*:19769:0:99999:7:::
uucp:*:19769:0:99999:7:::
proxy:*:19769:0:99999:7:::
www-data:*:19769:0:99999:7:::
backup:*:19769:0:99999:7:::
list:*:19769:0:99999:7:::
irc:*:19769:0:99999:7:::
gnats:*:19769:0:99999:7:::
nobody:*:19769:0:99999:7:::
systemd-network:*:19769:0:99999:7:::
systemd-resolve:*:19769:0:99999:7:::
messagebus:*:19769:0:99999:7:::
systemd-timesync:*:19769:0:99999:7:::
syslog:*:19769:0:99999:7:::
_apt:*:19769:0:99999:7:::
tss:*:19769:0:99999:7:::
uuidd:*:19769:0:99999:7:::
tcpdump:*:19769:0:99999:7:::
usbmux:*:19769:0:99999:7:::
dnsmasq:*:19769:0:99999:7:::
kernoops:*:19769:0:99999:7:::
avahi:*:19769:0:99999:7:::
cups-pk-helper:*:19769:0:99999:7:::
rtkit:*:19769:0:99999:7:::
whoopsie:*:19769:0:99999:7:::
fwupd-refresh:*:19769:0:99999:7:::
saned:*:19769:0:99999:7:::
colord:*:19769:0:99999:7:::
sddm:*:19769:0:99999:7:::
geoclue:*:19769:0:99999:7:::
pulse:*:19769:0:99999:7:::
hplip:*:19769:0:99999:7:::
admin:$6$ReFz5oqz6qJw13ly$HcS9/qOfBt88ZANfSd0vw6pG8jYMSGATSknIn8X4X4X/IpqQ0K.cWW3kJtbxSg5g3wtTCRpps4EWoDhbybpCG.:19830:0:99999:7:::
mysql:!:19850:0:99999:7:::
sshd:*:19850:0:99999:7:::
ntp:*:19851:0:99999:7:::
```

`root` лозинка почиње са `$y$` што значи да користи `yescrypt` алгоритам који
је много бољи у спречавању offline напада.

`admin` лозинка почиње са `$6$` тј. SHA-512 што је јак алгоритам.

Остали налози су `*` тј. login није могућ на овим налозима.

```bash
admin@admin-virtualbox:/bin$ cat /etc/pam.d/common-password
#
# /etc/pam.d/common-password - password-related modules common to all services
#
# This file is included from other service-specific PAM config files,
# and should contain a list of modules that define the services to be
# used to change user passwords.  The default is pam_unix.

# Explanation of pam_unix options:
# The "yescrypt" option enables
#hashed passwords using the yescrypt algorithm, introduced in Debian
#11.  Without this option, the default is Unix crypt.  Prior releases
#used the option "sha512"; if a shadow password hash will be shared
#between Debian 11 and older releases replace "yescrypt" with "sha512"
#for compatibility .  The "obscure" option replaces the old
#`OBSCURE_CHECKS_ENAB' option in login.defs.  See the pam_unix manpage
#for other options.

# As of pam 1.0.1-6, this file is managed by pam-auth-update by default.
# To take advantage of this, it is recommended that you configure any
# local modules either before or after the default block, and use
# pam-auth-update to manage selection of other modules.  See
# pam-auth-update(8) for details.

# here are the per-package modules (the "Primary" block)
password        [success=1 default=ignore]      pam_unix.so obscure yescrypt
# here's the fallback if no module succeeds
password        requisite                       pam_deny.so
# prime the stack with a positive return value if there isn't one already;
# this avoids us returning an error just because nothing sets a success code
# since the modules above will each just jump around
password        required                        pam_permit.so
# and here are more per-package modules (the "Additional" block)
password        optional        pam_gnome_keyring.so
# end of pam-auth-update config
```

Користи се `yescrypt`.

Добављање лозинке помоћу John-а:

```bash
admin@admin-virtualbox:/bin$ john
    [...]
--format=NAME              force hash type NAME: descrypt/bsdicrypt/md5crypt/
                           bcrypt/LM/AFS/tripcode/dummy/crypt
    [...]
admin@admin-virtualbox:/bin$ sudo john /etc/shadow
Created directory: /root/.john
Loaded 1 password hash (crypt, generic crypt(3) [?/64])
Will run 2 OpenMP threads
Press 'q' or Ctrl-C to abort, almost any other key for status
1234             (admin)
1g 0:00:00:05 100% 2/3 0.1945g/s 588.5p/s 588.5c/s 588.5C/s 123456..pepper
Use the "--show" option to display all of the cracked passwords reliably
Session completed
admin@admin-virtualbox:/bin$ sudo john /etc/shadow --show
admin:1234:19830:0:99999:7:::
```

`1234` је лозинка admin налога.

#### д.4.3 Reviewing the sudo configuration

```bash
admin@admin-virtualbox:/bin$ sudo egrep -v '^#|^$' /etc/sudoers
Defaults        env_reset
Defaults        mail_badpass
Defaults        secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin"
Defaults        use_pty
root    ALL=(ALL:ALL) ALL
%admin ALL=(ALL) ALL
%sudo   ALL=(ALL:ALL) ALL
@includedir /etc/sudoers.d
```

Видимо да `admin` група може приступити свим налозима (укључујући и root)
налогу, али мора имати лозинку.

### д.5 Services review

#### д.5.1 Identifying running services

```bash
admin@admin-virtualbox:~$ lsof -i UDP -n -P
admin@admin-virtualbox:~$ lsof -i TCP -n -P
COMMAND  PID  USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
uvicorn 1681 admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3 1683 admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3 1683 admin   18u  IPv4  27593      0t0  TCP 127.0.0.1:36812->127.0.0.1:3306 (ESTABLISHED)
```

Немамо UDP сервисе, а за TCP су сервиси API сервера. Nginx и MariaDB нису покренути као сервиси те се не налазе овде. Чак ни sshd није покренут као сервис.
Да су заиста покренути њихови процеси, види се овде:

```bash
admin@admin-virtualbox:~$ ps -edf | grep ngin
root        1690       1  0 14:06 ?        00:00:00 nginx: master process nginx
nobody      1691    1690  0 14:06 ?        00:00:00 nginx: worker process
admin       1700    1458  0 14:08 pts/0    00:00:00 grep --color=auto ngin
admin@admin-virtualbox:~$ ps -edf | grep mari
mysql        651       1  0 14:00 ?        00:00:02 /usr/sbin/mariadbd
admin       1704    1458  0 14:08 pts/0    00:00:00 grep --color=auto mari
admin@admin-virtualbox:~$ ps -edf | grep sshd
root         592       1  0 13:59 ?        00:00:00 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups
root        1379     592  0 14:00 ?        00:00:00 sshd: admin [priv]
admin       1456    1379  0 14:00 ?        00:00:01 sshd: admin@pts/0
admin       1706    1458  0 14:08 pts/0    00:00:00 grep --color=auto sshd
```

#### д.5.2 OpenSSH

```bash
admin@admin-virtualbox:~$ cat /etc/ssh/sshd_config | grep PermitRootLogin
#PermitRootLogin prohibit-password
# the setting of "PermitRootLogin without-password".
```

Закоментарисано се. Потребно је ставити `no` како `root` не би могао приступити SSH јер је то опасно.

Помоћу `protocol 2` се ограничава употреба на верзију 2 од SSH.

Порт је 22, изменом порта на неку другу вредност отежава нападачима да лоцирају SSH сокет. Потребно је променити и порт у VirtualBox-овом port forwarding опцијама, као и конфигурацију у PuTTY-ју.

Коначно, `AllowTcpForwarding` је постављен на `no`.

Пре рестартовања ssh, добро је проверити да нема грешака у конфигурацији:

```
sshd -t
```

Рестартовање на системима који користе systemd:

```
sudo systemctl restart sshd
```

#### д.5.3 MySQL

У конфигурацији `/etc/mysql/my.cnf` додато је:

```conf
[mysqld]
bind-address = 127.0.0.1
```

Конекција (лозинка је `root`):

```bash
sudo mariadb -u root -proot
```

```bash
MariaDB [(none)]> select @@version;
+----------------------------------+
| @@version                        |
+----------------------------------+
| 10.6.16-MariaDB-0ubuntu0.22.04.1 |
+----------------------------------+
1 row in set (0.000 sec)
```

За списак рањивости на овој верзији, погледати рањивости које су
печоване у новијим верзијама, на овом сајту: https://mariadb.com/kb/en/security/.

```bash
MariaDB [(none)]> select host, user, password from mysql.user;
+-----------+-------------+-------------------------------------------+
| Host      | User        | Password                                  |
+-----------+-------------+-------------------------------------------+
| localhost | mariadb.sys |                                           |
| localhost | root        | *81F5E21E35407D884A6CD4A731AEBFB6AF209E1B |
| localhost | mysql       | invalid                                   |
+-----------+-------------+-------------------------------------------+
3 rows in set (0.045 sec)
```

Видимо да `root` има лозинку која користи нови алгоритам за енкрипцију.

Крековање лозинке помоћу John-а је било једноставно с обзиром да је лозинка `root` и коришћен је `mysql-sha1`.

```
sudo unshadow ./mysql-password /etc/shadow > mypasswd
john mypasswd --format=crypt

Loaded 1 password hash (crypt, generic crypt(3) [?/64])
Will run 2 OpenMP threads
Press 'q' or Ctrl-C to abort, almost any other key for status
0g 0:00:00:08 57% 1/3 0g/s 33.21p/s 33.21c/s 33.21C/s Drroot..Root02
```

#### д.5.4 Nginx

Прво да видимо под којим корисником је покренут nginx:

```bash
admin@admin-virtualbox:~$ sudo lsof -nP -i
[sudo] password for admin:
COMMAND    PID            USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
systemd-r  335 systemd-resolve   13u  IPv4  21036      0t0  UDP 127.0.0.53:53
systemd-r  335 systemd-resolve   14u  IPv4  21037      0t0  TCP 127.0.0.53:53 (LISTEN)
avahi-dae  436           avahi   12u  IPv4  22158      0t0  UDP *:5353
avahi-dae  436           avahi   13u  IPv6  22159      0t0  UDP *:5353
avahi-dae  436           avahi   14u  IPv4  22160      0t0  UDP *:51154
avahi-dae  436           avahi   15u  IPv6  22161      0t0  UDP *:58342
NetworkMa  441            root   25u  IPv4  22484      0t0  UDP 10.0.2.15:68->10.0.2.2:67
cupsd      537            root    7u  IPv4  22363      0t0  TCP 127.0.0.1:631 (LISTEN)
cups-brow  667            root    7u  IPv4  23311      0t0  UDP *:631
uvicorn   1681           admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3   1683           admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3   1683           admin   18u  IPv4  27593      0t0  TCP 127.0.0.1:36812->127.0.0.1:3306 (CLOSE_WAIT)
nginx     1690            root    6u  IPv4  32986      0t0  TCP *:8300 (LISTEN)
nginx     1691          nobody    6u  IPv4  32986      0t0  TCP *:8300 (LISTEN)
sshd      1803            root    3u  IPv4  33416      0t0  TCP *:27 (LISTEN)
sshd      1803            root    4u  IPv6  33418      0t0  TCP *:27 (LISTEN)
sshd      1810            root    4u  IPv4  33430      0t0  TCP 10.0.2.15:27->10.0.2.2:51980 (ESTABLISHED)
sshd      1845           admin    4u  IPv4  33430      0t0  TCP 10.0.2.15:27->10.0.2.2:51980 (ESTABLISHED)
mariadbd  1959           mysql   20u  IPv4  34183      0t0  TCP 127.0.0.1:3306 (LISTEN)
```

Видимо `nginx` је под `root`. Ако погледамо конфигурацију:

```nginx
events {

}

http {
        server {
                listen 8300;

                location /api/ {
                        proxy_pass http://localhost:8000/;
                }
        }
}
```

ту видимо да нема ништа за подешавање корисника. Додавањем

```
user www-data;
```

на почетку поправља ствари. Међутим, постоје две инстанце nginx-a. Прва је
главни процес који не ради ништа већ се на основу њега fork-ују робови.

Рестартовањем nginx сервера видимо:

```bash
admin@admin-virtualbox:~$ sudo lsof -nP -i
COMMAND    PID            USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
systemd-r  335 systemd-resolve   13u  IPv4  21036      0t0  UDP 127.0.0.53:53
systemd-r  335 systemd-resolve   14u  IPv4  21037      0t0  TCP 127.0.0.53:53 (LISTEN)
avahi-dae  436           avahi   12u  IPv4  22158      0t0  UDP *:5353
avahi-dae  436           avahi   13u  IPv6  22159      0t0  UDP *:5353
avahi-dae  436           avahi   14u  IPv4  22160      0t0  UDP *:51154
avahi-dae  436           avahi   15u  IPv6  22161      0t0  UDP *:58342
NetworkMa  441            root   25u  IPv4  22484      0t0  UDP 10.0.2.15:68->10.0.2.2:67
cupsd      537            root    7u  IPv4  22363      0t0  TCP 127.0.0.1:631 (LISTEN)
cups-brow  667            root    7u  IPv4  23311      0t0  UDP *:631
uvicorn   1681           admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3   1683           admin    3u  IPv4  27557      0t0  TCP 127.0.0.1:8000 (LISTEN)
python3   1683           admin   18u  IPv4  27593      0t0  TCP 127.0.0.1:36812->127.0.0.1:3306 (CLOSE_WAIT)
sshd      1803            root    3u  IPv4  33416      0t0  TCP *:27 (LISTEN)
sshd      1803            root    4u  IPv6  33418      0t0  TCP *:27 (LISTEN)
sshd      1810            root    4u  IPv4  33430      0t0  TCP 10.0.2.15:27->10.0.2.2:51980 (ESTABLISHED)
sshd      1845           admin    4u  IPv4  33430      0t0  TCP 10.0.2.15:27->10.0.2.2:51980 (ESTABLISHED)
mariadbd  1959           mysql   20u  IPv4  34183      0t0  TCP 127.0.0.1:3306 (LISTEN)
nginx     2443            root    6u  IPv4  35676      0t0  TCP *:8300 (LISTEN)
nginx     2444        www-data    6u  IPv4  35676      0t0  TCP *:8300 (LISTEN)
```

Један процес је под www-data, али други је и даље root. Он мора бити root и
он не ради ништа сем што се fork-ује.

Токени се искључују тако што се се постави `server_tokens off` у конфигурацији.

Не постоји никаква контрола и ограничавање ресурса што олакшава DoS нападе.
Пожељно је додати директиве `client_body_buffer_size 1k`,
`client_header_buffer_size 1k`, `client_max_body_size 1k`,
`large_client_header_bufffers 2 1k`.

Треба ограничити које HTTP методе веб сервер прихвата:

```nginx
location /api/ {
    limit_except GET HEAD POST {
        deny all;
    }
    ...
}
```

#### д.5.5 Python Configuration

Сам апи сервер нема конфигурације, а садржај у nginx-у је покривен претходном
секцијом. TODO: Шта овде писати...?

#### д.5.6 Crontab

```bash
admin@admin-virtualbox:~/Desktop/web-kangaroo/api$ sudo crontab -u root -l
[sudo] password for admin:
no crontab for root
admin@admin-virtualbox:~/Desktop/web-kangaroo/api$ sudo crontab -u admin -l
no crontab for admin
```

```
admin@admin-virtualbox:/var/spool/cron$ sudo ls /var/spool/cron/crontabs/
admin@admin-virtualbox:/var/spool/cron$
```

Видимо да немамо ниједан крон џоб.

## е. Извлачење хеш лозинке

Описано је изнад.
