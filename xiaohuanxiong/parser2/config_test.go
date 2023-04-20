package parser2

import (
	"bufio"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type I interface {
	ask()
}

type T struct {
}

func (i *T) ask() {
	fmt.Println("ask")
}

func (i *T) echo() {
	fmt.Println("echo")
}

func test(t I) {
	t.ask()
}

var tttt = `
{"#time":"2023-02-08 12:16:08","#event_name":"watermarginstrike","#account_id":"7027003538","#distinct_id":"oc_YT43J9NLZbRxjmzdFoA1YdXIU","#type":"track","properties":{"appkey":"1655274499097","version":"1.02.09","normversion":"v3","stepnumid":"C02170","serverid":"7027","gamechannel":"1018802001","userid":"oqDIlwrvLBsSkOd04V3fvrY8v0ak","rolelevel":90,"factionid":"1886","operatetype":1,"towerid":"1","floors":322,"result":"2","timey":"834","nimingid":"dh268ecd-04ce-005d-f4hc-86amhd72cbck-1554747","rolename":"回眸一笑百媚生\"","jielevel":"23","money100":334000,"pvefight_detail":[{"orderid":"341","heroid":"10014","rank":15,"hero_power":653083},{"orderid":"12186","heroid":"11002","rank":6,"hero_power":320086},{"orderid":"55","heroid":"10006","rank":15,"hero_power":643860},{"orderid":"1139","heroid":"10105","rank":13,"hero_power":398687},{"orderid":"5072","heroid":"10045","rank":12,"hero_power":496529}],"battletime":31,"yfrzbbh":"1.0.0","timex":"2023-02-08 12:16:08","towerauto":0}}
`

func TestConfig(t *testing.T) {

	//f, _ := os.Open("/Users/suxiangdong/code/logbus2_plugin/changyou/xiaohuanxiong/demolog/s.log")
	//buf := bufio.NewReader(f)
	//i := 0
	//for {
	//	l, _, e := buf.ReadLine()
	//	if e == io.EOF {
	//		break
	//	}
	//	ll := bytes.Split(l, []byte("\\u0001"))
	//	fmt.Println(string(ll[37]))
	//	i++
	//}
	//return
	//InitConfig()
	//fmt.Println(vip.Get("templates"))
	//return
	fmt.Println(InitConfig())
	//fmt.Println(cfg.Templates["s7.log"].Configs[0].Fields["best5"].Params)
	//return
	root := "/Users/suxiangdong/code/logbus2_plugin/changyou/xiaohuanxiong/demolog"
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fn := strings.Split(info.Name(), ".")
		fna := fn[0] + "." + fn[1]
		if fna == ".DS_Store" {
			return nil
		}

		if fna != "s.log" {
			return nil
		}

		//if fna == "gethero.log" || fna == "guildsnap.log" {
		//	return nil
		//}
		//if fna != "pvpfight.log" {
		//	return nil
		//}
		f, _ := os.Open(path)
		buf := bufio.NewReader(f)
		i := 0
		for {
			i++
			l, _, e := buf.ReadLine()
			//if i == 1 {
			//	continue
			//}
			if len(l) <= 0 {
				break
			}
			if e == io.EOF {
				break
			}
			str := fna + "|ta|" + string(l)
			bx, err := Parse([]byte(str))
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println(gjson.Valid(tttt))
			fmt.Println(gjson.ValidBytes(bx))
			fmt.Println(string(bx))
		}

		//fmt.Println(path)
		return nil
	})
	//r := bufio.NewReader(bytes.NewBuffer([]byte{}))
	//r.ReadLine()
	//
	////fmt.Println(gjson.ValidBytes([]byte(`{\"#distinct_id\":\"oc_YT42hfgi1Zk5sf5GmqkuLcB6Q\",\"#time\":\"2022-10-11 00:01:00\",\"#account_id\":\"7001002509\",\"#event_name\":\"pvpfight\",\"#type\":\"track\",\"properties\":{\"appkey\":\"1655274499097\",\"version\":\"1.00.63\",\"normversion\":\"v3\",\"stepnumid\":\"B8320\",\"serverid\":\"7001\",\"gamechannel\":\"2018832001\",\"userid\":\"oqDIlwjgRaQq4JLYkKubnDZDUEyQ\",\"rolelevel\":\"58\",\"factionid\":\"2\",\"power\":\"1133783\",\"rankbf\":\"17\",\"rankaf\":\"17\",\"battleid\":\"1\",\"targetroleid\":\"7001002235\",\"targetname\":\"核平东京\",\"targetlevel\":\"47\",\"targetpower\":\"451319\",\"targetrankbf\":\"21\",\"targetrankaf\":\"21\",\"result\":\"1\",\"time\":\"17\",\"rankscore\":\"0\",\"scorebf\":\"1301\",\"score\":\"1308\",\"targetscorebf\":\"1081\",\"targetscore\":\"1078\",\"platform\":\"4\",\"timey\":\"087\",\"nimingid\":\"3e37987g-7555-c4a2-8d49-d27ei77d020d-1339876\",\"rolename\":\"弑神\",\"jielevel\":\"9\",\"money100\":\"0\",\"own_group_detail\":[{\"orderid\":77,\"heroid\":\"10006\",\"rank\":7,\"hero_power\":82613},{\"orderid\":1648,\"heroid\":\"10059\",\"rank\":7,\"hero_power\":67306},{\"orderid\":933,\"heroid\":\"10026\",\"rank\":6,\"hero_power\":44675},{\"orderid\":119,\"heroid\":\"10045\",\"rank\":7,\"hero_power\":77355},{\"orderid\":1468,\"heroid\":\"10035\",\"rank\":7,\"hero_power\":50800}],\"rival_group_detail\":[{\"orderid\":56,\"heroid\":\"10006\",\"rank\":6,\"hero_power\":30785},{\"orderid\":107,\"heroid\":\"10045\",\"rank\":5,\"hero_power\":25418},{\"orderid\":27,\"heroid\":\"10013\",\"rank\":6,\"hero_power\":51395},{\"orderid\":299,\"heroid\":\"10014\",\"rank\":5,\"hero_power\":23826},{\"orderid\":0,\"heroid\":\"-1\",\"rank\":0,\"hero_power\":0}],\"yfrzbbh\":\"1.0.0\",\"timex\":\"2022-10-11 00:01:00\",\"battletype\":\"1\",\"lszlscore\":\"null\n\"}}`)))
	//fmt.Println(`{"a":"b"}`)
	//fmt.Println(InitConfig())
	////fmt.Println(Init())
	//var b = []byte("err.log||2022-10-11 00:01:0016552744990971.00.63pvpfightv3B832070012018832001oqDIlwjgRaQq4JLYkKubnDZDUEyQ7001002509582oc_YT42hfgi1Zk5sf5GmqkuLcB6Q1133783171717001002235核平东京4745131921211170130113081081107840873e37987g-7555-c4a2-8d49-d27ei77d020d-1339876弑神9077|10006|7|82613,1648|10059|7|67306,933|10026|6|44675,119|10045|7|77355,1468|10035|7|5080056|10006|6|30785,107|10045|5|25418,27|10013|6|51395,299|10014|5|23826,-11.0.02022-10-11 00:01:001null\n::err.log||2022-10-11 00:01:00\u00011655274499097\u00011.00.63\u0001pvpfight\u0001v3\u0001B8320\u00017001\u00012018832001\u0001oqDIlwjgRaQq4JLYkKubnDZDUEyQ\u00017001002509\u000158\u00012\u0001oc_YT42hfgi1Zk5sf5GmqkuLcB6Q\u00011133783\u000117\u000117\u00011\u00017001002235\u0001核平东京\u000147\u0001451319\u000121\u000121\u00011\u000117\u00010\u00011301\u00011308\u00011081\u00011078\u00014\u0001087\u00013e37987g-7555-c4a2-8d49-d27ei77d020d-1339876\u0001弑神\u00019\u00010\u000177|10006|7|82613,1648|10059|7|67306,933|10026|6|44675,119|10045|7|77355,1468|10035|7|50800\u000156|10006|6|30785,107|10045|5|25418,27|10013|6|51395,299|10014|5|23826,-1\u00011.0.0\u00012022-10-11 00:01:00\u00011\u0001null\n")
	////var b = []byte("err.log||2022-10-11 00:01:00\u00011655274499097\u00011.00.63\u0001pvpfight\u0001v3\u0001B8320\u00017001\u00012018832001\u0001oqDIlwjgRaQq4JLYkKubnDZDUEyQ\u00017001002509\u000158\u00012\u0001oc_YT42hfgi1Zk5sf5GmqkuLcB6Q\u00011133783\u000117\u000117\u00011\u00017001002235\u0001核平东京\u000147\u0001451319\u000121\u000121\u00011\u000117\u00010\u00011301\u00011308\u00011081\u00011078\u00014\u0001087\u00013e37987g-7555-c4a2-8d49-d27ei77d020d-1339876\u0001弑神\u00019\u00010\u000177|10006|7|82613,1648|10059|7|67306,933|10026|6|44675,119|10045|7|77355,1468|10035|7|50800\u000156|10006|6|30785,107|10045|5|25418,27|10013|6|51395,299|10014|5|23826,-1\u00011.0.0\u00012022-10-11 00:01:00\u00011\u0001null\n::err.log||2022-10-11 00:02:01\u00011655274499097\u00011.00.63\u0001pvpfight\u0001v3\u0001B8320\u00017001\u00012018832001\u0001oqDIlwjgRaQq4JLYkKubnDZDUEyQ\u00017001002509\u000158\u00012\u0001oc_YT42hfgi1Zk5sf5GmqkuLcB6Q\u00011133783\u00010\u00010\u00011\u00017001000008\u0001隔壁家二叔\u000160\u00011059351\u00010\u00010\u00010\u000123\u00013\u000144841\u000144841\u000193410\u000193410\u00014\u0001088\u00013e37987g-7555-c4a2-8d49-d27ei77d020d-1339876\u0001弑神\u00019\u00010\u0001541|10014|5|38646,434|10027|4|35864,1111|10011|4|35392,2093|11003|5|39995,574|10037|5|38134_77|10006|7|82613,1648|10059|7|67306,637|10105|5|39645,526|10056|6|45936,27|10013|6|49314_1175|10023|6|45826,2669|10017|6|44558,933|10026|6|44675,1468|10035|7|50800,119|10045|7|77355\u00012122|10045|7|75146,132|10014|8|82026,2874|10037|7|44176,-1,27|10013|4|7500_62|10006|10|125001,899|10103|6|37307,1364|10026|6|37152,887|10020|4|31646,433|10029|4|7482_410|10017|6|70167,1110|10105|6|58177,-1,1256|10036|7|42464,-1\u00011.0.0\u00012022-10-11 00:02:01\u00012\u0001null\n")
	//bx, err := Parse(b)
	//fmt.Println(string(bx))
	//fmt.Println(err)
}
