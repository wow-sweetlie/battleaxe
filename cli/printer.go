package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
)

var (
	colorField = color.New(color.FgBlue, color.Bold).SprintFunc()
	colorValue = color.New().SprintFunc()
)

func printVersion() {
	fmt.Printf("v%s\n", version)
	os.Exit(0)
}

// PrintBody : print response body
func PrintBody(resp *http.Response, pretty bool) error {
	defer resp.Body.Close()
	if pretty {
		b := new(bytes.Buffer)
		b.ReadFrom(resp.Body)
		f := jsoncolor.NewFormatter()
		err := f.Format(os.Stdout, b.Bytes())
		if err != nil {
			return err
		}
		os.Exit(0)
	}
	_, err := io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// PrintHeader : print response header
func PrintHeader(resp *http.Response) {
	for field, value := range resp.Header {
		formatedValue := strings.Join(value, ",")
		fmt.Printf("%s: %s\n", colorField(field), colorValue(formatedValue))
	}
}