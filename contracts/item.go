//go:generate mockgen -package mocks -destination ../mocks/item.go . Timestamped

package contracts

type Timestamped interface {
	GetTimestamp() int64
}
