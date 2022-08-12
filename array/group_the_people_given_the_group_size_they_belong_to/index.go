package group_the_people_given_the_group_size_they_belong_to

func GroupThePeople(groupSizes []int) (ans [][]int) {
	groups := map[int][]int{}
	for i, size := range groupSizes {
		groups[size] = append(groups[size], i)
	}
	for size, people := range groups {
		for i := 0; i < len(people); i += size {
			ans = append(ans, people[i:i+size])
		}
	}
	return
}
