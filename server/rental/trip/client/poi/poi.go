package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"hash/fnv"

	"github.com/golang/protobuf/proto"
)

var poi = []string{
	"中关村",
	"天安门",
	"陆家嘴",
	"迪士尼",
	"天河体育中心",
	"广州塔",
}

// Manager defines a poi manager.
type Manager struct {
}

// Resolve resolves the given location.
func (*Manager) Resolve(c context.Context, loc *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(loc)
	if err != nil {
		return "", err
	}
	// 哈希算法 给定同一个地点，只会算出一个POI的值
	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
