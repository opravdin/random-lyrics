package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	mixer2 "taylor-swift-bot/internal/mixer"
	"taylor-swift-bot/internal/parser"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	songs := parser.ReadSongs()
	mixer := mixer2.New()

	wg := sync.WaitGroup{}
	for _, song := range songs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mixer.ProvideLines(song.Lines)
		}()

	}
	wg.Wait()
	err := mixer.ValidateTags()
	if err != nil {
		log.Println(mixer.TaggedLines)
		panic(err)
	}
	log.Print("mixer ready")

	http.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		recipe := r.URL.Query().Get("recipe")
		var res = ""
		var err error
		if recipe == "" {
			res, err = mixer.Mix()
		} else {
			res, err = mixer.Make(recipe)
		}
		if err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}
		fmt.Fprint(w, res)
	})

	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", mixer.RecipeBook.RecipeNames)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
