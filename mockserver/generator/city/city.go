// Package city implements city generator.
package city

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"imooc.com/ccmouse/learngo/mockserver/config"
	"imooc.com/ccmouse/learngo/mockserver/generator/profile"
)

// Generator represents the city generator.
type Generator struct {
	Tmpl       *template.Template
	ProfileGen *profile.Generator
}

const (
	itemCount = 20
	pageCount = 5
)

type params struct {
	City string `uri:"city" binding:"required"`
	Page int    `uri:"page"`
}

// HandleRequest is the gin request handler for city generation.
func (g *Generator) HandleRequest(c *gin.Context) {
	var p params
	err := c.BindUri(&p)
	if err != nil {
		log.Printf("BindUri(): %v.", err)
		return
	}
	if p.Page == 0 {
		p.Page = 1
	}
	g.handleRequest(c, p)
}

func (g *Generator) handleRequest(c *gin.Context, p params) {
	err := g.generate(p, c.Writer)

	if err != nil {
		log.Printf("Cannot generate page for city %q and page %d: %v.", p.City, p.Page, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (g *Generator) generate(p params, w io.Writer) error {
	h, err := hashCode(p)
	if err != nil {
		return fmt.Errorf("cannot calculate hash code: %v", err)
	}
	r := rand.New(rand.NewSource(h))

	items := make([]ContentItem, itemCount)
	for i := range items {
		id := r.Int63()
		p := g.ProfileGen.GenerateProfile(id)
		items[i] = ContentItem{
			ID:      id,
			Profile: p,
			URL:     fmt.Sprintf("http://%s/mock/album.zhenai.com/u/%d", config.ServerAddress, id),
		}
	}

	var pages []PageItem
	for i := 0; i < pageCount; i++ {
		targetPage := p.Page - 1 + i
		if targetPage < 1 {
			continue
		}
		url := ""
		if targetPage != p.Page {
			url = fmt.Sprintf("http://%s/mock/www.zhenai.com/zhenghun/%s/%d", config.ServerAddress, p.City, targetPage)
		}
		pages = append(pages, PageItem{
			URL:  url,
			Page: targetPage,
		})
	}

	return g.Tmpl.Execute(w, Content{
		Items: items,
		Pages: pages,
	})
}

func hashCode(v interface{}) (int64, error) {
	h := fnv.New64()
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)
	if err != nil {
		return 0, fmt.Errorf("cannot encode gob for param: %v", err)
	}
	_, err = h.Write(b.Bytes())
	if err != nil {
		return 0, fmt.Errorf("cannot write to hash: %v", err)
	}
	return int64(h.Sum64()), nil
}

// ContentItem defines content item bound into html template.
type ContentItem struct {
	ID      int64
	Profile *profile.PhotoProfile
	URL     string
}

// PageItem defines page item bound into html template.
type PageItem struct {
	URL  string
	Page int
}

// Content defines the master content bound into html template.
type Content struct {
	Items []ContentItem
	Pages []PageItem
}
