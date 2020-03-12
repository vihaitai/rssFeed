package seed

import "fmt"

type Zhihu struct {
	id       string
	zhuanlan string
}

func (z *Zhihu) Link(offset, limit int) string {
	if limit == 0 {
		limit = 10
	}
	return fmt.Sprintf("https://zhuanlan.zhihu.com/api/columns/%s/articles?offset=%d&limit=%d", z.id, offset, limit)
}

func (z *Zhihu) Home() string {
	return fmt.Sprintf("https://zhuanlan.zhihu.com/%s", z.id)
}

func (z *Zhihu) Name() string {
	return z.zhuanlan
}

func (z *Zhihu) Platform() string {
	return "zhihu"
}

func (z *Zhihu) Identifier() string {
	return z.id
}
