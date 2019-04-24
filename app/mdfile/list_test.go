package mdfile

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	// sw := sync.WaitGroup{}

	// for i := 0; i < 10; i++ {
	// 	sw.Add(1)
	// 	go func() {
	// 		l := Model
	// 		fmt.Println(l)
	// 		fmt.Println("-------------")
	// 		fmt.Println("")
	// 		sw.Done()
	// 	}()
	// }

	// sw.Wait()

	articles := Model.ArticlesAll()

	for _, article := range articles {
		fmt.Println(article.CreatedAt)
		fmt.Println(article.UpdatedAt)
		fmt.Println("-----------------------------")
	}

}
