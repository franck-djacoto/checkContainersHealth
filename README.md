# CHECK CONTAINERS HEALTH
This simple project check if containers deployed on certains url respond well and are not in a state like <br/>
**Gateway time out** or are not returning a 404 error for example. If an abnormal behavior is <br/>
returned by a certain containers, then a mail is sent to some people who can fix this. <br/>
***

# WHY THIS PROJECT  ?
I decided to code this simple program to solve a problem that we (my classroom mates and I) encountered during <br/>
our end of school year project. We developped a workflow devops containing many containers which sometimes go on a <br/>
**Gateway time out** status . To avoid checking containers state every time . I decided to ask my <br/>
program to do that for me and send an email alert to me when it detect a unexpected behavior.<br/>

***

# HOW TO RUN THE PROJECT ?
## Configuration
1. Create a `.env` file at the root of the porject and copy the content of `.env.example`inside it and set the values of env variables 
2. For those containers that may be checked via a connexion to an api endpoint for example, provide credentials for connexion in `.env` file
Take example on `PROD_AMDIN_ID` and `PROD_ADMIN_PASS`
3. If you want certains `env variable not to be null` make sure to add them to `notNullableEnvVar`variable in `main.go` file
4. Make sure to configure env varibales `SMTP_HOST,SMTP_PORT, SENDER_MAIL, SENDER_PASS, RECEIVERS` to receive a mail. You can have multiple receivers seperated with `,` 

## Run the project 
**NB** : When coding this project, I projected to run it as  cron job on my system. It's the best way so that you may not have to run it manually all the time. <br/>
<p>At the root of the project run <code>go run main.go</code> to launch the app. It'll display the logs about containers'health in the terminal.</p>