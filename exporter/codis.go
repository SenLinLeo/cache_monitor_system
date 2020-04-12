package exporter

import (
	"codis-monitor/common"
	"codis-monitor/conf"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (e *Exporter) extractCodisMetrics(addr string, alias string) error {

	var result common.IssuesSearchResult

	var mainInfo common.MainTopom

	errDash := common.SearchIssues(conf.DashboardUrl, &mainInfo)
	if nil != errDash {
		log.Fatal("地址：", conf.DashboardUrl, "error:", errDash)
	}

	modes := mainInfo.Stats.Proxy.Models

	for i := range modes {
		var issueUrl = "http://" +  modes[i].AdminAddr + "/proxy/stats"
		err := common.SearchIssues(issueUrl, &result)
		if err != nil {
			log.Fatal("获取admin_addr地址失败：", err)
			break
		}
		e.metrics["dashboard_qps_total"].WithLabelValues(modes[i].AdminAddr, "test", "test").Set(result.Ops.Qps)
		e.metrics["dashboard_sessions_total"].WithLabelValues(modes[i].AdminAddr, "test", "test").Set(result.Sessions.Alive)
	}

	proxyModels := mainInfo.Stats.Group.ProxyModels

	for i := range proxyModels {

		servers := proxyModels[i].Servers
		for num := range servers {
			var resultInfo common.GroupInfo
			var addr = servers[num].Server
			var serverUrl = conf.CodisServerUrl + servers[num].Server


			if err := common.SearchIssues(serverUrl, &resultInfo); err != nil {
				log.Fatal("地址：", serverUrl, " error:", err)
				break
			}
			db0 := strings.Split(resultInfo.Db0, ",")
			db1 := strings.Split(resultInfo.Db1, ",")

			keysNum0 := strings.Split(db0[0], "=")
			keysNum1 := strings.Split(db1[0], "=")

			if keysNum0 == nil || keysNum1 == nil {
				log.Error("字符串转化错误, error:db0[0]:%s,db1[0]:%s", db0[0], db1[0])
				break
			}

			num0, err0 := strconv.ParseFloat(keysNum0[0], 64)
			num1, err1 := strconv.ParseFloat(keysNum1[0], 64)

			usedMemory, errM := strconv.ParseFloat(resultInfo.UsedMemory, 64)

			if  err0 != nil  || err1 != nil || errM != nil {
				log.Error("字符串转化错误, error:%s", errM, err1, err0 )
				break
			}

			// totalSystemMemory, _ := strconv.ParseFloat(resultInfo.TotalSystemMemory, 64)

			e.metrics["dashboard_keys_num"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id), "test").Set(num0 + num1)
			e.metrics["dashboard_memeory_percent"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id), "test").Set(usedMemory)

			hits, _ := strconv.ParseFloat(resultInfo.KeyspaceHits, 64)
			misses, _ := strconv.ParseFloat(resultInfo.KeyspaceMisses, 64)
			hitRatio := hits / (hits + misses)
			e.metrics["dashboard_cache_hit_rate"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id), "test").Set(hitRatio)

			usedCpuSys, _ := strconv.ParseFloat(resultInfo.UsedCpuSys, 64)
			usedCpuSysChildren, _ := strconv.ParseFloat(resultInfo.UsedCpuSysChildren, 64)

			e.metrics["dashboard_used_cpu_sys"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id)).Set(usedCpuSys)
			e.metrics["dashboard_used_cpu_sys_children"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id)).Set(usedCpuSysChildren)
			kbps, err := strconv.ParseFloat(resultInfo.InstantaneousOutputKbps,  64)
			if  err != nil {
				log.Error("字符串转化错误, error:%s", err)
				break
			}

			e.metrics["dashboard_instantaneous_output_kbps"].WithLabelValues(addr, strconv.Itoa(proxyModels[i].Id), "test").Set(kbps)

			e.metrics["dashboard_cache_hit_rate"].WithLabelValues(addr,  strconv.Itoa(proxyModels[i].Id), "test").Set(hitRatio)
		}
	}

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error("get memeory use percent error:%s", err)
	}

	usedPercent := v.UsedPercent

	e.metricsMtx.RLock()

	e.metrics["dashboard_memeory_percent"].WithLabelValues(addr, "test", "test").Set(usedPercent)

	for i := range result.Ops.Cmd {
		strCmd := result.Ops.Cmd[i].Opstr
		e.metrics["dashboard_usec_per_call"].WithLabelValues(addr, strCmd, "test").Set(result.Ops.Cmd[i].UsecsPercall)
	}
	e.metricsMtx.RUnlock()
	return nil
}

func NewCodisExporter(host RedisHost, namespace, checkSingleKeys, checkKeys string) (*Exporter, error) {

	e := Exporter{
		redis:     host,
		namespace: namespace,
		keyValues: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "key_value",
			Help:      "The value of \"key\"",
		}, []string{"addr", "alias", "db", "key"}),
		//
		diskQps: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "redis_qps",
			Help: "redis sum qps",
		}, []string{"diskQps"}),
		//
		keySizes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "key_size",
			Help:      "The length or size of \"key\"",
		}, []string{"addr", "alias", "db", "key"}),
		scriptValues: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "script_value",
			Help:      "Values returned by the collect script",
		}, []string{"addr", "alias", "key"}),
		duration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "exporter_last_scrape_duration_seconds",
			Help:      "The last scrape duration.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total redis scrapes.",
		}),
		scrapeErrors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "exporter_last_scrape_error",
			Help:      "The last scrape error status.",
		}),
	}

	var err error

	if e.keys, err = parseKeyArg(checkKeys); err != nil {
		return &e, fmt.Errorf("Couldn't parse check-keys: %#v", err)
	}
	log.Debugf("keys: %#v", e.keys)

	if e.singleKeys, err = parseKeyArg(checkSingleKeys); err != nil {
		return &e, fmt.Errorf("Couldn't parse check-single-keys: %#v", err)
	}
	log.Debugf("singleKeys: %#v", e.singleKeys)

	e.initCodisGauges()
	return &e, nil
}

