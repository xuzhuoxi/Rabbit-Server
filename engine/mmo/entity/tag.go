// Package entity
// Create on 2023/8/30
// @author xuzhuoxi
package entity

import "github.com/xuzhuoxi/infra-go/slicex"

type TagsSupport struct {
	Tags []string
}

func (o *TagsSupport) SetTags(tags []string) {
	o.Tags = tags
}

func (o *TagsSupport) AddTag(tag string, tags ...string) {
	o.addTag(tag)
	if len(tags) > 0 {
		for index := range tags {
			o.addTag(tags[index])
		}
	}
}

func (o *TagsSupport) AddTags(tags []string) {
	if len(tags) == 0 {
		return
	}
	for index := range tags {
		o.addTag(tags[index])
	}
}

func (o *TagsSupport) RemoveTag(tag string) {
	if len(tag) == 0 {
		return
	}
	o.removeTag(tag)
}

func (o *TagsSupport) RemoveTags(tags []string) {
	if len(tags) == 0 {
		return
	}
	for index := range tags {
		o.removeTag(tags[index])
	}
}

func (o *TagsSupport) ContainsTag(tag string) bool {
	return slicex.ContainsString(o.Tags, tag)
}

func (o *TagsSupport) ContainsTags(tags []string, and bool) bool {
	if and {
		for index := range tags {
			if !slicex.ContainsString(o.Tags, tags[index]) {
				return false
			}
		}
		return true
	} else {
		for index := range tags {
			if slicex.ContainsString(o.Tags, tags[index]) {
				return true
			}
		}
		return false
	}
}

func (o *TagsSupport) addTag(tag string) {
	if len(tag) == 0 {
		return
	}
	if len(o.Tags) == 0 {
		o.Tags = append(o.Tags, tag)
		return
	}
	if slicex.ContainsString(o.Tags, tag) {
		return
	}
	o.Tags = append(o.Tags, tag)
}

func (o *TagsSupport) removeTag(tag string) {
	if len(tag) == 0 {
		return
	}
	if len(o.Tags) == 0 {
		return
	}
	for index := len(o.Tags) - 1; index >= 0; index-- {
		if o.Tags[index] == tag {
			o.Tags = append(o.Tags[:index], o.Tags[index+1:]...)
			return
		}
	}
}
