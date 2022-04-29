package redis

import (
	"hash/crc32"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"web_app/settings"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}


type ConsistentHashBalance struct {
	replicas int // 虚拟节点个数
	keys     UInt32Slice       // 已排序的节点hash切片（方便排序查找节点）
	Hash // 哈希算法
	HashMap map[uint32]string
	Clients map[string][]*redis.Conn // 取出连接
	MasterClients map[string]*redis.Client // sentinel 连接master，支持更多操作
	SentinelClient *redis.SentinelClient // 用于查询redis masterIP
	Password map[string]string // 对应masterName的密码
}


func NewConsistentHashBalance(hash Hash, replicas int, cfg *settings.RedisConfig) *ConsistentHashBalance {
	c := &ConsistentHashBalance{
		replicas: replicas,
		keys: UInt32Slice{},
		Hash:     hash,
		HashMap:        make(map[uint32]string),
		Clients:  make(map[string][]*redis.Conn),
		MasterClients: make(map[string]*redis.Client),
		SentinelClient: nil,
		Password: map[string]string{},
	}

	if c.Hash == nil {
		c.Hash = crc32.ChecksumIEEE // 设置默认hash算法
	}
	// 添加节点
	for i := 0; i < cfg.Masters.Counts; i++ {
		// 添加 master -> password的映射
		c.Password[cfg.Masters.MasterName[i]] = cfg.Masters.Passwords[i]
		config := &SentinelConfig{
			SentinelAddrs:    cfg.Sentinels,
			MasterName:       cfg.Masters.MasterName[i],
			SentinelPassword: cfg.Password,
			Password:         cfg.Masters.Passwords[i],
			PoolSize:         cfg.PoolSize,
			MaxOpenConns: cfg.MaxOpenConns,
		}
		c.AddMaster(config)
	}
	// 添加一个主sentinel
	c.SentinelClient = redis.NewSentinelClient(&redis.Options{
		Addr: cfg.Sentinels[0], // 随机取一个
		Password: cfg.Password,
	})

	return c
}


// 增加一个 master节点
func (c *ConsistentHashBalance) AddMaster(conf *SentinelConfig) (err error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       conf.MasterName,
		SentinelAddrs:    conf.SentinelAddrs,
		SentinelPassword: conf.SentinelPassword,
		PoolSize: conf.PoolSize,
		Password: conf.Password,
	})

	var clients []*redis.Conn
	// 每个服务端 maxOpenConns 个连接
	for i := 0; i < conf.MaxOpenConns; i++ {
		clients = append(clients, client.Conn(ctx))
	}
	// 将对应的 masterName（因为主从可能切换，但是MasterName一般不会进行改变）与client相关联
	c.Clients[conf.MasterName] = clients
	c.MasterClients[conf.MasterName] = client

	// 结合复制因子计算所有虚拟节点的hash值，并存入m.keys中，同时在m.hashMap中保存哈希值和key的映射
	for i := 0; i < c.replicas; i++ {
		hash := c.Hash([]byte(strconv.Itoa(i) + conf.MasterName))
		c.keys = append(c.keys, hash) // 方便排序查找节点
		c.HashMap[hash] = conf.MasterName
	}

	sort.Sort(c.keys)
	return nil
}

// 更新节点
func (c *ConsistentHashBalance) Update(conf *settings.RedisConfig) {
	// 新的 rb
	newRb := NewConsistentHashBalance(nil, conf.Replicas, conf)
	// 迁移数据
	c.Migrate(newRb)
	// 给全局变量赋予新值
	rb = newRb
}

// 迁移数据
func (c *ConsistentHashBalance) Migrate(newC *ConsistentHashBalance) {
	for masterName, client := range c.MasterClients {

		// 取出所有key,执行keys *
		keysResult, _ := client.Do(ctx, "key", "*").StringSlice()

		// 流水线
		pipeline := client.TxPipeline()
		// 循环执行迁移数据
		for _, key := range keysResult {
			// 得出key应该部署在哪一个节点
			newMasterName := newC.Get(key)
			if newMasterName == masterName { // 还是原来的节点，跳过该键
				continue
			}

			// 得出 newMasterName对应节点，并取出ip，port
			addr, port := c.GetHost(newMasterName)
			// 进行迁移
			_, _ = pipeline.Do(ctx, "migrate", addr, port,
				"", "0", "500", "replace", "auth", newC.Password[newMasterName],
				"keys", key).Result()
		}
		pipeline.Exec(ctx)
	}
	zap.L().Info("数据迁移成功！")
}

// Get 方法根据给定的对象获取最靠近它的那个节点
func (c *ConsistentHashBalance) Get(key string) string {
	// 计算 key的哈希值
	hashVal := c.Hash([]byte(key))

	// 通过二分查找获取最优节点，第一个"服务器hash"值大于"数据hash"值的就是最优"服务器节点"
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hashVal
	})

	// 如果查找结果 大于 服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(c.keys) {
		idx = 0
	}

	// 返回 对应的服务器节点的 masterName
	return c.HashMap[c.keys[idx]]
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) GetRandomConn(masterName string) *redis.Conn {
	length := len(c.Clients[masterName])
	// 随机取出一个连接
	return c.Clients[masterName][rand.Intn(length)]
}

type SentinelConfig struct {
	SentinelAddrs []string
	MasterName  string
	// sentinel Password
	SentinelPassword string
	// master Password
	Password string
	PoolSize int
	MaxOpenConns int
}

// // 利用正则表达式从命令 role 结果中获取 master addr 和 port
func (c *ConsistentHashBalance)GetHost(masterName string) (string, string) {
	// 通过sentinel来获取ip
	result := c.SentinelClient.GetMasterAddrByName(ctx, masterName).String()
	pattern := "(2[0-4]\\d|25[0-5]|[01]?\\d\\d?\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?.(\\d?){5})"
	host := regexp.MustCompile(pattern).FindString(result)
	newHost := strings.Split(host, " ")
	return newHost[0], newHost[1]
}