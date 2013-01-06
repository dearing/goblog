blog
====
**still in development**

blogging with go, markdown, redis and Oauth2

about
-----
- markdown posts
- redis for caching and metrics
- github oauth2 for administration
- more fun with go

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

example config
----
```
{
	"Verbose": true,
	"EnableWWW": false,
	"ContentFolder": "content",
	"TemplateFolder": "templates",
	"Suffix": ".md",
	"WWWHost": ":9002",
	"WWWRoot": "example",
	"RedisHost": "localhost:6379",
	"RedisPass": "",
	"RedisDB": -1,
	"ClientID": "",
	"ClientSecret": "",
	"RedirectURL": "/callback",
	"AdminLogin": ""
}
```
