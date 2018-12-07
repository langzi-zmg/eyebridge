# eyebridge

1.fluentd(get docker container log)-->kafka  
2.The message that consumes kafka will then send a message containing the keyword to the user.ex,panic(go),ERROR(java)  
3.Send to different people based on different service names  
4.According to the service name, if the mail is sent within a certain period of time, the mail will not be sent repeatedly.  
5.all data load into mysql  
