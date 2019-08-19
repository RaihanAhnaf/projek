package tagginghtml

import (
	"bytes"
	"os"
	"runtime"

	"os/exec"

	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func TaggingHtml(sourcepath string, destpath string, filename string) {
	switch runtime.GOOS {
	case "windows":
		pyfilepath := wd + "consoleapps/python/htmltools.py"
		htmlsource := sourcepath + "/" + filename
		htmlresult := destpath + "/" + filename

		cmdstr := []string{"/C", "python", pyfilepath, "-i", htmlsource, "-o", htmlresult}
		cmd := exec.Command("cmd", cmdstr...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			tk.Println("Failed: " + stderr.String())
		}

	case "linux":
		pyfilepath := wd + "consoleapps/python/htmltools.py"
		htmlsource := sourcepath + "/" + filename
		htmlresult := destpath + "/" + filename

		cmdstr := []string{pyfilepath, "-i", htmlsource, "-o", htmlresult}
		cmd := exec.Command("python3", cmdstr...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			tk.Println("Failed: " + stderr.String())
		}

	case "darwin":
		pyfilepath := wd + "consoleapps/python/htmltools.py"
		htmlsource := sourcepath + "/" + filename
		htmlresult := destpath + "/" + filename

		cmdstr := []string{pyfilepath, "-i", htmlsource, "-o", htmlresult}
		cmd := exec.Command("python3", cmdstr...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			tk.Println("Failed: " + stderr.String())
		}
	}

}
