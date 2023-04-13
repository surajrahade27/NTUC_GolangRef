package util

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToDateTime(t *testing.T) {
	t.Run("When valid date string is passed, it should return date", func(t *testing.T) {
		inputStr := time.Now().Format("2006-01-02 15:04:05")
		actualValue, err := ToDateTime(inputStr)
		ShouldEqual(actualValue, `2023-01-05 12:00:00 +0000 UTC`)
		ShouldBeNil(err)
	})
	t.Run("When invalid date string is passed, it should return error", func(t *testing.T) {
		_, err := ToDateTime(time.Now().String())
		ShouldNotBeNil(err)
	})

	t.Run("When empty date string is passed, it should return empty date and nil error", func(t *testing.T) {
		dateValue, err := ToDateTime("")
		ShouldNotBeNil(err, nil)
		ShouldEqual(dateValue, time.Time{})
	})
}

func Test_Difference(t *testing.T) {
	t.Run("test scenario", func(t *testing.T) {
		array1 := []int64{1, 2, 3, 4}
		array2 := []int64{2, 3, 4}
		expectedRes := []int64{1}
		response := Difference(array1, array2)
		ShouldEqual(response, expectedRes)
	})
}

func Test_GetUserID(t *testing.T) {
	t.Run("when user id is nil", func(t *testing.T) {
		ctx := context.Background()
		userID, err := GetUserID(ctx)
		expectedErr := "invalid user id"
		ShouldNotBeNil(err)
		ShouldEqual(userID, 0)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("when user id of interger type is given", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userId", 12345)
		userID, err := GetUserID(ctx)
		ShouldEqual(userID, 12345)
		ShouldBeNil(err)
	})

	t.Run("when user id of string type is given", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userId", "12345")
		userID, err := GetUserID(ctx)
		ShouldEqual(userID, 12345)
		ShouldBeNil(err)
	})

	t.Run("when user id of invalid type is given", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userId", int64(12345))
		userID, err := GetUserID(ctx)
		expectedErr := "invalid user id"
		ShouldNotBeNil(err)
		ShouldEqual(userID, 0)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
}
