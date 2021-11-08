package main

import (
	"encoding/json"
	"flag"
	"os"
	"regexp"

	crawler "github.com/odia/juscaba/crawler"
	extracttext "github.com/odia/juscaba/extracttext"
	fetcher "github.com/odia/juscaba/fetcher"
	shared "github.com/odia/juscaba/shared"
	log "github.com/sirupsen/logrus"
)

type arguments struct {
	jsonPath    string
	fm          *shared.FileManager
	exp         *shared.Expediente
	parseImages bool
}

func parseArguments() (*arguments, error) {
	var pdfsPath, expId string
	var err error
	args := arguments{}
	flag.StringVar(&args.jsonPath, "json", "", "json destination path")
	flag.StringVar(&pdfsPath, "pdfs", "", "pdfs destination path")
	flag.StringVar(&expId, "expediente", "", "expediente identifier (e.g.: \"182908/2020-0\")")
	flag.BoolVar(&args.parseImages, "images", true, "apply ocr")
	flag.Parse()

	log.WithFields(log.Fields{
		"json":        args.jsonPath,
		"pdfs":        pdfsPath,
		"expediente":  expId,
		"parseImages": args.parseImages,
	}).Print("arguments")

	args.fm = &shared.FileManager{Directory: pdfsPath}
	args.exp, err = crawler.GetExpediente(expId)
	if err != nil {
		return nil, err
	}

	return &args, nil
}

func main() {
	args, err := parseArguments()
	if err != nil {
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"expediente":  args.exp,
		"actuaciones": len(args.exp.Actuaciones),
	}).Printf("finished")

	blacklist := os.Getenv("BLACKLIST_REGEX")
	if blacklist != "" {
		_, err := regexp.Match(blacklist, []byte{})
		if err != nil {
			log.WithFields(log.Fields{
				"regex": blacklist,
				"error": err.Error(),
			}).Error("failed to parse blacklist regex")
			os.Exit(1)
		}
	}

	for _, act := range args.exp.Actuaciones {
		for _, doc := range act.Documentos {
			if blacklist != "" {
				if match, _ := regexp.Match(blacklist, []byte(doc.URL)); match {
					log.WithFields(log.Fields{
						"url": doc.URL,
					}).Info("skipping blacklisted URL")
					continue
				}
			}
			fetcher.Download(args.fm, doc.URL)
			reader, err := args.fm.GetReader(doc.URL)
			if err != nil {
				continue
			}
			doc.Content, _ = extracttext.GetDocumentText(reader, args.parseImages)
			doc.MirrorURL, _ = args.fm.DestinationURLforSourceURL(doc.URL)
		}
	}
	fp, err := os.Create(args.jsonPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		})
		os.Exit(2)
	}
	defer fp.Close()
	json.NewEncoder(fp).Encode(args.exp)
}
