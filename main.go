package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"time"

	"io/ioutil"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/viper"
)

type Registry struct {
	RegistryID   string `yaml:"registryId"`
	Region       string `yaml:"region"`
	PasswordFile string `yaml:"passwordFile"`
}

var configFile = flag.String("config", "/opt/config/ecr-token-refresh/config.yaml", "Configuration file location.")

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {

	flag.Parse()

	viper.SetConfigFile(*configFile)

	err := viper.ReadInConfig()

	exitOnError(err)

	err = os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	exitOnError(err)
}

func refreshCredential(registry Registry, sess *session.Session) {
	fmt.Printf("Refreshing credentials for %s in region %s.\n", registry.RegistryID, registry.Region)
	instance := ecr.New(sess, &aws.Config{
		Region: aws.String(registry.Region),
	})

	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(registry.RegistryID)},
	}

	authToken, err := instance.GetAuthorizationToken(input)

	if err != nil {
		log.Printf("Failed to get credential for %s in region %s (%s)", registry.RegistryID, registry.Region, err.Error())
		return
	}

	for _, data := range authToken.AuthorizationData {
		output, err := base64.StdEncoding.DecodeString(*data.AuthorizationToken)

		if err != nil {
			log.Printf("Failed to decode credential for %s in region %s (%s)", registry.RegistryID, registry.Region, err.Error())
			return
		}

		var password string
		split := strings.Split(string(output), ":")

		if len(split) == 2 {
			password = strings.TrimSpace(split[1])
		} else {
			log.Print("Failed to parse password.")
			return
		}

		err = ioutil.WriteFile(registry.PasswordFile, []byte(password), 0644)

		if err != nil {
			log.Printf("Failed to write password to file %s (%s)", registry.PasswordFile, err.Error())
		}

	}
}


func handleHttp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ECR Token Refresh")
}

func main() {

	sess, err := session.NewSession()

	exitOnError(err)

	var registries []Registry

	err = viper.UnmarshalKey("registries", &registries)

	exitOnError(err)

	interval, err := time.ParseDuration(viper.GetString("interval"))

	exitOnError(err)

	http.HandleFunc("/", handleHttp)
	go http.ListenAndServe(":3277", nil)

	if registries != nil {

		for _, registry := range registries {
			refreshCredential(registry, sess)
		}

		fmt.Printf("Starting periodic refresh of %d credentials every %s\n", len(registries), interval)
		c := time.Tick(interval)
		for range c {
			for _, registry := range registries {
				go refreshCredential(registry, sess)
			}
		}
	}

}
