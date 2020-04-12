package conf

var CodisServerUrl    = "http://192.168.5.181:18080/api/topom/group/info/"
var DashboardUrl      = "http://192.168.5.181:18080/topom"

var MetricPath string = "/metrics"
var Addrs             = []string{"192.168.5.182:6380"}
var Passwords         = []string{"passwad"}
var Aliases           = []string{"test"}


var Namespace         = "redis"
var CheckKeys         = "redis"
var CheckSingleKeys   = "redis"

/**  1 -- Codis
**   2 -- Redis
**/
var Flag              = 1
