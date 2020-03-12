package seed

import (
	"fmt"
)

type Seeder interface {
	Platform() string
	Link(int, int) string
	Home() string
	Name() string
	Identifier() string
}

type Zhihu struct {
	ID       string
	Zhuanlan string
}

func (z *Zhihu) Link(offset, limit int) string {
	if limit == 0 {
		limit = 10
	}
	return fmt.Sprintf("https://zhuanlan.zhihu.com/api/columns/%s/articles?offset=%d&limit=%d", z.ID, offset, limit)
}

func (z *Zhihu) Home() string {
	return fmt.Sprintf("https://zhuanlan.zhihu.com/%s", z.ID)
}

func (z *Zhihu) Name() string {
	return z.Zhuanlan
}

func (z *Zhihu) Platform() string {
	return "zhihu"
}

func (z *Zhihu) Identifier() string {
	return z.ID
}

func NewSeeder(platform, id string, name string) Seeder {
	switch platform {
	case "zhihu":
		return &Zhihu{ID: id, Zhuanlan: name}
	default:
		return nil
	}
}
