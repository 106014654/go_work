package delete_slice

import "fmt"

type Number interface {
	int | uint | string | int64
}

func DeleteSliceT[T Number](a []T, index int) []T {
	i := 0
	for k, v := range a {
		if k != index {
			a[i] = v
			i++
		}
	}
	return a[:i]
}

func main() {
	var a = []int{1, 2, 3, 4, 5, 6}
	b := DeleteSliceT[int](a, 3)

	var c = []string{"a", "b", "c", "d", "f"}
	d := DeleteSliceT[string](c, 3)

	fmt.Println(b, d)

}