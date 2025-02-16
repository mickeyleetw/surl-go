package dbcmd

import (
	"context"
	database "shorten_url/pkg/core/database"
	"shorten_url/pkg/core/redis"

	"github.com/spf13/cobra"
)

// ResetDBCmd is the command to reset the database
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

// ResetRedisCmd is the command to reset the redis
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

// resetDB is the function to reset the database
func resetDB() {
	database.GetDB()
}

// resetRedis is the function to reset the redis
func resetRedis() {
	redisClient := redis.GetRedis()
	redisClient.FlushAll(context.Background())
}
