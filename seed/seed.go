package seed

type Seeder interface {
	Platform() string
	Link(int, int) string
	Home() string
	Name() string
	Identifier() string
}

func NewSeeder(platform, id string, name string) Seeder {
	switch platform {
	case "zhihu":
		return &Zhihu{ID: id, Zhuanlan: name}
	default:
		return nil
	}
}
