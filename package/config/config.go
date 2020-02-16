package config

type Config struct {
	Listen  Listen
	Cluster Cluster
}
type Listen struct {
	Port int
}
type Cluster struct {
	Broadcast ClusterBroadcast
}
type ClusterBroadcast struct {
	Ports []int
}
