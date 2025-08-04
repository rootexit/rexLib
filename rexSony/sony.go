package rexSony

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
)

// note: 为了保持旧兼容
func NextId() string {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := sf.NextID()
	return fmt.Sprintf("%d", id)
}

// note: 索尼基于雪花算法实现的变种
// note: https://github.com/sony/sonyflake
func SonyNextId(node uint16) string {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return node, nil
		},
	})
	id, _ := sf.NextID()
	return fmt.Sprintf("%d", id)
}

// note: Twitter雪花算法的golang实现
// note: https://github.com/bwmarrin/snowflake
func TwitterSony(nodeCode int64) (string, error) {
	node, err := snowflake.NewNode(nodeCode)
	if err != nil {
		return "", err
	}
	return node.Generate().String(), nil
}

func TwitterSonyNode(nodeCode int64) (*snowflake.Node, error) {
	node, err := snowflake.NewNode(nodeCode)
	if err != nil {
		return nil, err
	}
	return node, nil
}
