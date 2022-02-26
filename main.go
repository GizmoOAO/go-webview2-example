package main

import (
	"go-webview2-example/internal/w32"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
)

func main() {
	dataPath, _ := filepath.Abs("./userdata")
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		DataPath:  dataPath,
		WindowOptions: webview2.WindowOptions{
			Title: "go-webview2 Example",
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()

	// update window icon
	w32.SendMessage(w.Window(), 0x0080, 1, w32.ExtractIcon(os.Args[0], 0))

	w.SetSize(800, 600, webview2.HintNone)

	chromium := getChromium(w)

	folderPath, _ := filepath.Abs("./public")
	webview := chromium.GetICoreWebView2_3()
	webview.SetVirtualHostNameToFolderMapping(
		"app.assets", folderPath,
		edge.COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY_CORS,
	)
	w.Navigate("http://app.assets/index.html")

	w.Run()
}

func getChromium(w webview2.WebView) *edge.Chromium {
	browser := reflect.ValueOf(w).Elem().FieldByName("browser")
	browser = reflect.NewAt(browser.Type(), unsafe.Pointer(browser.UnsafeAddr())).Elem()
	return browser.Interface().(*edge.Chromium)
}
