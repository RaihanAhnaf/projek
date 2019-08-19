package helpers

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/eaciit/knot/knot.v1"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)
var ContractManager = []string{
	"Thomas,Daniel",
	"Wilson,Thomas",
	"Jepp,Thomas",
	"Thenkarai Sankaran,Nagasubramanian",
	"S,Chidambaram Balaji",
	"Raju,Sriram",
	"Ramakrishnan,Sivaramakrishnan",
	"C V,Ramanatha Siva",
	"Kealeboga,Shadrack Ramontsho",
	"Venkatasubramanian,Vasudevan",
	"Kadam,Vaibhav Rajaram",
	"D,Ramkumar",
	"Muthusubramanian,Muthumaheswaran",
	"Charamba,Stewart Itai",
	"Lim,Chin Ee James",
	"Lim,Lai Yee",
	"Ong,Liok Lim",
	"Lim,Chee Choong",
}

func GetContractManager() string {
	return ContractManager[rand.Intn(len(ContractManager)-1)]
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf/app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}

func GetCluster(clusterstr string) (string, float64) {
	str := strings.Split(clusterstr, "\t")
	cluster := str[0]
	//confidence := str[1]
	//confstr := strings.Split(confidence, "\r\n")
	confre, err := regexp.Compile("0\\.[0-9]+")
	if err != nil {
		tk.Println(err.Error())
	}
	confStr := string(confre.Find([]byte(clusterstr)))
	conffloat, err := strconv.ParseFloat(confStr, 32)
	if err != nil {
		tk.Println(err.Error())
	}
	confpercent := conffloat * 100
	return cluster, confpercent
	/*
		if len(confstr) == 0 {
			//strs := strings.Split(confidence, " ")
			confre, _ := regexp.Compile("0\\.[0-9]+")
			confStr := string(confre.Find([]byte(clusterstr)))
			conffloat, err := strconv.ParseFloat(confStr, 32)
			if err != nil {
				tk.Println(err.Error())
			}
			confpercent := conffloat * 100
			return cluster, confpercent
		} else {
			strs := strings.Split(confstr[0], " ")
			conffloat, err := strconv.ParseFloat(strs[1], 32)
			if err != nil {
				tk.Println(err.Error())
			}
			confpercent := conffloat * 100
			return cluster, confpercent
		}*/

}

var (
	DebugMode bool
)

func CreateResult(success bool, data interface{}, message string) map[string]interface{} {
	if !success {
		tk.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func UploadHandlerCopy(r *knot.WebContext, filename, dstpath string) (error, string) {
	r.Request.ParseMultipartForm(32 << 20)
	file, handler, err := r.Request.FormFile(filename)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	dstSource := dstpath + tk.PathSeparator + handler.Filename
	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	io.Copy(f, file)

	return nil, handler.Filename
}
func UploadHandler(r *knot.WebContext, filename, dstpath string) (error, string) {
	file, handler, err := r.Request.FormFile(filename)

	if err != nil {
		return err, ""
	}
	defer file.Close()

	dstSource := dstpath + tk.PathSeparator + handler.Filename
	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err, ""
	}
	defer f.Close()
	io.Copy(f, file)

	return nil, handler.Filename
}
