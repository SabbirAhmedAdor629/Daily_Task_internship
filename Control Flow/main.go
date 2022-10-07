package main

import (
	// "fmt"
	// "io/ioutil"
	// "log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello Ador!"))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}



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