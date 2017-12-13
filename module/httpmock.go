// Httpmock can easy mocking of http responses
// from external resources.
package main

import (
	"net/http"
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
)

//Mock a simple HTTP request to prove its ability
//if it success ,you will get the new response string "I mocked it"
func simpleMock() {

	//Send the HTTP request to the target URL
	//It will print true response without mock
	resp, err := http.Get("http://www.baidu.com")
	fmt.Println(resp, err)

	//Begin mock
	//It will intercept all http request
	//before execute httpmock.DeactivateAndReset
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//Register the url and method to be mocked
	httpmock.RegisterResponder("GET", "http://www.baidu.com",
		httpmock.NewStringResponder(200, "I mocked it"))

	//Verification mock effect by request the url again
	mockResp, mockErr := http.Get("http://www.baidu.com")
	body, _ := ioutil.ReadAll(mockResp.Body)
	fmt.Println(string(body), mockErr)
}

func advancedMock() {

	//Begin mock
	//It will intercept all http request
	//before execute httpmock.DeactivateAndReset
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//Register the url and method to be mocked
	httpmock.RegisterResponder("GET", "http://www.baidu.com",
		//You can set different mock response
		//in different request
		//This example only filter request that not null
		func(req *http.Request) (*http.Response, error) {
			if req != nil {
				resp := httpmock.NewStringResponse(200, "I mocked it")
				return resp, nil
			}
			return nil, nil
		}, )

	//Verification mock effect by request the url
	mockResp, mockErr := http.Get("http://www.baidu.com")
	body, _ := ioutil.ReadAll(mockResp.Body)
	fmt.Println(string(body), mockErr)
}

func main() {
	//Mock a simple HTTP request to "http://www.baidu.com" to prove its ability
	//if it success ,you will get two responses :
	// before mock true response without mock,
	// after mock  mockedResponse string "I mocked it"
	fmt.Println("---------the simpleMock")
	simpleMock()

	//Mock a advanced HTTP request
	//if it success ,you will get a response is mock response
	//by your send different requests
	fmt.Println("---------the advanceMock")
	advancedMock()
}
