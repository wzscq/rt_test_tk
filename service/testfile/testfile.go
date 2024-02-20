package testfile

import (
	"os"
	"log"
	"bufio"
	"strconv"
	"io"
	"encoding/json"
	"strings"
	"time"
)

type TestFile struct {
	OutPath string
	DeviceID string
	TimeStamp string
	ContentFile *os.File
	StartTime int64
	LineCount int64
	//下面两个字段用于控制文件超时关闭
	lastLineCount int64
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	RGB string `json:"rgb"`
	Value float64  `json:"value"`
}

func (tf *TestFile) Close() {
	indexFileName:=tf.OutPath+"/"+tf.DeviceID+"_"+tf.TimeStamp+".idx"
	idxFile,err:=os.OpenFile(indexFileName,os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
	} else {
		idxFile.WriteString(strconv.FormatInt(tf.LineCount,10)+"\n")
		idxFile.Close()
	}
	tf.ContentFile.Close()
}

func (tf *TestFile) CloseReadOnly() {
	tf.ContentFile.Close()
}

func (tf *TestFile) WriteLine(lineContent string) {
	tf.ContentFile.WriteString(lineContent+"\n")
	tf.LineCount++
}

func (tf *TestFile)GetIdxContent()(string){
	indexFileName:=tf.OutPath+"/"+tf.DeviceID+"_"+tf.TimeStamp+".idx"
	idxFile,err:=os.Open(indexFileName)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	defer idxFile.Close()
		
	bytes:=make([]byte,1024)
	n,err:=idxFile.Read(bytes)
	if err != nil {
		log.Printf("Read file failed [Err:%s]\n", err.Error())
		return ""
	}

	if n==0 {
		return ""
	}

	return string(bytes[:n])
}

func (tf *TestFile)GetContent(from,to int64)([]string){
	lines:=[]string{}
	//文件复位
	_, err := tf.ContentFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
		return lines
	}

	scanner := bufio.NewScanner(tf.ContentFile)
	var n int64
	n = -1
	for scanner.Scan() {
		n++
		lineStr:=string(scanner.Bytes())
		log.Printf("from:%d,to:%d,n:%d,line:%s",from,to,n,lineStr)
        if n < from {
            continue
        }
		lines=append(lines,lineStr)
		if n >= to {
            break
        }
		log.Println(lines)
    }
	log.Println(lines)
	return lines
}

func (tf *TestFile)getDataItme(line map[string]interface{})(map[string]interface{}){
	data,ok:=line["data"].([]interface{})
	if !ok {
		return nil
	}

	if len(data)==0 {
		return nil
	}

	dataItem:=data[0].(map[string]interface{})
	return dataItem
}

func (tf *TestFile)getIndicator(line map[string]interface{},extractPath string)(*float64){
	//首先获取data数据的第一个数据，目前仅考虑第一个
	dataItem:=tf.getDataItme(line)
	if dataItem==nil {
		return nil
	}

	pathNodes:=strings.Split(extractPath, ".")
	lastIndex:=len(pathNodes)-1
	for idx,pathItem:=range pathNodes {
		if idx<lastIndex {
			dataTmp,ok:=dataItem[pathItem]
			if !ok {
				return nil
			}
			dataItem=dataTmp.(map[string]interface{})
		} else {
			break
		}
	}

	valueData,ok:=dataItem[pathNodes[lastIndex]].(string)
	if !ok {
		return nil
	}

	val,err:=strconv.ParseFloat(valueData,64)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &val
}

func (tf *TestFile)getLegendItem(value float64,legendItems []IndicatorLegendItem)(*IndicatorLegendItem){
	for _,legendItem:=range legendItems {
		if legendItem.End == "" {
			return &legendItem
		}

		endVal,_:=strconv.ParseFloat(legendItem.End,64)
		if endVal>=value {
			return &legendItem
		}
	}
	return nil
}

func (tf *TestFile)getPoint(line map[string]interface{},indicator Indicator)(*Point){
	robotInfo,ok:=line["robot_info"].(map[string]interface{})
	if(!ok){
		return nil
	}

	x,_:=robotInfo["pixel_x"].(float64)
	y,_:=robotInfo["pixel_y"].(float64)
	
	value:=tf.getIndicator(line,indicator.ExtractPath)
	if value==nil {
		return nil
	}

	color:="#000000"
	legendItem:=tf.getLegendItem(*value,indicator.Legend.List)
	if legendItem!=nil {
		color=legendItem.RGB
	}
	return &Point{
		X:x,
		Y:y,
		RGB:color,
		Value:*value,
	}
}

func (tf *TestFile)GetPoints(indicator Indicator)([]*Point){
	//文件复位
	_, err := tf.ContentFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	points:=[]*Point{}
	scanner := bufio.NewScanner(tf.ContentFile)
	for scanner.Scan() {
		var line map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &line)
		if err != nil {
			log.Println(err)
			return nil
		}
		point:=tf.getPoint(line,indicator)
		//if point != nil {
			points=append(points,point)
		//}
  }
	return points
}

func GetTestFile(outPath string,deviceID string,timeStamp string) *TestFile {
	contentFileName:=outPath+"/"+deviceID+"_"+timeStamp+".content"
	contentFile,err:=os.OpenFile(contentFileName,os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Open file failed [Err:%s]\n", err.Error())
		return nil
	}

	return &TestFile{
		DeviceID:deviceID,
		TimeStamp:timeStamp,
		ContentFile:contentFile,
		OutPath:outPath,
		StartTime:time.Now().Unix(),		
	}
}