package configuration

import "time"

type ConfigurationService interface {
	Read() error
	Load() Loader
}

type Loader interface {
	Get(string) interface{}
	GetBool(string) bool
	GetDuration(string) time.Duration
	GetFloat64(string) float64
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetIntSlice(string) []int
	GetSizeInBytes(string) uint
	GetString(string) string
	GetStringMap(string) map[string]interface{}
	GetStringMapString(string) map[string]string
	GetStringMapStringSlice(string) map[string][]string
	GetStringSlice(string) []string
	GetTime(string) time.Time
	GetUint(string) uint
	GetUint32(string) uint32
	GetUint64(string) uint64
}
