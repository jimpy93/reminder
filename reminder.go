package main

import(
    "github.com/0xAX/notificator"
)

var basepath string = "./resources/"

var notify *notificator.Notificator = notificator.New(notificator.Options{
                                        DefaultIcon: "icon/default.png",
                                        AppName:     "Reminder",
                                      })


type reminder struct{
    isAlarm bool
    msg string
}

func TriggerReminder(msg string){
    notify.Push("Reminder", msg, basepath + "alarm_icon.png", notificator.UR_CRITICAL)
    playWave(basepath + "reminder.wav")
}

func TriggerAlarm(){
    notify.Push("Reminder", "Alarm!", basepath + "alarm_icon.png", notificator.UR_CRITICAL)
    playWave(basepath + "alarm.wav")
}
