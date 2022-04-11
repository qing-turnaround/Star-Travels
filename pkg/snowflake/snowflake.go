package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int) (err error) {
	var st time.Time
	//时间为UTC时间，比中国慢8个小时
	st, err = time.Parse("2006-01-02 00:00:00", startTime)
	fmt.Println()
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(int64(machineID))
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
