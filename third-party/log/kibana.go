package log

const (
	CsTypeShow  = "show"
	CsTypeDebug = "debug"
)

// log a tps tag
func TagShowTps(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "tps", "env", "filebeat")
	Agent(msg, ctx...)
}

func TagShowIPs(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "ip", "env", "filebeat")
	Agent(msg, ctx...)
}

// tx pool
func TagShowTxPool(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "tx_pool")
	Agent(msg, ctx...)
}

// node location
func TagShowNodeLocation(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "node_location")
	Agent(msg, ctx...)
}

// 双花log
func TagShowDoubleExpend(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "double_expend")
	Agent(msg, ctx...)
}

// 女巫攻击
func TagWitchAttack(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeShow, "cs-item", "witch_attack")
	Agent(msg, ctx...)
}

// debug error logs
func TagError(msg string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeDebug, "cs-item", "error")
	Agent(msg, ctx...)
}

// trace feature states
func TagStateTrace(msg string, traceType string, ctx ...interface{}) {
	ctx = append(ctx, "cs-type", CsTypeDebug, "cs-item", "state", "trace_type", traceType)
	Agent(msg, ctx...)
}
