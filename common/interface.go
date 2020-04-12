package common

/** Codis dashboard 接口格式 结构体 **/
type IssuesSearchResult struct {
	Online    bool       `json:"online"`
	Closed    bool       `json:"closed"`
	Sentinels *Sentinels `json:"sentinels"`
	Ops       *Ops       `json:"ops"`
	Sessions  *Sessions  `json:"sessions"`
	Rusage    *Rusage    `json:"rusage"`
	Backend   *Backend   `json:"backend"`
	Runtime   *Runtime   `json:"runtime"`
}

type Runtime struct {
	General        *General  `json:"general"`
	Heap           *Heap     `json:"heap"`
	Gc             *Gc       `json:"gc"`
	NumProcs       int       `json:"num_procs"`
	NumGoroutines  int       `json:"num_procs"`
	NumCgoCall     int       `json:"num_cgo_call"`
	MemOffHeap     int       `json:"mem_offheap"`
}

type Gc struct {
	Num           int     `json:"num"`
	CpuFraction  float32 `json:"cpu_fraction"`
	TotalPausems int     `json:"total_pausems"`
}
type General struct {
	Alloc   int `json:"alloc"`
	Sys     int `json:"sys"`
	Lookups int `json:"lookups"`
	Mallocs int `json:"mallocs"`
	Frees   int `json:"frees"`
}
type Heap struct {
	Alloc   int `json:"alloc"`
	Sys     int `json:"sys"`
	Idle    int `json:"idle"`
	Inuse   int `json:"inuse"`
	Objects int `json:"objects"`
}
type Backend struct {
	PrimaryOnly bool `json:"primary_only"`
}

type Rusage struct {
	Now string  `json:"now"`
	Cpu float64 `json:"cpu"`
	Mem int     `json:"mem"`
	Raw *Raw    `json:"raw"`
}

type Raw struct {
	Utime       int `json:"utime"`
	Stime       int `json:"stime"`
	Cutime      int `json:"cutime"`
	Cstime      int `json:"cstime"`
	NumThreads int `json:"num_threads"`
	VmSize     int `json:"vm_size"`
	VmRss      int `json:"vm_rss"`
}
type Sentinels struct {
	Servers []string          `json:"servers"`
	Masters map[string]string `json:"masters"`
}

type CmdInfo struct {
	Opstr           string  `json:"opstr"`
	Calls           int     `json:"calls"`
	Usecs           int     `json:"usecs"`
	UsecsPercall   float64     `json:"usecs_percall"`
	Fails           int     `json:"fails"`
	RedisErrtype   int     `json:"redis_errtype"`

}

type Ops struct {
	Total   float64                      `json:"total"`
	Fails   float64                      `json:"fails"`
	Qps     float64                      `json:"qps"`
	Cmd     []CmdInfo                    `json:"cmd"`
	Redis   map[string]int               `json:"redis"`
	Servers []string
}

type Sessions struct {
	Total float64 `json:"total"`
	Alive float64     `json:"alive"`
}

/** api/topom/group/info  **/
type GroupInfo struct {
	Db0                       string   `json:"db0"`
	Db1                       string   `json:"db2"`
	UsedMemory                string   `json:"used_memory"`
	UsedMemoryLua             string   `json:"used_memory_lua"`
	KeyspaceHits              string   `json:"keyspace_hits"`
	KeyspaceMisses            string   `json:"keyspace_misses"`
	InstantaneousOutputKbps   string  `json:"instantaneous_output_kbps"`
	TotalSystemMemory		  string  `json:"total_system_memory"`
	UsedCpuSys		          string  `json:"used_cpu_sys"`
	UsedCpuSysChildren	      string  `json:"used_cpu_sys_children"`

}

type KeysDb struct {
	Keys     float64 `json:"keys"`
	Expires  float64 `json:"expires"`
	AvgTtl  float64 `json:"avg_ttl"`
}


/*** topom 主页面信息   ***/
type MainTopom struct {
	Version  string `json:"version"`
	Stats    *Stats `json:"stats"`
}

type Stats struct {
	Proxy *Proxy  `json:"proxy"`
	Group *Group  `json:"group"`
}

type Models struct {
	ID         int      `json:"id"`
	AdminAddr  string   `json:"admin_addr"`
	ProxyAddr  string   `json:"proxy_addr"`
}

/*  proxy 信息 */
type Proxy struct {
	Models []Models   `json:"models"`
}

type ProxyModels struct {
	Id           int        `json:"id"`
	Servers     []Servers   `json:"servers"`
}

type Servers struct {
	Server        string   `json:"server"`
	//Datacenter    string   `json:"datacenter"`
	//ReplicaGroup  bool     `json:"replica_group"`
}

/** 分组信息 */
type Group struct {
	ProxyModels []ProxyModels   `json:"models"`
}