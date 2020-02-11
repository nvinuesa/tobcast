package config

type Config struct {
	Listen  Listen
	Cluster Cluster
}
type Listen struct {
	Broadcast Broadcast
	Deliver   Deliver
}
type Cluster struct {
	Broadcast ClusterBroadcast
	Deliver   ClusterDeliver
}
type Broadcast struct {
	Port int
}
type Deliver struct {
	Port int
}
type ClusterBroadcast struct {
	Ports []int
}
type ClusterDeliver struct {
	Ports []int
}
