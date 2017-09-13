package influxdb1xx
import (
//	"github.com/efigence/influxdb-backup/common"
//	"github.com/efigence/influxdb-backup/input"
	influx "github.com/influxdata/influxdb/client/v2"
	"fmt"
	"time"
	"strings"
)


type Input struct {
	client influx.Client
	serverVersion string
	db string
}

func NewInput(addr string, user string, pass string, db string) (i *Input, err error) {
	i = &Input{
		db: db,
	}
	i.client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr: addr,
		Username: user,
		Password: pass,
	})
	if err != nil {return i ,err}
	_, i.serverVersion, err = i.client.Ping(time.Duration(time.Second*10))
	return i, err
}

func (i *Input) GetMeasurements() ([]string,error) {
	q := influx.NewQuery("SHOW MEASUREMENTS", i.db, "ns")
	resp, err := i.client.Query(q)
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	measurements := make([]string,len( resp.Results[0].Series[0].Values))
	for id,res :=  range resp.Results[0].Series[0].Values {
		measurements[id] = res[0].(string)
	}
	return measurements,nil
}

func PrintDebugResults (r *influx.Result) {
	if len(r.Messages) > 0 {
		fmt.Printf("Messages: %+v\n",r.Messages)
	}
	for id, ser := range r.Series {
		fmt.Printf("Series name [%s], ID: %d\n",ser.Name, id)
		fmt.Print("  Tags:\n")
		for k, v := range ser.Tags {
			fmt.Printf("    %s => %s\n",k,v)
		}
		fmt.Printf("%s\n",strings.Join(ser.Columns,"\t\t"))
		for id2, val := range ser.Values {
			fmt.Printf("%d -> %+v\n",id2, val)
		}

	}

}
