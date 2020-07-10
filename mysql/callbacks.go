package mysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	QuerySTKEY = "querystkey"
)

func (c *Client) RegisterMetricsCallbacks(client *gorm.DB) {
	client.Callback().Create().After("gorm:after_create").Register("kjlibs:metrics_after_create", func(scope *gorm.Scope) {
		if scope.HasError() {
			c.handleError(scope)
		}
	})

	client.Callback().Update().After("gorm:after_update").Register("kjlibs:metrics_after_update", func(scope *gorm.Scope) {
		if scope.HasError() {
			c.handleError(scope)
		}
	})

	client.Callback().Query().Before("gorm:query").Register("qms:before_query", func(scope *gorm.Scope) {
		scope.InstanceSet(QuerySTKEY, time.Now())
	})
	client.Callback().Query().After("gorm:after_query").Register("kjlibs:metrics_after_query", func(scope *gorm.Scope) {
		if scope.HasError() {
			c.handleError(scope)
		}
	})

	client.Callback().Delete().After("gorm:after_delete").Register("kjlibs:metrics_after_delete", func(scope *gorm.Scope) {
		if scope.HasError() {
			c.handleError(scope)
		}
	})

	client.Callback().RowQuery().After("gorm:row_query").Register("kjlibs:metrics_after_sql", func(scope *gorm.Scope) {
		if scope.HasError() {
			c.handleError(scope)
		}
	})
}

func (c *Client) handleError(scope *gorm.Scope) {
	var (
		schema = scope.Dialect().CurrentDatabase()
		fields = logrus.Fields{
			"schema":   schema,
			"sql_desc": fmt.Sprintf("%v", scope.DB().QueryExpr()),
		}
	)

	if err := scope.DB().Error; err == gorm.ErrRecordNotFound {
		// db miss
		c.log.WithFields(fields).WithError(err).Error("mysql miss")
		DBMiss.With(prometheus.Labels{"schema": schema, "table": scope.QuotedTableName()}).Inc()
	} else {
		// db error
		c.log.WithFields(fields).WithError(err).Error("mysql error")
		DBError.With(prometheus.Labels{"schema": schema, "table": scope.QuotedTableName()}).Inc()
	}
}
