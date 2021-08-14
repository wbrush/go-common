package db

import (
	"context"
	"github.com/fatih/color"
	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type LoggerHook struct{}

func (d LoggerHook) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d LoggerHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	//this check is used for not to unnecessary vars, if not needed
	if logrus.GetLevel() == logrus.TraceLevel {
		yellow := color.New(color.FgYellow).SprintFunc()
		q, _ := qe.FormattedQuery()
		logrus.Tracef("%s", yellow(q))
	}
	return nil
}

type RedisLoggerHook struct{}

func (rlh RedisLoggerHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (rlh RedisLoggerHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	//this check is used for not to unnecessary vars, if not needed
	if logrus.GetLevel() == logrus.TraceLevel {
		yellow := color.New(color.FgYellow).SprintFunc()
		logrus.Tracef("%s", yellow(cmd.String()))
	}

	return nil
}

func (rlh RedisLoggerHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (rlh RedisLoggerHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	//this check is used for not to unnecessary vars, if not needed
	if logrus.GetLevel() == logrus.TraceLevel {
		yellow := color.New(color.FgYellow).SprintFunc()
		for i := range cmds {
			logrus.Tracef("%s", yellow(cmds[i].String()))
		}
	}

	return nil
}
