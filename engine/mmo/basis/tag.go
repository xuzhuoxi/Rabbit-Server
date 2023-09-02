// Package basis
// Create on 2023/8/31
// @author xuzhuoxi
package basis

type ITagsSupport interface {
	// SetTags 批量设置Tag
	SetTags(tags []string)
	// AddTag 添加Tag
	AddTag(tag string, tags ...string)
	// AddTags 批量添加Tag
	AddTags(tags []string)
	// RemoveTag 移除Tag
	RemoveTag(tag string)
	// RemoveTags 批量移除Tag
	RemoveTags(tags []string)
	// ContainsTag 判断是否包含Tag
	ContainsTag(tag string) bool
	// ContainsTags 判断是否包含多个Tag
	// and: true代表与关系， false代表或关系
	ContainsTags(tags []string, and bool) bool
}
