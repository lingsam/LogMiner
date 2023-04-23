
package common

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/scylladb/go-set"
	"github.com/scylladb/go-set/strset"
	"github.com/thinkeridea/go-extend/exbytes"
)

// 是否空字符串
func IsEmptyString(str string) bool {
	return str == "null" || strings.TrimSpace(str) == ""
}

// 数组中是否包含某元素
func IsContainString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 获取数组元素下标
func GetStringSliceElementIndex(items []string, item string) int {
	for idx, eachItem := range items {
		if eachItem == item {
			return idx
		}
	}
	return -1
}

// 两个数组对比，判断 checkItems 是否是 originItems 子集
func IsSubsetString(originItems, checkItems []string) (bool, []string) {
	s1 := set.NewStringSet()
	for _, t := range originItems {
		s1.Add(strings.ToUpper(t))
	}
	s2 := set.NewStringSet()
	for _, t := range checkItems {
		s2.Add(strings.ToUpper(t))
	}
	isSubset := s1.IsSubset(s2)
	var notExists []string
	if !isSubset {
		notExists = strset.Difference(s2, s1).List()
	}
	return isSubset, notExists
}

// 过滤排除元素，返回新数组
func FilterDifferenceStringItems(originItems, excludeItems []string) []string {
	s1 := set.NewStringSet()
	for _, t := range originItems {
		s1.Add(strings.ToUpper(t))
	}
	s2 := set.NewStringSet()
	for _, t := range excludeItems {
		s2.Add(strings.ToUpper(t))
	}
	return strset.Difference(s1, s2).List()
}

// 过滤两个数组相同元素（交集），返回新数组
func FilterIntersectionStringItems(originItems, newItems []string) []string {
	s1 := set.NewStringSet()
	for _, t := range originItems {
		s1.Add(strings.ToUpper(t))
	}
	s2 := set.NewStringSet()
	for _, t := range newItems {
		s2.Add(strings.ToUpper(t))
	}
	return strset.Intersection(s1, s2).List()
}

// 删除字符中最后一个字母
func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

// 判断字符是否是数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// 字符转换
func StrconvIntBitSize(s string, bitSize int) (int64, error) {
	i, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return i, err
	}
	return i, nil
}

func StrconvUintBitSize(s string, bitSize int) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return i, err
	}
	return i, nil
}

func StrconvFloatBitSize(s string, bitSize int) (float64, error) {
	i, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return i, err
	}
	return i, nil
}

func StrconvRune(s string) (int32, error) {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return rune(r), err
	}
	return rune(r), nil
}

// 替换字符串引号字符
func ReplaceQuotesString(s string) string {
	return string(exbytes.Replace([]byte(s), []byte("\""), []byte(""), -1))
}

// 替换指定字符
func ReplaceSpecifiedString(s string, oldStr, newStr string) string {
	return string(exbytes.Replace([]byte(s), []byte(oldStr), []byte(newStr), -1))
}

// 忽略大小写切分字符串
func ReSplit(text string, cut string) []string {
	pattern := StringsBuilder("(?i)", cut)
	regex := regexp.MustCompile(pattern)
	result := regex.Split(text, -1)
	return result
}

// 字符数组转字符
func StringArrayToCapitalChar(strs []string) string {
	var newStrs []string
	for _, s := range strs {
		newStrs = append(newStrs, StringsBuilder("'", strings.ToUpper(s), "'"))
	}
	return strings.Join(newStrs, ",")
}

// 字符串拼接
func StringsBuilder(str ...string) string {
	var b strings.Builder
	for _, p := range str {
		b.WriteString(p)
	}
	return b.String() // no copying
}

// 字符串大写
func StringUPPER(str string) string {
	return strings.ToUpper(str)
}

// 字符串 JOIN
func StringJOIN(strs []string, strPrefix, strSuffix, joinS string) string {
	var tmpStr []string
	switch {
	case strPrefix == "" && strSuffix == "":
		for _, s := range strs {
			tmpStr = append(tmpStr, s)
		}
	case strPrefix != "" && strSuffix == "":
		for _, s := range strs {
			tmpStr = append(tmpStr, strPrefix+s)
		}
	case strPrefix == "" && strSuffix != "":
		for _, s := range strs {
			tmpStr = append(tmpStr, s+strSuffix)
		}
	default:
		for _, s := range strs {
			tmpStr = append(tmpStr, strPrefix+s+strSuffix)
		}
	}
	return strings.Join(tmpStr, joinS)
}

// 数组拆分
func SplitMultipleStringSlice(arr [][]string, num int64) [][][]string {
	var segmens = make([][][]string, 0)
	if num == 0 {
		segmens = append(segmens, arr)
		return segmens
	}
	max := int64(len(arr))
	if max < num {
		quantity := max / max
		end := int64(0)
		for i := int64(1); i <= max; i++ {
			qu := i * quantity
			if i != max {
				segmens = append(segmens, arr[i-1+end:qu])
			} else {
				segmens = append(segmens, arr[i-1+end:])
			}
			end = qu - i
		}
		return segmens
	}
	quantity := max / num
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}

