package CommonTools

import (
	"errors"
	"fmt"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
	~float32 | ~float64
}

// 定义支持的运算类型
type Operation string

const (
	Add      Operation = "+"
	Subtract Operation = "-"
	Multiply Operation = "*"
	Divide   Operation = "/"
	Avg      Operation = "avg"
)

// 计算函数
func Calculate[T any, V Number](items []T, getter func(T) V, op Operation) (V, error) {
	if len(items) == 0 {
		return 0, errors.New("empty slice")
	}

	// 初始值设置为第一个元素的值
	result := getter(items[0])
	var cnt V

	// 从第二个元素开始进行运算
	for i := 1; i < len(items); i++ {
		value := getter(items[i])
		switch op {
		case Add, Avg:
			result += value
			cnt += 1
		case Subtract:
			result -= value
		case Multiply:
			result *= value
		case Divide:
			if value == 0 {
				return 0, errors.New("division by zero")
			}
			result /= value
		default:
			return 0, fmt.Errorf("unsupported operation: %s", op)
		}
	}

	switch op {
	case Avg:
		result /= cnt
	}

	return result, nil
}
