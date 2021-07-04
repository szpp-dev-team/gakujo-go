package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) Home() error {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/home/home/initialize"
	params := make(url.Values)
	params.Set("EXCLUDE_SET", "")
	resp, err := c.client.PostForm(reqURL, params)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))

	return nil
}
