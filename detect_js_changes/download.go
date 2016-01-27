package detect_js_changes

import (
	"io"
	"net/http"
	"os"
)

func Download(url string, destination string) error {
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
