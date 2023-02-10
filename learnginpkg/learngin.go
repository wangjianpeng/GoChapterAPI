package learnginpkg

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ChapterResp struct {
	Error_Code int               `json:"error_code"`
	Error_Msg  string            `json:"error_msg"`
	Error_Ver  int               `json:"error_ver"`
	Data       map[string]string `json:"data"`
}

type FakeReq struct {
	Action string `json:"action"`
	Msg    string `json:"msg"`
}

func DoPingGin() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		var tempstr = `"message": "pong"`
		fmt.Println(tempstr)
		var in bytes.Buffer
		w := zlib.NewWriter(&in)
		w.Write([]byte(tempstr))
		w.Close()
		str := base64.StdEncoding.EncodeToString(in.Bytes())
		c.JSON(http.StatusOK, str)
	})

	r.POST("/fakepost", func(c *gin.Context) {

		bodydata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err.Error())
		}
		//base64 decode

		//unzip

		rawDecode1, err := base64.StdEncoding.DecodeString(string(bodydata))
		if err != nil {
			log.Println(err.Error())
		}

		rawDecode1bytes := bytes.NewReader(rawDecode1)
		reqbodyreader, err := zlib.NewReader(rawDecode1bytes)
		if err != nil {
			log.Println(err.Error())
		}

		reqbodydata, err := ioutil.ReadAll(reqbodyreader)
		if err != nil {
			log.Println(err.Error())
		}

		var fakeReq FakeReq
		err = json.Unmarshal(reqbodydata, &fakeReq)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(fakeReq.Action, "\t", fakeReq.Msg)

		time.Sleep(10 * time.Second)
		//strings.Contains(fakeReq.Msg, "ReadChapterPan") ||
		//|| strings.Contains(fakeReq.Msg, "node callback")
		// if strings.Contains(fakeReq.Msg, "shop frame") {
		// 	c.String(http.StatusRequestTimeout, "Time Out!")
		// 	return
		// }
		tempData := map[string]string{}

		tempData["name"] = "wjp"
		tempData["age"] = "32"
		tempChapterResp := &ChapterResp{
			Error_Code: 1,
			Error_Msg:  "success",
			Error_Ver:  123456,
			Data:       tempData,
		}
		rtbytes, _ := json.Marshal(&tempChapterResp)

		var in bytes.Buffer
		w := zlib.NewWriter(&in)
		w.Write(rtbytes)
		w.Close()
		str := base64.StdEncoding.EncodeToString(in.Bytes())
		// fmt.Println(str)
		// c.JSON(http.StatusOK, []byte(str))
		c.Data(http.StatusOK, "application/json", []byte(str))
		// c.Data(http.StatusOK, gin.MIMEJSON, []byte(str))
		// c.String(http.StatusOK, "hello world!")

	})

	r.Run("192.168.12.57:9999")
}

func DoBuildChapterResponseText() {
	inputstr := `{
		"user": {
			"_id": {
				"$id": "5e588f478295eced0d8b4567"
			},
			"uuid": "50831",
			"device_platform": "android",
			"channel_id": "AVG50005",
			"device_id": "008796764013583",
			"account_type": "tourists",
			"account_id": "tourists",
			"app_version": "1.6.9",
			"base_code_version": "1.6.9",
			"code_version": "19.0.2",
			"created_at": "2020-02-28 11:55:51",
			"updated_at": "2020-02-28 21:15:55",
			"bind_time": "2020-02-28 11:55:51",
			"bind_id": "008796764013583",
			"bind_type": "device",
			"ticket": 6,
			"diamond": 15,
			"register_time": 1582862151,
			"is_forbidden": 0,
			"login_status": 1,
			"is_new_user": 0,
			"is_new_guide": 0,
			"is_sync_diamond": 1,
			"archive_type": 1,
			"is_change_device": 0,
			"common_data_updated_at": 1597631804,
			"script_url": "http://cdnoss.stardustgod.com/chapters_test/file/202112/61c13ade560b6.jpg",
			"new_user_group_id": "new_user_group_b",
			"delete_time": 0,
			"is_deleted": 0,
			"is_delete_produce_new_user": 0
		},
		"switch": {
			"on_off_configs": [
				{
					"name": "is_cdkey",
					"value": "0",
					"code_version": "1"
				},
				{
					"channel": "AVG10003",
					"name": "language_data",
					"value": "[]"
				}
			],
			"version": [],
			"ab": {
				"ab_iapversion_2": 1,
				"ab_syncloading_1": 1,
				"ab_chapend_reward_g1_2": 1,
				"ab_chapend_reward_g2_1": 1,
				"ab_newbie_hall_1": 1,
				"ab_7day_daypass_1": 1
			},
			"member": {
				"member_type": 1
			},
			"summary_version": 16,
			"error_code": 2,
			"error_msg": "success",
			"error_ver": 1582895800
		}
	}`

	var inputbytes bytes.Buffer
	w := zlib.NewWriter(&inputbytes)
	w.Write([]byte(inputstr))
	w.Close()
	result := base64.StdEncoding.EncodeToString(inputbytes.Bytes())
	fmt.Println(result)

}

func DoReadFile(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s is not exist \n", path)
		}
		return ""
	} else {
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		return string(b)
	}
}

func DoReadFileByte(path string) []byte {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s is not exist \n", path)
		}
		return nil
	} else {
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		return b
	}
}
