//go:generate mockgen -source=provider.go -destination=mocks/mock.go
package hash

type HashProvider interface {
	Make(string) string
	Compare(string, string) error
}
