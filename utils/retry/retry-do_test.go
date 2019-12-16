package retry

import (
	"fmt"
	"github.com/pkg/errors"
	"shangwoa.com/utils/number"
	"testing"
)
/// go test -v retry-do_test.go retry-do.go
func TestRetryDoInteralFunc(t *testing.T) {
	exec := func() (error, interface{}){
		r := number.Random(0, 100000000)
		fmt.Println("r is", r)
		if r < 100000000 / 2{
			return errors.New("xx"), nil
		}
		return nil, map[string]string{"abc":"wocao"}
	}
	err, res, count := RetryDoInteralTime(exec, 5, 100)
	fmt.Printf("TestRetryDoInteralFunc error is %s res is %v count is %d", err, res, count)
}


