package os2
import (
	"io/ioutil"
	"os"
	"encoding/json"
)

// LoadFileToStruct 加载本地文件数据到结构体里
// dest必须是结构体的指针
func LoadFileToStruct(filename string, dest interface{}) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}
