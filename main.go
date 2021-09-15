package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type smtpServer struct {
	host string
	port string
}

type Sender struct {
	email    string
	password string
}

func (s *smtpServer) Adress() string {
	return s.host + ":" + s.port
}

type Container struct {
	name          string
	url           string
	method        string
	identifiant   string
	password      string
	idFieldName   string
	passFieldName string
}

var containers map[string]Container
var sender Sender
var receivers []string
var server smtpServer
var notNullableEnvVar = []string{
	"GOGS_ID", "GOGS_PASS", "JENKINS_ID", "JENKINS_PASS", "GRAFANA_ID", "GRAFANA_PASS",
	"PORTAINER_ID", "PORTAINER_PASS", "PROD_AMDIN_ID", "PROD_ADMIN_PASS",
	"PREPROD_ADMIN_PASS", "PREPROD_FRONT_ID", "PREPROD_FRONT_PASS", "PROD_FRONT_ID", "PROD_FRONT_ID", "PREPROD_FRONT_PASS",
	"PHPMYADMIN_ID", "PHPMYADMIN_PASS", "SENDER_MAIL", "SENDER_PASS", "SMTP_HOST", "SMTP_PORT", "RECEIVERS",
}

func checkAndInitConfiguration() error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	for _, v := range notNullableEnvVar {
		if os.Getenv(v) == "" {
			return errors.New("Env variable " + v + "can't be null")
		}
	}

	containers = map[string]Container{
		"gogs": {name: "gogs",
			url:           "https://gogs.dsp-archiwebo20-mt-ma-ca-fd.fr/user/login",
			method:        "post",
			identifiant:   os.Getenv("GOGS_ID"),
			password:      os.Getenv("GOGS_PASS"),
			idFieldName:   "user_name",
			passFieldName: "password",
		},

		"jenkins": {name: "jenkins",
			url:           "https://jenkins.dsp-archiwebo20-mt-ma-ca-fd.fr/login?from=%2F",
			method:        "post",
			identifiant:   os.Getenv("JENKINS_ID"),
			password:      os.Getenv("JENKINS_PASS"),
			idFieldName:   "j_username",
			passFieldName: "j_password",
		},

		"grafana": {name: "grafana",
			url:           "https://grafana.dsp-archiwebo20-mt-ma-ca-fd.fr",
			method:        "post",
			identifiant:   os.Getenv("GRAFANA_ID"),
			password:      os.Getenv("GRAFANA_PASS"),
			idFieldName:   "user",
			passFieldName: "password",
		},

		"portainer": {name: "portainer",
			url:           "https://portainer.dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method:        "post",
			identifiant:   os.Getenv("PORTAINER_ID"),
			password:      os.Getenv("PORTAINER_PASS"),
			idFieldName:   "username",
			passFieldName: "password",
		},

		"prod-admin": {name: "prod-admin",
			url:           "https://prod-admin.dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method:        "post",
			identifiant:   os.Getenv("PRO_AMDIN_ID"),
			password:      os.Getenv("PRO_ADMIN_PASS"),
			idFieldName:   "email",
			passFieldName: "password",
		},

		"preprod-admin": {name: "preprod-admin",
			url:           "https://preprod-admin.dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method:        "post",
			identifiant:   os.Getenv("PREPROD_ADMIN_ID"),
			password:      os.Getenv("PREPROD_FRONT_PASS"),
			idFieldName:   "email",
			passFieldName: "password",
		},

		"preprod-frontend": {name: "preprod-frontend",
			url:    "https://frontend.dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method: "post", identifiant: os.Getenv("PREPROD_FRONT_ID"),
			password:      os.Getenv("PREPROD_FRONT_PASS"),
			idFieldName:   "email",
			passFieldName: "password",
		},

		"prod-frontend": {name: "prod-frontend",
			url:           "https://dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method:        "post",
			identifiant:   os.Getenv("PROD_FRONT_ID"),
			password:      os.Getenv("PREPROD_FRONT_PASS"),
			idFieldName:   "email",
			passFieldName: "password",
		},

		"php-my-amdin": {name: "php-my-admin",
			url:           "https://phpadmin.dsp-archiwebo20-mt-ma-ca-fd.fr/index.php?route=/",
			method:        "post",
			identifiant:   os.Getenv("PHPMYADMIN_ID"),
			password:      os.Getenv("PHPMYADMIN_PASS"),
			idFieldName:   "pma_username",
			passFieldName: "pma_password",
		},

		"traefik": {name: "traefik",
			url:    "https://traefik.dsp-archiwebo20-mt-ma-ca-fd.fr/",
			method: "get",
		},

		"cAdvisor": {name: "cAdvisor",
			url:    "https://cadvisor.dsp-archiwebo20-mt-ma-ca-fd.fr",
			method: "get",
		},

		"prometheus": {name: "prometheus",
			url:    "https://prometheus.dsp-archiwebo20-mt-ma-ca-fd.fr",
			method: "get",
		},
	}

	sender.email = os.Getenv("SENDER_MAIL")
	sender.password = os.Getenv("SENDER_PASS")
	server.host = os.Getenv("SMTP_HOST")
	server.port = os.Getenv("SMTP_PORT")
	receivers = strings.Split(os.Getenv("RECEIVERS"), ",")

	return nil
}

func checkContainerHealth(container Container) (error, *http.Response) {
	var resp *http.Response
	var err error
	if _, ok := containers[container.name]; container.name == "" || container.url == "" || !ok {
		return errors.New("Provide valid container name and valid container Url"), resp
	}

	if container.method == "post" {
		requestBody, err := json.Marshal(map[string]string{
			container.idFieldName:   container.identifiant,
			container.passFieldName: container.password,
		})

		if err != nil {
			return err, nil
		}

		fmt.Printf("Trying to connect to %s container", container.name)
		resp, err = http.Post(container.url, "application/json", bytes.NewBuffer(requestBody))
	}

	if container.method == "get" {
		resp, err = http.Get(container.url)
	}

	if err != nil {
		return err, resp
	}

	return nil, resp
}

func sendMailAlert(server smtpServer, sender Sender, receiver []string, message []byte) error {
	auth := smtp.PlainAuth("", sender.email, sender.password, server.host)
	err := smtp.SendMail(server.Adress(), auth, sender.email, receiver, message)

	if err != nil {
		return err
	}
	return nil
}

func main() {

	err := checkAndInitConfiguration()

	if err != nil {
		log.Fatalf("Error on configuring app %v", err)
	}

	for containerName, container := range containers {
		err, resp := checkContainerHealth(container)
		if err != nil {
			log.Fatalf("Could'nt get %s container health State\n", containerName)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("Container %s is healthy", containerName)
		}

		if resp.StatusCode >= 400 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error while retreiving %s container response body", containerName)
			}

			log.Printf("!!! Container %s is not responding as expected  !!! => status code : %d", containerName, resp.StatusCode)
			log.Printf("%s container's body : ==> %s", containerName, string(bodyBytes))

			message := []byte(
				"Subject: Problem found on container  " + containerName + " (code : " + resp.Status + ")" + "\r\n\r\n" +
					string(bodyBytes),
			)

			err = sendMailAlert(server, sender, receivers, message)

			if err != nil {
				fmt.Errorf("An error occured when sending mail alert for container %s => %v \n", containerName, err)
			}

			log.Printf("Mail alert send successfully for container : %s", containerName)
		}
	}

}
