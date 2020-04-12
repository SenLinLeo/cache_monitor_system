package conf

var CodisServerUrl    = "http://172.30.2.181:18080/api/topom/group/info/"
//var CodisServerUrl = "http://172.30.2.181:18080/api/topom/group/info/172.30.2.182:6380"
var DashboardUrl      = "http://172.30.2.181:18080/topom"

var MetricPath string = "/metrics"
var Addrs             = []string{"172.30.2.182:6380"}
var Passwords         = []string{"iboxpay"}
var Aliases           = []string{"test"}


var Namespace         = "redis"
var CheckKeys         = "redis"
var CheckSingleKeys   = "redis"

/**  1 -- Codis
**   2 -- Redis
**/
var Flag              = 1