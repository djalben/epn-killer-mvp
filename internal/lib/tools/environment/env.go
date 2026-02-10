package environment

import (
	"os"
	"strconv"
)

const (
	envListenPortHTTPMetric = "LISTEN_PORT_HTTP_METRICS"
	envLogLevel             = "LOG_LEVEL"
	envLogPlain             = "LOG_PLAIN"

	// PostgreSQL.
	envPostgresDSN = "POSTGRES_DSN"
)

// GetListenPortHTTPMetric возвращает порт для HTTP метрик из переменной окружения LISTEN_PORT_HTTP_METRIC.
// Если переменная окружения не установлена, возвращает 0.
// Если переменная окружения установлена, возвращает значение, преобразованное в uint16.
// Пример использования:
// port := environment.GetListenPortHTTPMetric()
//
//	if port != 0 {
//	    fmt.Println("HTTP Metric Port:", port)
//	} else {
//
//	    fmt.Println("HTTP Metric Port is not set")
//	}
func GetListenPortHTTPMetric() uint16 {
	portStr, exists := os.LookupEnv(envListenPortHTTPMetric)
	if !exists {
		return 0
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return 0
	}

	return uint16(port)
}

// GetLog возвращает уровень логирования и флаг, указывающий на использование плоского формата логов.
// Если переменная окружения LOG_LEVEL не установлена, возвращает "debug2".
// Если переменная окружения LOG_PLAIN установлена, возвращает значение, преобразованное в bool.
// Пример использования:
// level, isPlain := environment.GetLog()
//
//	if isPlain {
//	    fmt.Println("Logging in plain format at level:", level)
//	} else {
//
//	    fmt.Println("Logging in structured format at level:", level)
//	}
func GetLog() (level string, isPlain bool) {
	level, exists := os.LookupEnv(envLogLevel)
	if !exists {
		level = "debug2"
	}

	isPlainStr, exists := os.LookupEnv(envLogPlain)
	if exists {
		isPlain, _ = strconv.ParseBool(isPlainStr)
	}

	return level, isPlain
}
