package utils

import (
	"github.com/upyun/go-sdk/v3/upyun"
)

type FormUploadConfig struct {
	LocalPath      string                   // 待上传的文件路径
	SaveKey        string                   // 保存路径
	ExpireAfterSec int64                    // 签名超时时间
	NotifyUrl      string                   // 结果回调地址
	Apps           []map[string]interface{} // 异步处理任务
	Options        map[string]interface{}   // 更多自定义参数
}

func Upload(filePath string, filename string) error {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "image-chatchat",
		Operator: "ziyu",
		Password: "BmpDdXw4QvpWtwB96UuHanTuIGoD63Yu",
	})
	FR, err := up.FormUpload(&upyun.FormUploadConfig{
		LocalPath: filePath,
		SaveKey:   "/chatchatUsers/" + filename,
		NotifyUrl: "",
	})
	if err != nil {
		return err
	}
	println(FR)
	return nil
	// 上传文件

}
