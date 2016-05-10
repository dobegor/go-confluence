package confluence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ContentUpdate struct {
	Body struct {
		Storage struct {
			Representation string `json:"representation"`
			Value          string `json:"value"`
		} `json:"storage"`
	} `json:"body"`
	Id      string `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
}

type Content struct {
	Expandable struct {
		Ancestors   string `json:"ancestors"`
		Children    string `json:"children"`
		Descendants string `json:"descendants"`
		History     string `json:"history"`
		Metadata    string `json:"metadata"`
	} `json:"_expandable"`
	Links struct {
		Base       string `json:"base"`
		Collection string `json:"collection"`
		Self       string `json:"self"`
		Tinyui     string `json:"tinyui"`
		Webui      string `json:"webui"`
	} `json:"_links"`
	Body struct {
		Expandable struct {
			Editor     string `json:"editor"`
			ExportView string `json:"export_view"`
			Storage    string `json:"storage"`
		} `json:"_expandable"`
		View struct {
			Expandable struct {
				Content string `json:"content"`
			} `json:"_expandable"`
			Representation string `json:"representation"`
			Value          string `json:"value"`
		} `json:"view"`
	} `json:"body"`
	Container struct {
		Expandable struct {
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
			Icon        string `json:"icon"`
		} `json:"_expandable"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
		ID   int    `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"container"`
	Id    string `json:"id"`
	Space struct {
		Expandable struct {
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
			Icon        string `json:"icon"`
		} `json:"_expandable"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
		ID   int    `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"space"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Version struct {
		By struct {
			DisplayName    string `json:"displayName"`
			ProfilePicture struct {
				Height    int    `json:"height"`
				IsDefault bool   `json:"isDefault"`
				Path      string `json:"path"`
				Width     int    `json:"width"`
			} `json:"profilePicture"`
			Type     string `json:"type"`
			Username string `json:"username"`
		} `json:"by"`
		Message   string `json:"message"`
		MinorEdit bool   `json:"minorEdit"`
		Number    int    `json:"number"`
		When      string `json:"when"`
	} `json:"version"`
}

func (w *Wiki) contentEndpoint(contentID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + contentID)
}

func (w *Wiki) DeleteContent(contentID string) error {
	contentEndPoint, err := w.contentEndpoint(contentID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", contentEndPoint.String(), nil)
	if err != nil {
		return err
	}

	_, err = w.sendRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wiki) GetContent(contentID string, expand []string) (*Content, error) {
	contentEndPoint, err := w.contentEndpoint(contentID)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("expand", strings.Join(expand, ","))
	contentEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", contentEndPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}
	dataRes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(dataRes))

	var content Content
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (w *Wiki) UpdateContent(content *ContentUpdate) (*ContentUpdate, error) {
	jsonbody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	contentEndPoint, err := w.contentEndpoint(content.Id)
	req, err := http.NewRequest("PUT", contentEndPoint.String(), strings.NewReader(string(jsonbody)))
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(string(jsonbody))

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var newContent ContentUpdate
	err = json.Unmarshal(res, &newContent)
	if err != nil {
		return nil, err
	}

	return &newContent, nil
}
