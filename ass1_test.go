package main

import (
	//////////////	"net/http"
	"testing"
)

/*Mocked data*/
/*var Mock = ResponsePayload{
	Name: "kafka",
	Owner: Owne{
		Login: "apache",
	},
	CommitInfo: Commit{
		Login:         "ijuma",
		Contributions: 306,
	},
	Languages: []string{"Batchfile", "HTML", "Java", "Python", "Scala", "Shell", "XSLT"},
}*/

/*func TestStringToArray(t *testing.T) {
	test := "!! foo, bar, 21321,    baz"
	expected := []string{"foo", "bar", "baz"}
	actual := toArray(test)
	if actual[0] != expected[0] && actual[1] != expected[1] && actual[2] != expected[2] {
		t.Error("Expected and actual does not match")
	}
}*/

func TestCheckLength(t *testing.T) {
	/*case valid*/
	test := []string{"localhost", "projectinfo", "v1", "github.com", "apache", "kafka"}
	expected := true
	actual := checkLength(test, 6)
	if actual != expected {
		t.Error("Check length did not respond correctly")
	}

	/*case fail*/
	expected = false
	actual = checkLength(test, 5)
	if actual != expected {
		t.Error("Check length did not respond correctly")
	}
}

/*TODO getURL*/
func TestGetURL(t *testing.T) {
	/*Should return a valid URL*/
	test := "localhost/projectinfo/v1/github.com/apache/kafka"
	expected := "https://api.github.com/repos/apache/kafka"
	actual, _ := getURL(test, "")
	if actual != expected {
		t.Error("GetURL did not function as expected, case 1")
	}

	/*Should return a valid URL*/
	expected = "https://api.github.com/repos/apache/kafka/contributors"
	actual, _ = getURL(test, "contributors")
	if actual != expected {
		t.Error("getURL did not function as expected, case 2")
	}

	/*Should return -1, incorrect URL*/
	test = "localhost/projectinfo/v1/github.com/apache"
	expected2 := -1
	_, err := getURL(test, "")
	if err != expected2 {
		t.Error("getURL did not function as expected, case 3")
	}

	/*Should return -1, incorrect URL*/
	_, err = getURL(test, "contributors")
	if err != expected2 {
		t.Error("getURL did not function as expected, case 4")
	}
}

/*TODO getData */
