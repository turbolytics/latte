package partition

import "net/url"

type Partition interface {
	URI() *url.URL
}
