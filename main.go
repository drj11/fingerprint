package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	client := ec2.New(session.New())

	instanceID := os.Args[1]

	resp, err := client.GetConsoleOutput(
		&ec2.GetConsoleOutputInput{
			InstanceId: aws.String(instanceID)})
	if err != nil {
		log.Println(err.Error())
		return
	}

	output, err := base64.StdEncoding.DecodeString(*resp.Output)

	re := regexp.MustCompile("BEGIN SSH HOST KEY FINGERPRINTS(.|\\n)*END SSH HOST KEY FINGERPRINTS")

	block := re.Find(output)
	if block == nil {
		log.Println("No fingerprints found")
		return
	}

	f := regexp.MustCompile("(:?[[:xdigit:]]{2}){16}")
	matches := f.FindAll(block, -1)

	for _, match := range matches {
		fmt.Println(string(match))
	}
}
