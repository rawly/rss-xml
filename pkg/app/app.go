package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

var log zap.SugaredLogger

// Define flags
var dirPath, outPath = Flags()

func Run() error {
	{
		zaplogger, err := zap.NewProductionConfig().Build()
		if err != nil {
			return err
		}
		log = *zaplogger.Sugar()
		defer log.Sync()
	}

	outPtr, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed create file: %w", err)
	}
	defer outPtr.Close()

	feedfiles, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed reading directory: %w", err)
	}

	var g gofeed.Feed
	for _, filename := range feedfiles {
		if filename.IsDir() {
			continue
		}

		file, err := os.Open(path.Join(dirPath, filename.Name()))
		if err != nil {
			return fmt.Errorf("failed reading file: %w", err)
		}
		defer file.Close()

		feed, err := gofeed.NewParser().Parse(file)
		if err != nil {
			log.Errorf("failed parsing file: %w", err)
			continue
		}

		g.Items = append(g.Items, feed.Items...)
	}

	enc := json.NewEncoder(outPtr)
	if err := enc.Encode(g); err != nil {
		return fmt.Errorf("failed encoding: %w", err)
	}

	return nil
}

func Flags() (dir, out string) {
	dirPtr := flag.String("path", ".", "the file directory")
	outPtr := flag.String("out", "out.json", "the output file")
	flag.Parse()
	return *dirPtr, *outPtr
}
