// Package basis
// Create on 2023/8/31
// @author xuzhuoxi
package basis

type ITagsSupport interface {
	SetTags(tags []string)
	AddTag(tag string, tags ...string)
	AddTags(tags []string)
	RemoveTag(tag string)
	RemoveTags(tags []string)
	ContainsTag(tag string) bool
	ContainsTags(tags []string, and bool) bool
}
