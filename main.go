package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Build struct {
	Timestamp int64 `json:"timestamp"`
}

type JenkinsJob struct {
	Builds []Build `json:"builds"`
}

type Config struct {
	JenkinsURL string   `yaml:"jenkinsURL"`
	Username   string   `yaml:"username"`
	Token      string   `yaml:"token"`
	Jobs       []string `yaml:"jobs"`
}

func main() {
	config := Config{}

	// 从配置文件中读取配置
	file, err := ioutil.ReadFile("conf/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	totalBuildCount := 0
	thisMonth := time.Now().Month()

	for _, jobName := range config.Jobs {
		// 创建 HTTP 请求
		req, err := http.NewRequest("GET", config.JenkinsURL+"/job/"+jobName+"/api/json?tree=builds[*]", nil)
		if err != nil {
			log.Fatal(err)
		}

		// 设置 HTTP Basic Authentication
		req.SetBasicAuth(config.Username, config.Token)

		// 发送 HTTP 请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 解析响应内容
		job := JenkinsJob{}
		err = json.Unmarshal(body, &job)
		if err != nil {
			log.Fatal(err)
		}

		// 计算本月的构建数
		monthlyBuildCount := 0
		for _, build := range job.Builds {
			buildTime := time.Unix(build.Timestamp/1000, 0)
			if buildTime.Month() == thisMonth {
				monthlyBuildCount++
			}
		}

		fmt.Printf("Job %s 的本月构建数: %d\n", jobName, monthlyBuildCount)
		totalBuildCount += monthlyBuildCount
	}

	fmt.Printf("总的本月构建数: %d\n", totalBuildCount)
}
