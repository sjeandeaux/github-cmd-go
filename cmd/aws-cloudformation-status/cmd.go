package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//commandLine the arguments command line
type commandLine struct {
	profile string
	name    string
	stdout  io.Writer
	stderr  io.Writer
	stdin   *os.File
}

func (c *commandLine) init() {

	//flag
	log.SetPrefix("[aws-cloudformation-status]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.name, "stack-name", "", "")
	flag.StringVar(&c.profile, "profile", "", "")

	flag.Parse()

}

func (c *commandLine) main() int {
	opt := session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		Profile: c.profile,
	}

	sess := session.Must(session.NewSessionWithOptions(opt))

	cldf := cloudformation.New(sess)

	input := &cloudformation.DescribeStacksInput{}
	input.SetStackName(c.name)

	output, err := cldf.DescribeStacks(input)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}

	if output != nil {
		for _, stack := range output.Stacks {
			fmt.Fprintf(c.stdout, "%q %q \n", stack.StackName, stack.StackStatus)
		}
	}

	return 0
}
