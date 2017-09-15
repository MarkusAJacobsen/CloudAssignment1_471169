package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"
)

//ResponsePayload represents root in github api
type ResponsePayload struct {
	Name       string   `json: "name"`
	Owner      Owne     `json: "owner"`
	CommitInfo Commit   `json: "contributors_url"`
	Languages  []string `json: "language"`
}

//Owne object from github
type Owne struct {
	Login string `json: "login"`
}

//Commit object from github
type Commit struct {
	Login         string `json: "login"`
	Contributions int    `json: "contributions"`
}

//Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go\n")
}

func infoPage(w http.ResponseWriter, r *http.Request) {

	/*Configure request URLs*/
	var requestProjectURL, requestContributorsURL, requestLanguagesURL string
	var err0, err1, err2 int
	requestProjectURL, err0 = getURL(r.URL.Path, "")
	requestContributorsURL, err1 = getURL(r.URL.Path, "contributors")
	requestLanguagesURL, err2 = getURL(r.URL.Path, "languages")

	if err0 != 0 && err1 != 0 && err2 != 0 {
		status := 400
		http.Error(w, http.StatusText(status), status)
		return
	}

	/*Set up data types*/
	contribution := new([1]Commit)
	generalInfo := &ResponsePayload{}

	/*Fetch data*/
	if error := getData(requestContributorsURL, contribution); error != nil {
		printError(w, error)
	}
	if error := getData(requestProjectURL, generalInfo); error != nil {
		printError(w, error)
	}

	lang := getLang(w, requestLanguagesURL)

	/*Append contribution structs data to generalInfo Struct*/
	generalInfo.CommitInfo.Login = contribution[0].Login
	generalInfo.CommitInfo.Contributions = contribution[0].Contributions
	generalInfo.Languages = toArray(lang)

	/*Encode struct and print it on screen*/
	returnResponse(w, generalInfo)
}

func getURL(u string, postfix string) (string, int) {
	//Needs to build a get request out of the url
	base := "https://api.github.com/repos/"
	parts := strings.Split(u, "/")

	/*If a supplied request URL is not on the form:
	 *"/foo/bar/baz/foz/kaa" Then its not valid
	 */
	if !checkLength(parts, 6) {
		return "", -1
	}
	if postfix != "" {
		return (base + parts[4] + "/" + parts[5] + "/" + postfix), 0
	} else {
		return (base + parts[4] + "/" + parts[5]), 0
	}
}

func checkLength(s []string, length int) bool {
	if len(s) != length {
		return false
	}
	return true
}

/*Strip away any non letter characters from a string, and return array*/
func toArray(s string) []string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	return strings.FieldsFunc(s, f)
}

/*printError prints an error message on the screen*/
func printError(w http.ResponseWriter, e error) {
	fmt.Fprintf(w, "Something went wrong %s\n", e)
}

/*returnResponse encodes a interface and writes it on screen*/
func returnResponse(w http.ResponseWriter, r interface{}) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		printError(w, err)
	}
	fmt.Fprintf(w, "%s\n", b)
}

func getLang(w http.ResponseWriter, url string) string {
	response, err := http.Get(url)
	if err != nil {
		printError(w, err)
	}
	defer response.Body.Close()
	lang, err := ioutil.ReadAll(response.Body)
	if err != nil {
		printError(w, err)
	}

	/*Following code snippet courtesy:
	 * https://stackoverflow.com/questions/40429296/converting-string-to-json-or-struct-in-golang
	 */
	in := []byte(lang)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)
	out, _ := json.Marshal(raw)
	return string(out)
}

/*getData takes in a url and an interface, decodes the response and puts it in the interface*/
func getData(url string, payload interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/", startPage)
	http.HandleFunc("/projectinfo/v1/", infoPage)
	panic(http.ListenAndServe(":8080", nil))
}

/* Reference material:
https://golang.org/pkg/net/http/
https://gist.github.com/Tinker-S/52ae0f981d7b86e0b34f
https://golang.org/doc/articles/wiki/#tmp_4
https://golang.org/pkg/encoding/json/
https://blog.golang.org/json-and-go
https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang
https://stackoverflow.com/questions/20866817/golang-decoding-json-into-custom-structure
https://sosedoff.com/2016/07/16/golang-struct-tags.html
*/
