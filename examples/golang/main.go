package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ajvb/kala/client"
	"github.com/ajvb/kala/job"
	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("Remote job scheduler test")
	e := echo.New()
	e.Logger.SetLevel(1)
	e.GET("/", func(c echo.Context) error {
		jobid := c.QueryParam("jobid")
		log.Println("Callback called timer expired ", jobid)
		return c.String(http.StatusOK, "Callback Called!!!")
	})

	kalaClient := client.New("http://localhost:8000")
	scheduleTime := time.Now().Add(time.Minute * 1)
	parsedTime := scheduleTime.Format(time.RFC3339)
	delay := "PT1M"
	scheduleStr := fmt.Sprintf("R0/%s/%s", parsedTime, delay)
	log.Println("Schedule string is : ", scheduleStr)
	body := &job.Job{
		Name:     "OneMinute Timer Once Only",
		Schedule: scheduleStr,
		JobType:  job.RemoteJob,
		RemoteProperties: job.RemoteProperties{
			Url:                   "http://localhost:1323/",
			Method:                http.MethodGet,
			ExpectedResponseCodes: []int{200},
		},
	}
	jobid, err := kalaClient.CreateJob(body)
	log.Println("jobid", jobid, err)
	e.Logger.Fatal(e.Start(":1323"))
}
