package tools

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"polar_reflow/logger"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const customTimeFormat = "2006-01-02T15:04:05"

func FormatTime(t time.Time) string {
	return t.Format(customTimeFormat)
}

func ParseTime(t string) time.Time {
	timePoint, err := time.Parse(customTimeFormat, t)
	logger.Error(err.Error())
	return timePoint
}

func Dumper(v any) {
	spew.Dump(v)
}

func UnpackArchive(src, dest string) error {
	switch {
	case strings.HasSuffix(src, ".zip"):
		return unzip(src, dest)
	case strings.HasSuffix(src, ".tar"):
		return untar(src, dest, "")
	case strings.HasSuffix(src, ".tar.gz"), strings.HasSuffix(src, ".tgz"):
		return untar(src, dest, "gzip")
	case strings.HasSuffix(src, ".tar.bz2"):
		return untar(src, dest, "bzip")
	default:
		return fmt.Errorf("unsupported file format")
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fpath := filepath.Join(dest, file.Name)

		if !filepath.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func untar(src, dest, decompressor string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var r io.Reader = file

	switch decompressor {
	case "gzip":
		r, _ = gzip.NewReader(file)
	case "bzip":
		r = bzip2.NewReader(file)
	}

	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fpath := filepath.Join(dest, header.Name)

		if !filepath.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			return fmt.Errorf("unsupported type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

func OpenFile(path string) ([]byte, error) {
	reader, err := os.Open(path)
	if err != nil {
		logger.Error(err.Error())
	}
	return io.ReadAll(reader)
}
