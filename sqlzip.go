// sqlzip
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	zipk "github.com/alexmullins/zip" //带加密口令的zip
	//"crypto/aes"
	//"crypto/cipher"
	//"encoding/base64"
)

var ZPssWord string = "XYUYXA1234***{}_+)_)11111111111111111111111111111110"
var DAction string
var FPath1 string
var FPath2 string
var DPssWord string
var autoRun bool = false

func main() {
	//按照命令参数连接
	comargstr := os.Args

	//设置运行目录
	Chdir()

	if len(comargstr) >= 2 {
		DAction = comargstr[1] //获取动作
		fmt.Println("Commandline:", comargstr)
		autoRun = true //要求自动关闭
		if len(comargstr) >= 3 {
			FPath1 = comargstr[2] //获取1文件名
		}
		if len(comargstr) >= 4 {
			FPath2 = comargstr[3] //获取2文件名
		}
		if len(comargstr) >= 5 {
			DPssWord = comargstr[4] //获取密码
		}
	} else {
		//无参数模式
		fmt.Println("<<sqlzip>> 压缩命令启动： sqlzip.exe -z filepathname zipname mypassword")
		fmt.Println("<<sqlzip>> 解压命令启动： sqlzip.exe -u zipfile unzipfile mypassword")
		fmt.Println("<<sqlzip>> 缺省_替代示例： sqlzip.exe -u zipfile _ mypassword")
		fmt.Println(">压缩文件夹：sql -all")
	}

	//压缩文件  --可以使用未过期的文件避免重复压缩
	sqldir := "./sql"
	file := "./source.sydp"
	sqlexp := "./sqlexp"
	//指令逻辑处理
	if DPssWord != "" && DPssWord != "_" {
		ZPssWord = DPssWord
	}

	if DAction == "-z" {
		//压缩
		if FPath1 != "" && FPath1 != "_" {
			sqldir = FPath1
		}
		if FPath2 != "" && FPath2 != "_" {
			file = FPath2
		}
		//goto GTOZIP
		fmt.Println(">压缩文件夹：", sqldir, file)
		ZipAction(sqldir, file)
		return
	} else if DAction == "-u" {
		//解压
		if FPath1 != "" && FPath1 != "_" {
			file = FPath1
		}
		if FPath2 != "" && FPath2 != "_" {
			sqlexp = FPath2
		}
		//goto GTOUNZIP
		fmt.Println(">解压文件：", file, sqlexp)
		UnZipAction(file, sqlexp)
		return
	}

	//文件夹下是否有文件
	var filelist []string
	filelist, _ = GetAllFile(sqldir, filelist)
	if len(filelist) == 0 {
		fmt.Println("  警告 " + sqldir + " 是空文件夹, 即将退出！")
		WaitEnter("准备退出")
		return
	}
	//加密处理

	search := GetFileName(file)

	if search != "" {
		fmt.Println("  存在 " + file + " 文件")
		fmt.Println(">>>>重新压缩文件：", file)
	} else {
		fmt.Println(">>>>创建压缩文件：", file)
	}

	//压缩
	ZipAction(sqldir, file)
	// itms := time.Now()
	// //err = Zip(sqldir, file)
	// err = ZipEncrypt(sqldir, file, ZPssWord)

	// if err != nil {
	// 	fmt.Println("X 压缩失败：", sqldir, ">->", file, err.Error())
	// }
	// itme := time.Now()
	// ms1 := (itme.UnixNano() - itms.UnixNano()) / 1e6
	// strresult := fmt.Sprintf("压缩耗时：%v ms", NumberFormat(strconv.FormatInt(ms1, 10)))
	// fmt.Println(strresult)
	if !autoRun {
		time.Sleep(200000000) //等待0.2秒,显示内容
		WaitEnter("准备解压")
	} else {
		//自动运行直接结束
		return
	}

	//解压缩
	fmt.Println(">>>>解压文件：", file, " 释放到 ", sqlexp)
	UnZipAction(file, sqlexp)
	// itms2 := time.Now()
	// //err = UnZip(file, sqlexp)
	// err = UnZipDecrypt(file, sqlexp, ZPssWord)
	// if err != nil {
	// 	fmt.Println("X 解压失败：", file, ">->", sqlexp, err.Error())
	// }
	// itme2 := time.Now()
	// ms2 := (itme2.UnixNano() - itms2.UnixNano()) / 1e6
	// strresult2 := fmt.Sprintf("解压耗时：%v ms", NumberFormat(strconv.FormatInt(ms2, 10)))
	// fmt.Println(strresult2)
	if !autoRun {
		WaitEnter("程序执行完成")
	}
}

