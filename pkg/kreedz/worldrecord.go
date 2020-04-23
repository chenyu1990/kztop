package kreedz

import (
	"bufio"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io"
	"io/ioutil"
	"kztop/pkg/util"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var SortName = []string{
	"TEST",
	"XJ",
	"CC",
	"WS",
}

var Name = []string{
	"TEST",
	"Xtreme-Jumps",
	"Cosy-Climbing",
	"World-Surf",
}

var organizations = []string{
	"http://kztop:8080/debug.txt",
	"https://xtreme-jumps.eu/demos.txt",
	"https://cosy-climbing.net/demoz.txt",
	"http://world-surf.com/demos.txt",
}

var localFile = []string{
	"./data/kreedz/debug.txt",
	"./data/kreedz/xj.txt",
	"./data/kreedz/cc.txt",
	"./data/kreedz/ws.txt",
}

var cacheTime []int64

func init() {
	cacheTime = make([]int64, len(localFile))

	err := os.MkdirAll("./data/kreedz", 0755)
	if err != nil {
		panic(err)
	}
}

func (a *WorldRecord) FirstSync() (bool, []*RecordInfo, error) {
	file, err := os.Open(localFile[a.Organization])
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil && os.IsNotExist(err) {
		err = a.downloadFile(localFile[a.Organization])
		if err != nil {
			return false, nil, err
		}
		newFile, err := os.Open(localFile[a.Organization])
		if err != nil {
			panic(err)
		}
		defer newFile.Close()

		var records []*RecordInfo
		br := bufio.NewReader(newFile)
		//_, _, _ = br.ReadLine()
		for {
			record, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			recordInfo := a.unserialize(string(record))
			records = append(records, recordInfo)
		}

		return true, records, nil
	}

	return false, nil, nil
}

func (a *WorldRecord) fileTimeout() bool {
	// 本地文件不存在，直接返回“过期”
	stat, err := os.Stat(localFile[a.Organization])
	if err != nil {
		panic(err)
	}

	// 没有缓存时间，初次启动，读取本地时间
	if cacheTime[a.Organization] == 0 {
		cacheTime[a.Organization] = stat.ModTime().Unix()
	}

	a.getFileInfo()

	location, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	// Thu, 02 Apr 2020 13:00:44 GMT
	fileTime, err := time.ParseInLocation(time.RFC1123,
		a.recordFileHeader.Get("Last-Modified"), location)
	if err != nil {
		panic(err)
	}
	a.NewsDate = fileTime

	if fileTime.Unix() > cacheTime[a.Organization] {
		return true
	}

	return false
}

func (a *WorldRecord) getDiff() error {
	old, err := ioutil.ReadFile(localFile[a.Organization])
	if err != nil {
		panic(err)
	}
	new, _ := ioutil.ReadFile(localFile[a.Organization] + "new")
	if err != nil {
		panic(err)
	}

	// map[mapname][holderAndExHolder]*recordInfo
	news := make(map[string]map[string]*RecordInfo)
	diffs := Diff(string(old), string(new))
	for _, d := range diffs {
		var recordInfos []*RecordInfo
		if strings.Contains(d.Text, "\n") {
			split := strings.Split(d.Text, "\n")
			for _, record := range split {
				recordInfo := a.unserialize(record)
				if recordInfo == nil {
					continue
				}
				recordInfos = append(recordInfos, recordInfo)
			}
		} else {
			recordInfo := a.unserialize(d.Text)
			if recordInfo == nil {
				continue
			}
			recordInfos = append(recordInfos, recordInfo)
		}
		for _, recordInfo := range recordInfos {
			if d.Type != diffmatchpatch.DiffEqual {
				if news[recordInfo.MapName] == nil {
					news[recordInfo.MapName] = make(map[string]*RecordInfo)
				}
			}
			switch d.Type {
			case diffmatchpatch.DiffInsert:
				a.NewRecords = append(a.NewRecords, recordInfo)
				news[recordInfo.MapName]["holder"] = recordInfo
			case diffmatchpatch.DiffDelete:
				news[recordInfo.MapName]["exHolder"] = recordInfo
			}
		}
	}

	for mapname, records := range news {
		if a.News == nil {
			a.News = make(map[string]map[string][]*RecordInfo)
		}
		holder := records["holder"].Holder
		if holder == "n/a" {
			continue
		}

		if a.News[holder] == nil {
			a.News[holder] = make(map[string][]*RecordInfo)
		}
		if a.News[holder][mapname] == nil {
			a.News[holder][mapname] = make([]*RecordInfo, 2)
		}
		records["holder"].MapName = ""
		a.News[holder][mapname][0] = records["holder"]
		if _, ok := records["exHolder"]; ok {
			records["exHolder"].MapName = ""
			a.News[holder][mapname][1] = records["exHolder"]
		}

		//for who, recordInfo := range records {
		//	switch who {
		//	case "holder":
		//		if recordInfo.Time > 0 {
		//			a.News[recordInfo.Holder][recordInfo.MapName][0] = recordInfo
		//		}
		//	case "exHolder":
		//		if records["holder"] != nil {
		//			a.News[recordInfo.Holder][recordInfo.MapName][1] = recordInfo
		//		}
		//	}
		//}

	}
	return nil
}

func (a *WorldRecord) getFileInfo() {
	resp, err := http.Head(organizations[a.Organization])
	if err != nil {
		util.HandleHttpError(err)
		return
	}

	defer resp.Body.Close()
	a.recordFileHeader = resp.Header
}

func (a *WorldRecord) downloadFile(savePath string) error {
	resp, err := http.Get(organizations[a.Organization])
	if err != nil {
		util.HandleHttpError(err)
		return err
	}

	f, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, resp.Body)
	return nil
}

func (a *WorldRecord) CheckUpdate(organization *Organization) (bool, error) {
	if organization == nil {
		panic("have to input Organization")
	}
	a.Organization = *organization

	if a.fileTimeout() == true {
		err := a.downloadFile(localFile[a.Organization] + "new")
		if err != nil {
			return false, nil
		}
		a.getDiff()
		if len(a.News) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (a *WorldRecord) CopyFile() {
	file, err := ioutil.ReadFile(localFile[a.Organization] + "new")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(localFile[a.Organization], file, os.ModePerm)
}

func (a *WorldRecord) unserialize(record string) *RecordInfo {
	split := strings.Split(record, " ")
	if len(split) < 4 {
		return nil
	}

	time, err := strconv.ParseFloat(split[1], 10)
	if err != nil {
		panic(err)
	}
	return &RecordInfo{
		MapName: split[0],
		Holder:  split[2],
		Region:  split[3],
		Time:    time,
	}
}
