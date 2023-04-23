
package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// 表过滤器的表过滤规则
// 过滤匹配成功是接受(positive)
// 过滤匹配不成功是拒绝(negative)
type tableRule struct {
	table matcher
}

// matcher 表规则过滤接口
type matcher interface {
	matchString(name string) bool
}

// stringMatcher 字符串 mather 接口实现
type stringMatcher string

func (m stringMatcher) matchString(name string) bool {
	// 忽略大小写
	return strings.ToUpper(string(m)) == strings.ToUpper(name)
}

// trueMatcher 匹配所有匹配器 The `*` pattern.
type trueMatcher struct{}

func (trueMatcher) matchString(string) bool {
	return true
}

// regexpMatcher 是基于正则表达式的匹配器
type regexpMatcher struct {
	pattern *regexp.Regexp
}

func newRegexpMatcher(pat string) (matcher, error) {
	if pat == "(?i)(^|([\\s\\t\\n]+))(.*$)" {
		// special case for '*'
		return trueMatcher{}, nil
	}

	pattern, err := regexp.Compile(pat)
	if err != nil {
		return nil, fmt.Errorf("newRegexpMatcher regexp compile failed: %v", err)
	}

	return regexpMatcher{pattern: pattern}, nil
}

func (m regexpMatcher) matchString(name string) bool {
	return m.pattern.MatchString(name)
}
