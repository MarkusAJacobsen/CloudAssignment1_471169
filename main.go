package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	//"unicode"
)

//ResponsePayload represents root in github api
type ResponsePayload struct {
	Name       string
	Owner      Owne
	CommitInfo Commit
	Languages  []string
}

//Owne object from github
type Owne struct {
	Login string
}

//Commit object from github
type Commit struct {
	Login         string
	Contributions int
}

//Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go\n")
}

func infoPage(w http.ResponseWriter, r *http.Request) {
	incomplete := false
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
	lang := new(map[string]interface{})

	/*Fetch data*/
	if error := getData(requestProjectURL, generalInfo); error != nil {
		printError(w, error)
		incomplete = true
	}

	if error := getData(requestContributorsURL, contribution); error != nil {
		printError(w, error)
		incomplete = true
	}

	//lang := getLang(w, requestLanguagesURL)
	if error := getData(requestLanguagesURL, lang); error != nil {
		printError(w, error)
		incomplete = true
	}

	/*Append contribution struct data to generalInfo Struct
	 * If no commits, give verbose message instead of ""
	 */
	if contribution[0].Login != "" {
		generalInfo.CommitInfo.Login = contribution[0].Login
		generalInfo.CommitInfo.Contributions = contribution[0].Contributions
	} else {
		generalInfo.CommitInfo.Login = "No commits registered"
		generalInfo.CommitInfo.Contributions = 0
	}

	//toArray(lang)
	for r := range *lang {
		generalInfo.Languages = append(generalInfo.Languages, r)
	}

	/*Encode struct and print it on screen*/
	if !incomplete {
		status := 404
		http.Error(w, http.StatusText(status), status)
		return
	}
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
	}
	return (base + parts[4] + "/" + parts[5]), 0
}

func checkLength(s []string, length int) bool {
	if len(s) != length {
		return false
	}
	return true
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
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, "%s\n", b)
}

/*getData takes in a url and an interface, decodes the response and puts it in the interface*/
func getData(url string, payload interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return err
	}
	return err
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/", startPage)
	http.HandleFunc("/projectinfo/v1/", infoPage)
	panic(http.ListenAndServe(":"+port, nil))
	//panic(http.ListenAndServe(":8080", nil))
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
