package snowflake

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeProvider struct {
	node *snowflake.Node
}

func (sp SnowflakeProvider) GenerateSnowflake() string {
	return strconv.FormatInt(sp.node.Generate().Int64(), 10)
}

// NewSnowflakeProvider generates a new snowflake ID service.
func NewSnowflakeProvider() (SnowflakeProvider, error) {
	sp := SnowflakeProvider{}
	node, err := snowflake.NewNode(1)

	if err != nil {
		return sp, err
	}

	sp.node = node
	return sp, nil
}
