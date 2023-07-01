package app

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

var log zap.SugaredLogger

func Run() error {
	{
		l, err := zap.NewProductionConfig().Build()
		if err != nil {
			return err
		}
		log = *l.Sugar()
	}
	defer log.Sync()

	// Define flags
	dirPtr := flag.String("path", ".", "the file directory")
	outPtr := flag.String("out", "out.json", "the output file")
	flag.Parse()
	dirPath := *dirPtr
	outPath := *outPtr
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed reading directory: %s", err)
	}
	var g gofeed.Feed
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		filePath := path.Join(dirPath, v.Name())
		filePtr, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed reading file: %s", err)
		}
		defer filePtr.Close()
		feed, err := gofeed.NewParser().Parse(filePtr)
		if errors.Is(err, gofeed.ErrFeedTypeNotDetected) {
			log.Errorf("failed parsing feed %s", filePath)
			continue
		}
		if err != nil {
			return fmt.Errorf("failed parsing feed %s: %s", filePath, err)
		}
		log.Debugw("read file", "items", feed.Items)
		g.Items = append(g.Items, feed.Items...)
	}

	filePtr, err := os.Create(outPath)
	defer filePtr.Close()
	if err != nil {
		return fmt.Errorf("failed reading out: %s", err)
	}
	enc := json.NewEncoder(filePtr)
	if err := enc.Encode(g); err != nil {
		return fmt.Errorf("failed writing out: %s", err)
	}

	return nil
}