func (e *Exporter) initCodisGauges() {

	e.metrics = map[string]*prometheus.GaugeVec{}
	e.metrics["instance_info"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "instance_info",
		Help:      "Information about the Redis instance",
	}, []string{"addr", "alias", "role", "redis_version", "redis_build_id", "redis_mode", "os"})
	e.metrics["slave_info"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "slave_info",
		Help:      "Information about the Redis slave",
	}, []string{"addr", "alias", "master_host", "master_port", "read_only"})
	e.metrics["start_time_seconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "start_time_seconds",
		Help:      "Start time of the Redis instance since unix epoch in seconds.",
	}, []string{"addr", "alias"})
	e.metrics["master_link_up"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name:      "master_link_up",
		Help:      "Master link status on Redis slave",
	}, []string{"addr", "alias"})

	// add ++
	e.metrics["slow_log_list"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "slow_log_list",
		Help: "redis slowlog get",
	}, []string{"addr", "args", "used_time", "execute"})

	e.metrics["dashboard_memeory_percent"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_memeory_percent",
		Help: "codis dashboard memeory percent",
	}, []string{"addr","args", "percent"}) // +

	e.metrics["dashboard_sessions_total"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_sessions_total",
		Help: "codis dashboard sessions total",
	}, []string{"addr", "alias", "percent"}) // +

	e.metrics["dashboard_qps_total"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dashboard_qps_total",
		Help: "codis dashboard qps total",
	}, []string{"addr", "alias", "percent"}) // +

	e.metrics["dashboard_keys_num"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_keys_num",
		Help: "codis dashboard keys num",
	}, []string{"addr", "value", "alias"}) // +

	e.metrics["dashboard_cache_hit_rate"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_cache_hit_rate",
		Help: "dashboard cache hit rate",
	}, []string{"addr", "value", "alias"}) // +

	e.metrics["dashboard_usec_per_call"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_usec_per_call",
		Help: "dashboard usec per call",
	}, []string{"addr", "cmd", "alias"})

	e.metrics["dashboard_instantaneous_output_kbps"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_instantaneous_output_kbps",
		Help: "codis dashboard instantaneous output kbps",
	}, []string{"addr", "alias", "kbps"})

	e.metrics["dashboard_used_cpu_sys"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_used_cpu_sys",
		Help: "dashboard used cpu sys",
	}, []string{"addr", "alias"})

	e.metrics["dashboard_used_cpu_sys_children"]= prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: e.namespace,
		Name: "dashboard_used_cpu_sys_children",
		Help: "codis used cpu sys children",
	}, []string{"addr", "alias"})

	// add++
	/*
		e.metrics["connected_slave_offset"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "connected_slave_offset",
			Help:      "Offset of connected slave",
		}, []string{"addr", "alias", "slave_ip", "slave_port", "slave_state"})
		e.metrics["connected_slave_lag_seconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "connected_slave_lag_seconds",
			Help:      "Lag of connected slave",
		}, []string{"addr", "alias", "slave_ip", "slave_port", "slave_state"})
		e.metrics["db_keys"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "db_keys",
			Help:      "Total number of keys by DB",
		}, []string{"addr", "alias", "db"})
		e.metrics["db_keys_expiring"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "db_keys_expiring",
			Help:      "Total number of expiring keys by DB",
		}, []string{"addr", "alias", "db"})
		e.metrics["db_avg_ttl_seconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "db_avg_ttl_seconds",
			Help:      "Avg TTL in seconds",
		}, []string{"addr", "alias", "db"})

		// Latency info
		e.metrics["latency_spike_last"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "latency_spike_last",
			Help:      "When the latency spike last occurred",
		}, []string{"addr", "alias", "event_name"})
		e.metrics["latency_spike_milliseconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "latency_spike_milliseconds",
			Help:      "Length of the last latency spike in milliseconds",
		}, []string{"addr", "alias", "event_name"})

		e.metrics["commands_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "commands_total",
			Help:      "Total number of calls per command",
		}, []string{"addr", "alias", "cmd"})
		e.metrics["commands_duration_seconds_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "commands_duration_seconds_total",
			Help:      "Total amount of time in seconds spent per command",
		}, []string{"addr", "alias", "cmd"})
		e.metrics["slowlog_length"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "slowlog_length",
			Help:      "Total slowlog",
		}, []string{"addr", "alias"})
		e.metrics["slowlog_last_id"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "slowlog_last_id",
			Help:      "Last id of slowlog",
		}, []string{"addr", "alias"})
		e.metrics["last_slow_execution_duration_seconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: e.namespace,
			Name:      "last_slow_execution_duration_seconds",
			Help:      "The amount of time needed for last slow execution, in seconds",
		}, []string{"addr", "alias"})
		*/
}