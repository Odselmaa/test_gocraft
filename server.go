package main

import (
  "net/http"
  "log"

  "github.com/test/lib"
  "github.com/gocraft/work"
)

var enqueuer = work.NewEnqueuer("work", lib.RedisPool)
var task_name = "test"
var task_finished  = false
var task_sent = false

func startTask(){
  _, err := enqueuer.Enqueue(task_name, work.Q{"message": "ping"})
  if err != nil {
    log.Fatal(err)
  }
}

//endpoint for receiving task result
func taskFinished(w http.ResponseWriter, r *http.Request){
   task_finished = true;
   w.Write([]byte("OK"))
}

//endpoint for informing state of task
func sayHello(w http.ResponseWriter, r *http.Request) {
  var msg = "WAIT"
  if(!task_sent){
    startTask()
    task_sent = true
  }

  if(task_finished){
    msg = "PONG"
    task_sent = false
    task_finished = false
  }
  w.Write([]byte(msg))
}

func main() {
  http.HandleFunc("/hello", sayHello)
  http.HandleFunc("/finish", taskFinished)

  if err := http.ListenAndServe(":8080", nil); err != nil {
      panic(err)
  }
}
