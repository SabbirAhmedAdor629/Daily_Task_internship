package main

import (
	// "fmt"
	// "io/ioutil"
	// "log"
	"encoding/json"
	"log"
	"net/http"
)

type Book struct {
	Title string `json:"title"`
	Author string `json:"author"`
	Pages int `json:"pages"`
}


func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	book := Book { Title: "The Gunslinger", Author: "Stephen King", Pages: 304 }

	json.NewEncoder(w).Encode(book)

}



func Hello(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "text/html")
	w.Write([]byte ("<h1 style='color: steelblue'> Hello </h1>"))
}

func main(){
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/book", GetBook)
	
	log.Fatal(http.ListenAndServe(":5002", nil))	
}







// func main(){
// 	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request){
// 		rw.Write([]byte("Hi Ador!"))
// 	})
// 	log.Fatal(http.ListenAndServe(":5002", nil))	
// }



		// Panic
// func main(){
// 	fmt.Println("start")
// 	panic("dont know what to do now")
// 	fmt.Println("end")
// }



		// Uses of defer
// func main(){
// 	res, err := http.Get("https://www.google.com/robots.txt")
// 	if (err != nil){
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	robots, err := ioutil.ReadAll(res.Body)
// 	if err != nil{
// 		log.Fatal(err)
// 	}
	
// 	fmt.Printf("%s",robots)
// }







		// just Practicing the code again
// func main(){
// 	res, err := http.Get("https://www.google.com/robots.txt")
// 	if(err != nil){
// 		log.Fatal(err)
// 	}

// 	robots, err := ioutil.ReadAll(res.Body)
// 	if(err != nil){
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%s",robots)
// }