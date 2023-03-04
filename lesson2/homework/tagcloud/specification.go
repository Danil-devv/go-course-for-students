package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	// Словарь состоит из пар вида тег-индекс
	// Тег является строкой
	// map[tag] - это индекс вхождения тега tag в слайс occurrence
	tagIndexInOccurrenceArray map[string]int

	// Слайс частоты появления тегов, отсортированный по убыванию
	// occurrence[i] - это частота появления тега с индексом i
	// ПРИМЕР: occurrence[tagIndexInOccurrenceArray["t1"]] - частота появления тега "t1"
	occurrence []int

	// Слайс тегов
	// indexToTag[i] возвращает тег(строку), частотность которого равна occurrence[i]
	// ПРИМЕР: indexToTag[5] - тег, который встречается occurrence[5] раз
	indexToTag []string

	// Данный метод хранения тегов необходим для того, чтобы не хранить
	// одинаковые теги по несколько экземпляров и быстро добавлять тег в облако,
	// а также быстро находить N самых частотных тегов

	// O(log N) - асимптотика функции AddTag(tag)
	// O(n) - асимптотика функции TopN(n)
}

// TagStat represent statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{
		tagIndexInOccurrenceArray: make(map[string]int),
		occurrence:                make([]int, 0),
		indexToTag:                make([]string, 0),
	}
}

// AddTag add a tag to the cloud if it wasn't present and increase tag occurrence count
func (cloud *TagCloud) AddTag(tag string) {
	index, ok := cloud.tagIndexInOccurrenceArray[tag]
	if ok {
		cloud.occurrence[index]++

		// После увеличения частотности тега возможна такая ситуация,
		// что слайс occurrence больше не отсортирован по убыванию.
		// В таком случае есть ровно один элемент, с которым нам нужно
		// поменять местами occurrence[index], чтобы слайс снова стал остортированным.
		// Воспользуемся для этого бинарным поиском
		left, right := 0, index
		for left < right {
			mid := (left + right) / 2
			if cloud.occurrence[mid] >= cloud.occurrence[index] {
				left = mid + 1
			} else {
				right = mid
			}
		}

		// Смена местами тега с индексом index и тега с индексом left
		cloud.occurrence[left], cloud.occurrence[index] = cloud.occurrence[index], cloud.occurrence[left]
		cloud.indexToTag[left], cloud.indexToTag[index] = cloud.indexToTag[index], cloud.indexToTag[left]
		cloud.tagIndexInOccurrenceArray[tag] = left
		cloud.tagIndexInOccurrenceArray[cloud.indexToTag[index]] = index

	} else {
		// Если тег встречается впервые, то его необходимо просто добавить в конец массива occurrence
		cloud.tagIndexInOccurrenceArray[tag] = len(cloud.occurrence)
		cloud.occurrence = append(cloud.occurrence, 1)
		cloud.indexToTag = append(cloud.indexToTag, tag)
	}
}

// TopN return top N most frequent tags ordered in descending order by occurrence count
func (cloud *TagCloud) TopN(n int) []TagStat {
	// Обработка случая, когда N больше кол-ва тегов в облаке
	if n > len(cloud.occurrence) {
		n = len(cloud.occurrence)
	}

	top := make([]TagStat, n)

	// Так как массив occurrence отсортирован по убыванию,
	// нам достаточно взять первые N тегов
	for i := 0; i < n; i++ {
		top[i].Tag = cloud.indexToTag[i]
		top[i].OccurrenceCount = cloud.occurrence[i]
	}

	return top
}
