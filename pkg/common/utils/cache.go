package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func GetUserLogin(client *redis.Client, channel string, loginId string) (error, LoginData) {
	result := client.Get(fmt.Sprintf("LOGIN:%s:%s", channel, loginId))
	value, err := result.Result()
	if err != nil {
		return err, LoginData{}
	}

	var userInfo LoginData
	err = json.Unmarshal([]byte(value), &userInfo)
	if err != nil {
		return err, LoginData{}
	}
	if userInfo.MDN == "" {
		return errors.New("User data Not Found/Empty"), LoginData{}
	}
	return nil, userInfo
}

func SetData(client *redis.Client, key string, value LoginData) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	diff := endOfDay.Sub(now)
	diffSecond := int(diff.Seconds())

	jsonData, err := json.Marshal(value)
	err = client.Set(key, jsonData, time.Second*time.Duration(diffSecond)).Err()
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

type LoginData struct {
	UserId                 string    `json:"user_id"`
	LoginId                string    `json:"login_id"`
	MDN                    string    `json:"MDN"`
	PartnerId              int       `json:"partner_id"`
	PartnerCode            string    `json:"partner_code"`
	UserGroupId            int       `json:"user_group_id"`
	ClusterId              int       `json:"cluster_id"`
	ClusterName            string    `json:"cluster_name"`
	ClusterCode            string    `json:"cluster_code"`
	RegionId               int       `json:"region_id"`
	RegionName             string    `json:"region_name"`
	RegionCode             string    `json:"region_code"`
	LogonAttempt           int       `json:"logon_attempt"`
	IsLogin                int       `json:"is_login"`
	IsAllowedTime          int       `json:"is_allowed_time"`
	Status                 int       `json:"status"`
	ForgotPassword         int       `json:"forgot_password"`
	LastTimeChangePassword time.Time `json:"last_time_change_password"`
	ChangePin              int       `json:"change_pin"`
	ExpiredTime            int       `json:"expired_time"`
	LastActivity           string    `json:"last_activity"`
}
