package common

import (
	"archive/zip"
	"fmt"
	"image/png"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/shinomontaz/pixel"
)

type Opener interface {
	Open(name string) (fs.File, error)
}

type Loader struct {
	opener Opener
	root   string
}

type Option func(*Loader)

func WithZip(name string) Option {
	return func(l *Loader) {
		l.opener, _ = zip.OpenReader(name)
	}
}

func NewLoader(root string, opts ...Option) *Loader {
	root = fmt.Sprintf("%s/", strings.TrimSuffix(root, "/")) // root must be with trailing slash /

	l := &Loader{
		root: root,
	}
	for _, opt := range opts {
		opt(l)
	}

	if l.opener == nil {
		l.opener = os.DirFS(root)
		l.root = ""
	}

	return l
}

func (l *Loader) Open(name string) (fs.File, error) {
	return l.opener.Open(fmt.Sprintf("%s%s", l.root, name))
}

func (l *Loader) LoadPicture(name string) (pixel.Picture, error) {
	file, err := l.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func (l *Loader) LoadFileToString(name string) (string, error) {
	file, err := l.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
