package sqlite

import (
	//	"sync"
	"fmt"
	"github.com/efigence/influxdb-backup/src/common"
)

type router struct {
	routingTable map[string]chan *common.Field

}

func (s *SQLiteOut) route(in chan *common.Field){
	var r router
	r.routingTable = make(map[string]chan *common.Field)
	for field := range in {
		routingKey := s.SeriesNameGen(field.Name)
		if ch, ok := r.routingTable[routingKey]; ok {
			ch <- field
		} else {
			ch, err := s.workers.GetRouteFor(routingKey)
			if (err != nil) {
				panic(fmt.Sprintf("Err when getting route for %s: %s", routingKey,err))
			}
			r.routingTable[routingKey] = ch
			ch <- field
		}
	}
}
