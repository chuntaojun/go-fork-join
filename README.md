## Fork-Join

![License](https://img.shields.io/github/license/chuntaojun/go-fork-join.svg)
![stars](https://img.shields.io/github/stars/chuntaojun/go-fork-join.svg)
![forks](https://img.shields.io/github/forks/chuntaojun/go-fork-join.svg)


#### 使用示例

```go
package example

import (
	"fmt"
	"fork-join"
	"github.com/smartystreets/assertions/assert"
	"github.com/smartystreets/assertions/should"
	"testing"
	"time"
)


var taskPool = fork_join.NewForkJoinPool("pool", 10)

type SumAdd struct {
	start int64
	end   int64
	fork_join.ForkJoinTask
}

func (s *SumAdd) Compute() interface{} {

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("here is err %#v\n", p)
		}
	}()

	var sum int64
	if s.end-s.start < 1000 {
		tmp := int64(0)
		for i := s.start; i <= s.end; i ++ {
			tmp += i
		}
		sum = tmp
	} else {
		mid := (s.start + s.end) / 2
		sTask1 := &SumAdd{start: s.start, end: mid}
		sTask2 := &SumAdd{start: mid + 1, end: s.end}
		sTask1.Build(taskPool).Run(sTask1)
		sTask2.Build(taskPool).Run(sTask2)
		ok1, r1 := sTask1.Join()
		ok2, r2 := sTask2.Join()
		if ok1 && ok2 {
			sum = r1.(int64) + r2.(int64)
		}
	}
	return sum
}

func TestForkJoin(t *testing.T) {

	t1 := time.Now()
	v1 := int64(0)
	for i := int64(1); i <= 100000000; i ++ {
		v1 += i
	}
	elapsed := time.Since(t1)
	fmt.Println("Costumer App elapsed: ", elapsed)

	s := &SumAdd{start: 1, end: 100000000}
	t2 := time.Now()
	v2 := s.Compute()
	elapsed2 := time.Since(t2)
	fmt.Println("ForkJoin App elapsed: ", elapsed2)

	result := assert.So(v2, should.Equal, v1)
	fmt.Println(result.Log())

}
```