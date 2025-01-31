package dbCmd

import (
	"context"
	database "shorten_url/pkg/core/database"
	"shorten_url/pkg/core/redis"

	"github.com/spf13/cobra"
)

var (
	ResetDBCmd = &cobra.Command{
		Use:   "resetdb",
		Short: "Database commands",
		Long:  "Database commands",
		Run: func(cmd *cobra.Command, args []string) {
			resetDB()
		},
	}
)

var (
	ResetRedisCmd = &cobra.Command{
		Use:   "resetredis",
		Short: "Redis commands",
		Long:  "Redis commands",
		Run: func(cmd *cobra.Command, args []string) {
			resetRedis()
		},
	}
)

func resetDB() {
	database.GetDB()
}

func resetRedis() {
	redis_client := redis.GetRedis()
	redis_client.FlushAll(context.Background())
}
