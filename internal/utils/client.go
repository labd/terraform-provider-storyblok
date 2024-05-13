package utils

import "github.com/labd/storyblok-go-sdk/sbmgmt"

func GetClient(data any) sbmgmt.ClientWithResponsesInterface {
	c, ok := data.(sbmgmt.ClientWithResponsesInterface)
	if !ok {
		panic("invalid client type")
	}
	return c
}
