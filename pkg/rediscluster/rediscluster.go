package rediscluster

func NewRedisClient() *redis.Client {
	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"}, // Redis Cluster 节点地址
	})
}
