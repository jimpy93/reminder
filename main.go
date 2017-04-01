package main

import(
    "github.com/jimpy93/scheduler"
    "time"
    "fmt"
    "bufio"
    "os"
    "strings"
)

func addAlarm(q *scheduler.Scheduler){
    var hour, min int
    fmt.Print("Set Alarm for (use format HH:MM where HH varies from 0-24): ")
    fmt.Scanf("%d:%d", &hour, &min)
    cur_time := time.Now()
    t := cur_time
    t = t.Add(-time.Duration(t.Hour()) * time.Hour - time.Duration(t.Minute()) * time.Minute)
    t = t.Add(time.Duration(hour) * time.Hour + time.Duration(min) * time.Minute)
    t = t.Add(-time.Duration(t.Second()) * time.Second)
    if cur_time.Hour() > hour && cur_time.Minute() > min{
        t = t.Add(24 * time.Hour)
    }
    rem := reminder{isAlarm: true, msg: "Alarm"}
    q.Add(scheduler.Task{Time: t, Task: rem})
}

func addReminder(q *scheduler.Scheduler){
    var date, month, year, hour, min int
    fmt.Print("Reminder msg: ")
    in := bufio.NewReader(os.Stdin)
    msg, err := in.ReadString('\n')
    if err != nil{
        fmt.Errorf("Unable to read msg. Got: %s\n", err)
        return
    }
    msg = strings.TrimRight(msg, "\n")
    fmt.Print("Set Reminder for (use format DD-MM-YYYY HH:MM, where HH varies from 0-24): ")
    fmt.Scanf("%d-%d-%d %d:%d", &date, &month, &year, &hour, &min)
    t := time.Date(year, time.Month(month), date, hour, min, 0, 0, time.Now().Location())
    rem := reminder{isAlarm: false, msg: msg}
    q.Add(scheduler.Task{Time: t, Task: rem})
}

func removeAlarm(q *scheduler.Scheduler){
    var index int
    for i,t := range q.Tasks(){
        var rem reminder
        if r, ok := t.Task.(reminder); ok {
            rem = r
        } else {
            rem = reminder {msg: "Unable to parse this reminder"}
        }
        fmt.Printf("%d. set for %s, msg: %s\n", i, t.Time, rem.msg)
    }
    fmt.Print("Choose an alarm/reminder to remove: ")
    fmt.Scanf("%d", &index)
    suc := q.RemoveFromPosition(index-1)
    if suc{
        fmt.Println("Successfully removed!")
    } else {
        fmt.Println("Not removed! Please check that you provided an valid index")
    }
}

func consumeTasks(q *scheduler.Scheduler, exitChan chan bool){
    exit := false
    for !exit && q.IsRunning(){
        select{
        case v:= <- q.TriggerChan() :
            if r, ok := v.Task.(reminder); ok{
                if r.isAlarm {
                    TriggerAlarm()
                } else{
                    TriggerReminder(r.msg)
                }
            } else {
                fmt.Errorf("Unable to detect the task triggered!!!\n")
            }
        case exit = <- exitChan:
        }
    }
}

func main() {
    q := scheduler.NewScheduler();
    q.Start()
    menu := `
    Choose an option:
    1. Add an alarm
    2. Add an reminder
    3. Remove an Alarm/Reminder
    4. Stop all alarms/reminders and Exit!
    `
    exit := false
    exitChan := make(chan bool)
    go consumeTasks(&q, exitChan)
    for !exit{
        var opt int
        fmt.Println(menu)
        fmt.Scanf("%d", &opt)
        switch opt{
        case 1:
            addAlarm(&q)
        case 2:
            addReminder(&q)
        case 3:
            removeAlarm(&q)
        case 4:
            exit = true
            exitChan <- true
            q.Stop()
        default:
            fmt.Println("Please choose an valid option!")
        }

    }
}
