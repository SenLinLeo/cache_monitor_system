# cache_monitor_system
基于Grafana+Prometheus的缓存（Redis、Codis）监控系统

概述
Prometheus是一个开源的服务监控系统，它通过HTTP协议从远程的机器收集数据并存储在本地的时序数据库上。
    多维数据模型（时序列数据由metric名和一组key/value组成）
    在多维度上灵活的查询语言(PromQl)
    不依赖分布式存储，单主节点工作.
    通过基于HTTP的pull方式采集时序数据
    可以通过push gateway进行时序列数据推送(pushing)
    可以通过服务发现或者静态配置去获取要采集的目标服务器
    多种可视化图表及仪表盘支持
    Prometheus通过安装在远程机器上的cache_monitor来收集监控数据，后面我们将使用到cache_monitor收集系统数据。

Grafana 是一个开箱即用的可视化工具，具有功能齐全的度量仪表盘和图形编辑器，有灵活丰富的图形化选项，可以混合多种风格，支持多个数据源特点。

安装
  下载：
     git clone https://github.com/SenLinLeo/cache_monitor_system.git
  运行
    go run main.goPrometheus
下载地址：https://prometheus.io/download

执行以下命令：

## 下载
wget https://github.com/prometheus/prometheus/releases/download/v2.0.0-rc.3/prometheus-2.0.0-rc.3.linux-amd64.tar.gz
## 可自定义解压目录
tar -xvf prometheus-2.0.0-rc.3.linux-amd64.tar.gz
配置prometheus，vi prometheus.yml

global:
  scrape_interval:     15s
  evaluation_interval: 15s
  
  - job_name: prometheus
    static_configs:
      - targets: ['localhost:9090']
        labels:
          instance: prometheus
          
  - job_name: linux1
    static_configs:
      - targets: ['192.168.1.120:9100']
        labels:
          instance: sys1
          
  - job_name: linux2
    static_configs:
      - targets: ['192.168.1.130:9100']
        labels:
          instance: sys2
IP对应的是我们内网的服务器，端口则是对应的exporter的监听端口。

运行Prometheus
./prometheus 
level=info ts=2017-11-07T02:39:50.220187934Z caller=main.go:215 msg="Starting Prometheus" version="(version=2.0.0-rc.2, branch=HEAD, revision=ce63a5a8557bb33e2030a7756c58fd773736b592)"
level=info ts=2017-11-07T02:39:50.22025258Z caller=main.go:216 build_context="(go=go1.9.1, user=root@a6d2e4a7b8da, date=20171025-18:42:54)"
level=info ts=2017-11-07T02:39:50.220270139Z caller=main.go:217 host_details="(Linux 3.10.0-514.16.1.el7.x86_64 #1 SMP Wed Apr 12 15:04:24 UTC 2017 x86_64 iZ2ze74fkxrls31tr2ia2fZ (none))"
level=info ts=2017-11-07T02:39:50.223171565Z caller=web.go:380 component=web msg="Start listening for connections" address=0.0.0.0:9090
......
启动成功以后我们可以通过Prometheus内置了web界面访问，http://ip:9090 ，如果出现以下界面，说明配置成功
