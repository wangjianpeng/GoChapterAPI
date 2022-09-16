package api

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func DoConsume(ticketNum int, diamondNum int) {
	fmt.Printf("do it")
}

func httpDo() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "baidu.com", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

// application/x-www-form-urlencoded
//signature
//platform
//u3dapi
func Do1() {
	values := map[string]string{"uuid": "John Doe", "value_type": "diamond", "value_change": "1", "channel_id": "AVG50005"}
	// var sigvalue string = ""
	sigvalue := ""
	for _, v := range values {
		sigvalue = sigvalue + v
	}
	//add signature
	sigvalue = sigvalue + `MARS_SECRET_AVG_!@#$%`
	hash := md5.Sum([]byte(sigvalue))
	values["signature"] = hex.EncodeToString([]byte(hash[:]))
	json_data, err := json.Marshal(values)
	requestBody := base64.StdEncoding.EncodeToString(StringToAsciiBytes(string(json_data)))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://dev-chapters-int.stardustgod.com/syncValueApi.Class.php", "application/json",
		bytes.NewBuffer([]byte(requestBody)))

	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp.Header.Set("platform", "ios")
	resp.Header.Set("u3dapi", "true")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("MyHeader:	", resp.Header.Get("Content-Encoding"))
	// var res map[string]interface{}

	// json.NewDecoder(resp.Body).Decode(&res)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	bodystr := string(body)
	fmt.Println("body:	", bodystr)

	// base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(bodystr)))

	// n, _ := base64.StdEncoding.Decode(base64Text, []byte(bodystr))
	bodystr = `Nzc2NzExMzA0NDkwMjA0ZWY0MDI0Zjc5MmY1NjlmNzV4nKtWSi0qyi+KT85PSVWy0jU1MNCBiuQWpytZKQWnFpWlFimAhZRgUkARJStDM2NzQxMLY1PjWgAvhhZl`
	s2, _ := base64.StdEncoding.DecodeString(bodystr)

	//nr := bytes.NewReader(s2[2:])
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~!!!!!!!!!!!!!!")
	// //r, _ := gzip.NewReader(nr)
	// //result, _ := ioutil.ReadAll(r)
	// rt, _ := gUnzipData(StringToAsciiBytes(string(s2[2:])))
	// fmt.Println("base64Text:", string(rt))
	// for i, b := range rt[2:] {
	// 	fmt.Println(i, b)
	// }
	// fmt.Println(string(Inflate(s2)))

	// comp := compressor{s2[2:]}
	// buffer := comp.decompress()
	// fmt.Print("Uncompressed data: ")
	// fmt.Println(string(buffer))

	// fmt.Println(s2[:4])
	// r := flate.NewReader(bytes.NewReader(StringToAsciiBytes(string(s2))))
	// if err != nil {
	// 	panic(err)
	// }
	// enflated, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(enflated))
	for i, b := range StringToAsciiBytes(string(s2)) {
		fmt.Println(i, b)
	}

	content := []byte{55, 55, 54, 55, 49, 49, 51, 48, 52, 52, 57, 48, 50, 48, 52, 101, 102, 52, 48, 50, 52, 102, 55, 57, 50, 102, 53, 54, 57, 102, 55, 53, 120, 156, 171, 86, 74, 45, 42, 202, 47, 138, 79, 206, 79, 73, 85, 178, 210, 53, 53, 48, 208, 129, 138, 228, 22, 167, 43, 89, 41, 5, 167, 22, 149, 165, 22, 41, 128, 133, 148, 96, 82, 64, 17, 37, 43, 67, 51, 99, 115, 67, 19, 11, 99, 83, 227, 90, 0, 47, 134, 22, 101}

	enflated, err := ioutil.ReadAll(flate.NewReader(bytes.NewReader(content[2 : len(content)-4])))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enflated))

}

// if (!dic.ContainsKey("channel_id"))
//             {
//                 dic.Add("channel_id", ConfigSetting.Channel);
//             }
//             foreach (var item in dic)
//             {
//                 values += item.Value;
//             }
//             if (autoSign)
//             {
//                 string md5Hash = Extensions.CalculateMD5Hash(Encoding.UTF8.GetBytes(values + ConfigSetting.NetKey));
//                 if (dic.ContainsKey("signature"))
//                     dic["signature"] = md5Hash;
//                 else
//                     dic.Add("signature", md5Hash);
//             }
//             string sendJson = JsonMapper.ToJson(dic);
//             byte[] sendStr = CommonTools.ZipFromStrToBytes(sendJson);
//             request.RawData = Encoding.UTF8.GetBytes(Convert.ToBase64String(sendStr));
//             request.SetHeader("platform", "ios");
//             request.SetHeader("u3dapi", "true");

func StringToAsciiBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++

	}

	return t
}

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func gZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}

// Inflate utility that decompresses a string using the flate algo
func Inflate(deflated []byte) []byte {
	var b bytes.Buffer
	fmt.Println("bt length before:", len(deflated))
	brd := bytes.NewReader(deflated)
	fmt.Println("bt length before2:", brd.Len())
	r := flate.NewReader(brd)
	b.ReadFrom(r)
	r.Close()
	fmt.Println("bt length:", len(b.Bytes()))
	return b.Bytes()
}

type compressor struct {
	content []byte
}

func (r *compressor) decompress() []byte {
	dc := flate.NewReader(bytes.NewReader(r.content))
	defer dc.Close()
	rb, err := ioutil.ReadAll(dc)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Fatalf("Err %v\n read %v", err, rb)
		}
	}
	return rb
}

func compressFlate(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, 9)
	if err != nil {
		return nil, err
	}
	w.Write(data)
	w.Close()
	return b.Bytes(), nil
}

func decompressString(str string) (string, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer([]byte(str)))
	defer gr.Close()
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DoPostWithHeader() {
	rawstr := "eJwdzD0OgCAMQOG7dHZoKIXi5uQNXA3hR0kMg8pkvLvE9eXleyDsvtZ0rCXCCNMyMyIyDHCVrfq7naln44iz5oRokFBQOycxZKWIyVqRvrf2A0xaEbwfTegXgg=="
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", "http://dev-chapters-int.stardustgod.com/Controllers/mall/GetMsgProductShelfListApi.php", bytes.NewBuffer([]byte(rawstr)))
	// request.Header.Set("Content-type", "application/json")
	request.Header.Set("platform", "ios")
	request.Header.Set("u3dapi", "true")
	// request.Header.Set("language", "en-US")
	// request.Header.Set("mcf", "3CN0")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
