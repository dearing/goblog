**still in development**

blogging with go, markdown and redis

about
-----
simple blogging; markdown for design; redis for caching and metrics

try it out
---
```
go get github.com/dearing/blog
go install

cd example
blog --help

Usage of blog:
  -conf="blog.conf": JSON configuration
  -generate=false: generate a new config as conf is set
  
```
```
{
  	"ContentFolder": "content",
	"TemplateFolder": "templates",
	"Suffix": ".md",
	"WWWHost": ":9000",
	"WWWRoot": "wwwroot",
	"RedisHost": "127.0.0.1:6379",
	"RedisPass": "",
	"RedisDB": -1,
	"Verbose": true
}
```
