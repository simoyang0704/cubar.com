package conf

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CBCategory struct {
	Cate       Category   `json:"cate"`
	ChildCates []Category `json:"child_cates"`
}

var (
	CBCates = []CBCategory{
		CBCategory{
			Cate: Category{
				Id:   1000,
				Name: "运动",
			},
			ChildCates: []Category{
				Category{
					Id:   100010001,
					Name: "篮球",
				},
				Category{
					Id:   100010002,
					Name: "足球",
				},
				Category{
					Id:   100010003,
					Name: "跑步",
				},
				Category{
					Id:   100010004,
					Name: "羽毛球",
				},
				Category{
					Id:   100010005,
					Name: "乒乓球",
				},
				Category{
					Id:   100010006,
					Name: "游泳",
				},
			},
		},
		CBCategory{
			Cate: Category{
				Id:   1001,
				Name: "娱乐",
			},
			ChildCates: []Category{
				Category{
					Id:   100110001,
					Name: "麻将",
				},
				Category{
					Id:   100110002,
					Name: "扑克",
				},
				Category{
					Id:   100110003,
					Name: "k歌",
				},
			},
		},
		CBCategory{
			Cate: Category{
				Id:   1002,
				Name: "旅游",
			},
			ChildCates: []Category{
				Category{
					Id:   100210001,
					Name: "自驾游",
				},
				Category{
					Id:   100210002,
					Name: "骑行",
				},
				Category{
					Id:   100210003,
					Name: "自由行",
				},
			},
		},
		CBCategory{
			Cate: Category{
				Id:   1003,
				Name: "技艺",
			},
			ChildCates: []Category{
				Category{
					Id:   100310001,
					Name: "厨艺",
				},
				Category{
					Id:   100310002,
					Name: "育儿",
				},
				Category{
					Id:   100310003,
					Name: "分享",
				},
			},
		},
		CBCategory{
			Cate: Category{
				Id:   CATE_OTHER,
				Name: "其他",
			},
		},
	}

	CATE_OTHER = 1999
)
