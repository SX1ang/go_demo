package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) {
	var st time.Time
	var err error

	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	snowflake.Epoch = st.Unix() / 1000000

	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return
	}
}

func GenID() int64 {
	return node.Generate().Int64()
}
