package rexUtils

// note: 快速对比2个新旧列表，返回需要删除和添加的元素列表

func DiffStringSlice(oldList, newList []string) (toDelete, toAdd []string) {
	oldSet := make(map[string]struct{}, len(oldList))
	newSet := make(map[string]struct{}, len(newList))

	for _, v := range oldList {
		oldSet[v] = struct{}{}
	}
	for _, v := range newList {
		newSet[v] = struct{}{}
	}

	// 需要删除：old 有，new 没有
	for v := range oldSet {
		if _, ok := newSet[v]; !ok {
			toDelete = append(toDelete, v)
		}
	}

	// 需要添加：new 有，old 没有
	for v := range newSet {
		if _, ok := oldSet[v]; !ok {
			toAdd = append(toAdd, v)
		}
	}

	return
}
