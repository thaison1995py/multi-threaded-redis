package core

import "github.com/thaison199py/multi-threaded-redis/internal/data_structure"

var dictStore *data_structure.Dict
var setStore map[string]*data_structure.SimpleSet
var zsetStore map[string]*data_structure.SortedSet
var cmsStore map[string]*data_structure.CMS
var bloomStore map[string]*data_structure.Bloom

func init() {
	dictStore = data_structure.CreateDict()
	setStore = make(map[string]*data_structure.SimpleSet)
	zsetStore = make(map[string]*data_structure.SortedSet)
	cmsStore = make(map[string]*data_structure.CMS)
	bloomStore = make(map[string]*data_structure.Bloom)
}
