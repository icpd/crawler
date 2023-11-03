package worker

import (
	"errors"

	"fmt"
	"log"

	"imooc.com/ccmouse/learngo/crawler/config"
	"imooc.com/ccmouse/learngo/crawler/engine"
	xcar "imooc.com/ccmouse/learngo/crawler/xcar/parser"
	zhenai "imooc.com/ccmouse/learngo/crawler/zhenai/parser"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

func SerializeResult(
	r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		result.Requests = append(result.Requests,
			SerializeRequest(req))
	}
	return result
}

func DeserializeRequest(
	r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

func DeserializeResult(
	r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error deserializing "+
				"request: %v", err)
			continue
		}
		result.Requests = append(result.Requests,
			engineReq)
	}
	return result
}

func deserializeParser(
	p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		return engine.NewFuncParser(
			zhenai.ParseCityList,
			config.ParseCityList), nil
	case config.ParseCity:
		return engine.NewFuncParser(
			zhenai.ParseCity,
			config.ParseCity), nil

	case config.ParseProfile:
		if userName, ok := p.Args.(string); ok {
			return zhenai.NewProfileParser(
				userName), nil
		} else {
			return nil, fmt.Errorf("invalid "+
				"arg: %v", p.Args)
		}
	case config.ParseCarDetail:
		return engine.NewFuncParser(
			xcar.ParseCarDetail,
			config.ParseCarDetail), nil
	case config.ParseCarModel:
		return engine.NewFuncParser(
			xcar.ParseCarModel,
			config.ParseCarModel), nil
	case config.ParseCarList:
		return engine.NewFuncParser(
			xcar.ParseCarList,
			config.ParseCarList), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	default:
		return nil, errors.New(
			"unknown parser name")
	}
}
