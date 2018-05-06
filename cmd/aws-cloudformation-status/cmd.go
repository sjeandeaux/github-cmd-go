package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

//commandLine the arguments command line
type commandLine struct {
	internalcmd.CommandLine

	profile string
	name    string
}

func (c *commandLine) init() {

	//flag
	c.Init("[aws-cloudformation-status]")

	flag.StringVar(&c.name, "name", "", "")
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
		return c.Fatal(err)
	}

	if output != nil {
		for _, stack := range output.Stacks {
			fmt.Fprintln(c.Stdout, stack.StackName, " ", stack.StackStatus)
		}
	}

	return 0
}
