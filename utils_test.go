package go_workflow

import "testing"

func TestMd5Sum(t *testing.T) {
	cases := []struct {
		str      string
		expected string
	}{
		{
			str:      "1aaa",
			expected: "068896804551b982511aff19931628b9",
		}, {
			str:      "123456",
			expected: "e10adc3949ba59abbe56e057f20f883e",
		},
	}
	for _, c := range cases {
		value := md5Sum(c.str)
		if value != c.expected {
			t.Errorf("md5Sum error,actual:%s,expected:%s", value, c.expected)
		}
	}
}

func TestListRemove(t *testing.T) {
	cases := []struct {
		list    []string
		element string
	}{
		{list: []string{"1", "2", "3"}, element: "3"},
		{list: []string{"1", "2", "3"}, element: "2"},
		{list: []string{"1", "2", "3"}, element: "1"},
	}

	for _, c := range cases {
		t.Logf("%+v--%s-->%+v", c.list, c.element, listRemove(c.list, c.element))
	}
}
