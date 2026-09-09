package version

import "testing"

func TestParse(t *testing.T) {
	v, err := Parse("26.3.1")
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if v.Year != 26 || v.Major != 3 || v.Minor != 1 {
		t.Errorf("解析结果不对: %+v", v)
	}
	if v.String() != "26.3.1" {
		t.Errorf("还原不对: %q", v.String())
	}
	if v.FullYear() != 2026 {
		t.Errorf("年份展开不对: %d", v.FullYear())
	}
	if _, err := Parse("v26.3.1"); err != nil {
		t.Errorf("带 v 前缀应能解析: %v", err)
	}
	for _, bad := range []string{"", "26.3", "26.3.1.0", "2026.3.1", "26.x.1", "26.-1.0"} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("%q 应报错", bad)
		}
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"26.1.0", "26.1.0", 0},
		{"26.1.0", "26.1.1", -1},
		{"26.2.0", "26.10.0", -1}, // 数字比较，10 > 2
		{"26.10.0", "26.9.9", 1},  // 字符串比较会误判，数字比较正确
		{"25.9.9", "26.1.0", -1},  // 跨年
		{"27.1.0", "26.9.9", 1},
	}
	for _, c := range cases {
		if got := CompareStrings(c.a, c.b); got != c.want {
			t.Errorf("CompareStrings(%q,%q)=%d, want %d", c.a, c.b, got, c.want)
		}
	}
}
