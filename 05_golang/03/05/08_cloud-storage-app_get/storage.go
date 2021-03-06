package storage

import (
	"io"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

const bucketName = "serious-water-88716.appspot.com"

func init() {
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/get", handleGet)

}

func getCloudContext(appengineContext context.Context) context.Context {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(appengineContext, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: appengineContext},
		},
	}

	return cloud.NewContext(appengine.AppID(appengineContext), hc)
}

func handlePut(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: ctx},
		},
	}

	cctx := cloud.NewContext(appengine.AppID(ctx), hc)

	writer := storage.NewWriter(cctx, bucketName, "example2.txt")
	io.WriteString(writer, "Another file")
	err := writer.Close()
	if err != nil {
		http.Error(res, "ERROR WRITING TO BUCKET: "+err.Error(), 500)
		return
	}
}


func handleGet(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	cctx := getCloudContext(ctx)

	rdr, err := storage.NewReader(cctx, bucketName, "example.txt")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	defer rdr.Close()

	io.Copy(res, rdr)
}
