package lib

import (
	"fmt"
	"time"	
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
)

var URL = "http://localhost:8080/finish"

type Context struct {
	Message string
}

var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var Pool = work.NewWorkerPool(Context{}, 10, "work", RedisPool)
var client = work.NewClient("work", RedisPool)


func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) FindMessage(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["message"]; ok {
		c.Message = job.ArgString("message")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	return next()
}
//send http request to inform that task has been processed
func sendResult(){
	resp, _ := http.Get(URL)
	fmt.Println("Response: ", resp)

}

// func IsBusy(name string) bool {
// 	observations, _ := client.WorkerObservations()
// 	var busy = true;
// 	for _, ob := range observations {
// 	      fmt.Println("Checkin: ", ob.IsBusy, ob.JobName , ob.JobID, ob.StartedAt)
// 	      if(name==ob.JobName){
// 	      	busy = ob.IsBusy
// 	      	break
// 	      }
// 	 }
// 	 return busy
// }

// process task
func (c *Context) SendMessage(job *work.Job) error {
	msg := job.ArgString("message")
	if err := job.ArgError(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	fmt.Println("message: ", msg)
	fmt.Println("result:  pong")
	sendResult()
	return nil
}

