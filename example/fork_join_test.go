package example

import (
	"fmt"
	"fork-join"
	"testing"
	"time"
)

var taskPool = fork_join.NewForkJoinPool("pool", 10)
var Id = int32(0)

type SumAdd struct {
	start int64
	end   int64
	fork_join.ForkJoinTask
}

func (s *SumAdd) Compute() interface{} {
	var sum int64
	if s.end-s.start < 2 {
		tmp := int64(0)
		for i := s.start; i <= s.end; i ++ {
			tmp += i
		}
		sum = tmp
	} else {
		mid := (s.start + s.end) / 2
		sTask1 := &SumAdd{start: s.start, end: mid}
		sTask2 := &SumAdd{start: mid + 1, end: s.end}
		sTask1.Build(taskPool, &Id).Run(sTask1)
		sTask2.Build(taskPool, &Id).Run(sTask2)
		sum = sTask1.Join().(int64) + sTask2.Join().(int64)
	}
	return sum
}

func TestForkJoin(t *testing.T) {
	t1 := time.Now()
	v1 := int64(0)
	for i := int64(1); i <= 1000000; i ++ {
		v1 += i
	}
	fmt.Printf("result v1 is %#v\n", v1)
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)

	s := &SumAdd{start: 1, end: 100}
	t2 := time.Now()
	v2 := s.Compute()
	fmt.Printf("result v2 is %#v\n", v2)
	elapsed2 := time.Since(t2)
	fmt.Println("App elapsed: ", elapsed2)

}
