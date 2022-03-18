package wxws

func StrSearch(needle string, haystack []string) (bool, int) {
	for index, hay := range haystack {
		if hay == needle {
			return true, index
		}
	}
	return false, 0
}

func SliceRemove(slice []string, index int) []string {
	slice[len(slice)-1], slice[index] = "", slice[len(slice)-1]
	return slice[:len(slice)-1]
}
