package l337c0d3

func twoSum(nums []int, target int) []int {
	return question_1_impl_1(nums, target)
}

func question_1_impl_1(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}

	return nil
}

func question_1_impl_2(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		numMap[nums[i]] = i
	}
	for i := 0; i < len(nums); i++ {
		complement := target - nums[i]
		if idx, ok := numMap[complement]; ok {
			return []int{i, idx}
		}
	}

	return nil
}

func question_1_impl_3(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		complement := target - nums[i]
		if idx, ok := numMap[complement]; ok {
			return []int{idx, i}
		}
		numMap[nums[i]] = i
	}

	return nil
}
