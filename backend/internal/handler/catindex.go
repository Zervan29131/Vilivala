package handler

import (
	"sort"

	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

const maxCategoryDepth = 3

// catIndex is an in-memory view of the (small) category table.
type catIndex struct {
	byID     map[uint]*model.Category
	children map[uint][]uint
	roots    []uint
}

func loadCats() *catIndex {
	var list []model.Category
	config.DB.Order("sort ASC, id ASC").Find(&list)
	ix := &catIndex{byID: map[uint]*model.Category{}, children: map[uint][]uint{}}
	for i := range list {
		c := &list[i]
		ix.byID[c.ID] = c
	}
	for i := range list {
		c := &list[i]
		if c.ParentID == 0 || ix.byID[c.ParentID] == nil {
			ix.roots = append(ix.roots, c.ID)
		} else {
			ix.children[c.ParentID] = append(ix.children[c.ParentID], c.ID)
		}
	}
	return ix
}

// chain returns the categories from the root down to id.
func (ix *catIndex) chain(id uint) []*model.Category {
	var rev []*model.Category
	seen := map[uint]bool{}
	for id != 0 && !seen[id] {
		seen[id] = true
		c := ix.byID[id]
		if c == nil {
			break
		}
		rev = append(rev, c)
		id = c.ParentID
	}
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}

func (ix *catIndex) depth(id uint) int { return len(ix.chain(id)) }

// descendants returns id and all its descendant ids.
func (ix *catIndex) descendants(id uint) []uint {
	out := []uint{id}
	for _, c := range ix.children[id] {
		out = append(out, ix.descendants(c)...)
	}
	return out
}

// height is the number of levels in the subtree rooted at id (leaf = 1).
func (ix *catIndex) height(id uint) int {
	h := 0
	for _, c := range ix.children[id] {
		if x := ix.height(c); x > h {
			h = x
		}
	}
	return h + 1
}

func (ix *catIndex) ref(id uint) *model.CategoryRef {
	ch := ix.chain(id)
	if len(ch) == 0 {
		return nil
	}
	leaf := ch[len(ch)-1]
	r := &model.CategoryRef{ID: leaf.ID, Name: leaf.Name, Slug: leaf.Slug}
	for _, c := range ch {
		r.Path = append(r.Path, model.PathNode{ID: c.ID, Name: c.Name, Slug: c.Slug})
	}
	return r
}

func (ix *catIndex) pathString(id uint) string {
	s := ""
	for i, c := range ix.chain(id) {
		if i > 0 {
			s += " / "
		}
		s += c.Name
	}
	return s
}

// fillCategories sets CategoryOut (with breadcrumb) on each post.
func fillCategories(posts []model.Post) {
	ix := loadCats()
	for i := range posts {
		if posts[i].CategoryID != nil {
			posts[i].CategoryOut = ix.ref(*posts[i].CategoryID)
		}
	}
}

func sortedKeysDesc(m map[int]bool) []int {
	k := make([]int, 0, len(m))
	for x := range m {
		k = append(k, x)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(k)))
	return k
}