func ZipAction(sqldir string, file string) {
	var err error
	//压缩
	itms := time.Now()
	//err = Zip(sqldir, file)
	err = ZipEncrypt(sqldir, file, ZPssWord)

	if err != nil {
		fmt.Println("X 压缩失败：", sqldir, ">->", file, err.Error())
	}
	itme := time.Now()
	ms1 := (itme.UnixNano() - itms.UnixNano()) / 1e6
	strresult := fmt.Sprintf("压缩耗时：%v ms", NumberFormat(strconv.FormatInt(ms1, 10)))
	fmt.Println(strresult)

}
func UnZipAction(file string, sqlexp string) {
	//解压缩
	var err error
	fmt.Println(">>>>解压文件：", file, " 释放到 ", sqlexp)
	itms2 := time.Now()
	//err = UnZip(file, sqlexp)
	err = UnZipDecrypt(file, sqlexp, ZPssWord)
	if err != nil {
		fmt.Println("X 解压失败：", file, ">->", sqlexp, err.Error())
	}
	itme2 := time.Now()
	ms2 := (itme2.UnixNano() - itms2.UnixNano()) / 1e6
	strresult2 := fmt.Sprintf("解压耗时：%v ms", NumberFormat(strconv.FormatInt(ms2, 10)))
	fmt.Println(strresult2)
}

func WaitEnter(msg string) {
	fmt.Println(msg, ",输入回车继续：")
	var sinput string
	//fmt.Scanf("字符%s", &sinput)
	fmt.Scanln(&sinput)
}

func GetFileName(filename string) string {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println(filename, err)
		return ""
	} else {
		return fileInfo.Name()
	}

}

//取文件名
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

//格式化数值    1,234,567,898.55
func NumberFormat(str string) string {
	length := len(str)
	if length < 4 {
		return str
	}
	arr := strings.Split(str, ".") //用小数点符号分割字符串,为数组接收
	length1 := len(arr[0])
	if length1 < 4 {
		return str
	}
	count := (length1 - 1) / 3
	for i := 0; i < count; i++ {
		arr[0] = arr[0][:length1-(i+1)*3] + "," + arr[0][length1-(i+1)*3:]
	}
	return strings.Join(arr, ".") //将一系列字符串连接为一个字符串，之间用sep来分隔。
}

// Chdir 将程序工作路径修改成程序所在位置
func Chdir() (err error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	err = os.Chdir(dir)
	return
}

//zip压缩 filePath 为需要压缩的文件路径，zipPath为压缩后文件路径
func FileToZip(filePath string, zipPath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	z, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer z.Close()

	wr := zip.NewWriter(z)
	// 因为filePath是一个路径，所以会创建路径中的所有文件夹
	w, err := wr.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	return nil
}

//压缩文件夹
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

//解压文件
func UnZip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

//先压缩+再加密文件  先压缩,一旦你加密文件，将生成一个随机数据流，这是不可压缩的。压缩过程依赖于在数据中找到可压缩模式。
func ZipEncrypt(srcFile string, destZip string, passwd string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zipk.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zipk.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zipk.Deflate
		}

		// 不加密的文件头
		// writer, err := archive.CreateHeader(header)
		// if err != nil {
		// 	return err
		// }
		//加密单个文件头 FileHeader
		writer, err := archive.Encrypt(header.Name, passwd)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

//解压+解密文件
func UnZipDecrypt(archive, target string, passwd string) error {
	reader, err := zipk.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		file.SetPassword(passwd)
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
