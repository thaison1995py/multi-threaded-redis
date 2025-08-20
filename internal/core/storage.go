package core

import "github.com/thaison199py/multi-threaded-redis/internal/data_structure"

var dictStore *data_structure.Dict

func init() {
	dictStore = data_structure.CreateDict()
}
