package detect_js_changes

import (
	"fmt"
	"os"
	"path"
)

func Reset(dir string) {
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		os.Remove(path.Join(dir, file))
	}
}
