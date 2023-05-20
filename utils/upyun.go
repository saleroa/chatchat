package utils

import (
	"github.com/upyun/go-sdk/v3/upyun"
	"io"
)

type FormUploadConfig struct {
	LocalPath      string                   // 待上传的文件路径
	SaveKey        string                   // 保存路径
	ExpireAfterSec int64                    // 签名超时时间
	NotifyUrl      string                   // 结果回调地址
	Apps           []map[string]interface{} // 异步处理任务
	Options        map[string]interface{}   // 更多自定义参数
}

func Upload(w io.Reader, filename string) error {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "image-chatchat",
		Operator: "ziyu",
		Password: "BmpDdXw4QvpWtwB96UuHanTuIGoD63Yu",
	})
	err := up.Put(&upyun.PutObjectConfig{
		Reader: w,
		Path:   "/chatchatUsers/" + filename,
	})
	if err != nil {
		return err
	}
	return nil
	// 上传文件

}

func Delete(filename string) error {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "image-chatchat",
		Operator: "ziyu",
		Password: "BmpDdXw4QvpWtwB96UuHanTuIGoD63Yu",
	})
	err := up.Delete(&upyun.DeleteObjectConfig{
		Async: true,
		Path:  "/chatchatUsers/" + filename + "/",
	})
	if err != nil {
		return err
	}
	return nil
	// 上传文件

}
