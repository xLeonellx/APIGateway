package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type API struct {
	Rout *mux.Router
}

type ResponseDetailed struct {
	NewsDetailed struct {
		ID      int    `json:"ID"`
		Title   string `json:"Title"`
		Content string `json:"Content"`
		PubTime int    `json:"PubTime"`
		Link    string `json:"Link"`
	} `json:"NewsDetailed"`
	Comments []struct {
		ID              int    `json:"ID"`
		NewsID          int    `json:"newsID"`
		ParentCommentID int    `json:"parentCommentID"`
		Content         string `json:"content"`
		PubTime         int    `json:"pubTime"`
	} `json:"Comments"`
}

const limit = 10 //limit posts on one page
const newsService = "http://news:8080"
const commentService = "http://comments:8081"
const censorAddr = "http://censorship:8082"

// Создание объекта api
func New() *API {
	api := API{
		Rout: mux.NewRouter(),
	}
	api.endpoints()
	return &api

}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.Rout.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rout.HandleFunc("/news/latest", api.newsLatestHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rout.HandleFunc("/news/search", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rout.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)

}

func MakeRequest(req *http.Request, method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Header
	return client.Do(r)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.Rout
}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	sParam := r.URL.Query().Get("s")

	resp, err := MakeRequest(r, http.MethodGet, newsService+"/news"+"?page="+pageParam+"&"+"s="+sParam, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)

}

func (api *API) newsLatestHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	resp, err := MakeRequest(r, http.MethodGet, newsService+"/news/latest"+"?page="+pageParam, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)

}

func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "search parameters are required", http.StatusBadRequest)
		return
	}
	chNews := make(chan *http.Response, 2)
	chComments := make(chan *http.Response, 2)
	chErr := make(chan error, 2)
	var response ResponseDetailed
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		respNews, err := MakeRequest(r, http.MethodGet, newsService+"/news/search"+"?id="+idParam, nil)
		chErr <- err
		chNews <- respNews

	}()

	go func() {
		defer wg.Done()

		respComments, err := MakeRequest(r, http.MethodGet, commentService+"/comments"+"?news_id="+idParam, nil)
		chErr <- err
		chComments <- respComments

	}()
	wg.Wait()
	close(chErr)

	for err := range chErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

block:
	for {
		select {
		case respNews := <-chNews:
			body, err := ioutil.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := ioutil.ReadAll(respComment.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.Comments)
		default:
			break block
		}

	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close() //  must close
	Body1 := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	Body2 := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	respCensor, err := MakeRequest(r, http.MethodPost, censorAddr+"/comments/add", Body1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "incorrect comment content", respCensor.StatusCode)
		return
	}

	resp, err := MakeRequest(r, http.MethodPost, commentService+"/comments/add", Body2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
