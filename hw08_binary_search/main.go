package main

import (
	"fmt"
)

func main() {
	array := []int{1, 3, 5, 7, 100, 150}
	fmt.Println(binarySearchRecursive(array, 0, len(array)-1, 7))
	fmt.Println(binarySearchIterative(array, 7))
}

func binarySearchRecursive(array []int, left, right, target int) int {
	if left > right {
		return -1
	}

	mid := left + (right-left)/2

	if array[mid] == target {
		return mid
	}

	if array[mid] < target {
		return binarySearchRecursive(array, mid+1, right, target)
	}

	return binarySearchRecursive(array, left, mid-1, target)
}

func binarySearchIterative(array []int, target int) int {
	left := 0
	right := len(array) - 1

	for left <= right {
		mid := left + (right-left)/2

		if array[mid] == target {
			return mid
		}

		if array[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}
