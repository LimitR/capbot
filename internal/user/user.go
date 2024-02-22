package user

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	MAX = 100
)

type User struct {
	Id        int64
	ChatId    int64
	MessageId int
	Nums      map[string]int
	num1      int
	num2      int
	sum       int
	Time      time.Time
}

func NewUser(userId int64, chatId int64) *User {
	num1 := rand.Intn(MAX)
	num2 := rand.Intn(MAX)
	return &User{
		Id:     userId,
		ChatId: chatId,
		Nums:   getRandSum(num1 + num2),
		sum:    num1 + num2,
		num1:   num1,
		num2:   num2,
		Time:   time.Now(),
	}
}

func (u *User) GetString() string {
	return fmt.Sprintf("%d + %d = ?", u.num1, u.num2)
}

func (u *User) Validate(num string) bool {
	inputNum, ok := u.Nums[num]
	if !ok {
		return false
	}
	if inputNum == u.sum {
		return true
	}

	return false
}

func getRandSum(sum int) map[string]int {
	res := make(map[string]int, 4)
	m := make(map[int]struct{}, 4)
	s := false
	counter := 0
	for {
		num := rand.Intn(4)
		if !s {
			name := strconv.Itoa(num)
			res[name] = sum
			s = true
			m[num] = struct{}{}
			counter++
			continue
		}
		_, ok := m[num]
		if !ok {
			m[num] = struct{}{}
			name := strconv.Itoa(num)
			res[name] = rand.Intn(MAX + MAX)
			counter++
		}
		if counter > 3 {
			break
		}
	}
	return res
}
