package contracts

type HostsStore interface {
	Init()
	Flush()
	Parse(content string)
	Query(IP string) (hosts []string)
	Bind(IP, host string)
	IPs() []string
	Forget(IP string)
	Count() int
}

type HostsConfigLoader interface {
	Path() (path string)
	Load(path string, hosts HostsStore)
	Print(hosts HostsStore)
}