// 多重数组拆分
func SplitMultipleSlice(arr []interface{}, num int64) [][]interface{} {
	var segmens = make([][]interface{}, 0)
	if num == 0 {
		segmens = append(segmens, arr)
		return segmens
	}
	max := int64(len(arr))
	if max < num {
		quantity := max / max
		end := int64(0)
		for i := int64(1); i <= max; i++ {
			qu := i * quantity
			if i != max {
				segmens = append(segmens, arr[i-1+end:qu])
			} else {
				segmens = append(segmens, arr[i-1+end:])
			}
			end = qu - i
		}
		return segmens
	}
	quantity := max / num
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}

// 用于版本号比较
func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}

// 用于对比 struct 是否相等
func DiffStructArray(structA, structB interface{}) ([]interface{}, []interface{}, bool) {
	var (
		addDiffs    []interface{}
		removeDiffs []interface{}
	)
	aVal := reflect.ValueOf(structA)
	bVal := reflect.ValueOf(structB)

	if !aVal.IsValid() && !bVal.IsValid() {
		return addDiffs, removeDiffs, true
	}

	if aVal.Kind() == reflect.Struct && bVal.Kind() == reflect.Struct {
		if !reflect.DeepEqual(structA, structB) {
			addDiffs = append(addDiffs, structA)
		}
		return addDiffs, removeDiffs, false
	}

	if aVal.IsNil() && bVal.IsNil() {
		return addDiffs, removeDiffs, true
	}

	if aVal.IsNil() && !bVal.IsNil() {
		if bVal.Kind() == reflect.Slice || bVal.Kind() == reflect.Array {
			for bi := 0; bi < bVal.Len(); bi++ {
				removeDiffs = append(removeDiffs, bVal.Index(bi).Interface())
			}
		}
	}

	if !aVal.IsNil() && bVal.IsNil() {
		if aVal.Kind() == reflect.Slice || aVal.Kind() == reflect.Array {
			for ai := 0; ai < aVal.Len(); ai++ {
				addDiffs = append(addDiffs, aVal.Index(ai).Interface())
			}
		}
	}

	if !aVal.IsNil() && !bVal.IsNil() {
		if (aVal.Kind() == reflect.Slice && bVal.Kind() == reflect.Slice) || (aVal.Kind() == reflect.Array && bVal.Kind() == reflect.Array) {
			dict := make(map[interface{}]bool)
			for bi := 0; bi < bVal.Len(); bi++ {
				dict[bVal.Index(bi).Interface()] = true
			}
			for ai := 0; ai < aVal.Len(); ai++ {
				if _, ok := dict[aVal.Index(ai).Interface()]; !ok {
					addDiffs = append(addDiffs, aVal.Index(ai).Interface())
				}
			}
		}
	}
	if len(addDiffs) == 0 && len(removeDiffs) == 0 {
		return addDiffs, removeDiffs, true
	}
	return addDiffs, removeDiffs, false
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 如果存在特殊字符，直接在特殊字符前添加\
/**
判断是否为字母： unicode.IsLetter(v)
判断是否为十进制数字： unicode.IsDigit(v)
判断是否为数字： unicode.IsNumber(v)
判断是否为空白符号： unicode.IsSpace(v)
判断师傅是特殊符号：unicode.IsSymbol(v)
判断是否为Unicode标点字符 :unicode.IsPunct(v)
判断是否为中文：unicode.Han(v)
*/
func SpecialLettersUsingMySQL(bs []byte) string {

	var (
		b     strings.Builder
		chars []rune
	)
	for _, r := range bytes.Runes(bs) {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r) {
			// mysql/tidb % 字符, /% 代表 /%，% 代表 % ,无需转义
			// mysql/tidb _ 字符, /_ 代表 /_，_ 代表 _ ,无需转义
			if r == '%' || r == '_' {
				chars = append(chars, r)
			} else {
				chars = append(chars, '\\', r)
			}
		} else {
			chars = append(chars, r)
		}
	}

	b.WriteString(string(chars))

	return b.String()
}

func SpecialLettersUsingOracle(bs []byte) string {

	var (
		b     strings.Builder
		chars []rune
	)
	for _, r := range bytes.Runes(bs) {
		if r == '\'' {
			chars = append(chars, '\'', r)
		} else {
			chars = append(chars, r)
		}
	}

	b.WriteString(string(chars))

	return b.String()
}

// 判断文件夹是否存在，不存在则创建
func PathExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("file dir MkdirAll failed: %v", err)
		}
	}
	return err
}
