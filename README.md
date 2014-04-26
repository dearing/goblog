blog
====
**still in development**

blogging with go, redis and Oauth2

about
-----
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
	"ContentFolder": "content",
	"TemplateFolder": "templates",
	"Suffix": ".md",
	"WWWHost": ":9002",
	"RedisHost": "localhost:6379",
	"RedisPass": "",
	"RedisDB": -1,
	"ClientID": "abcdef",
	"ClientSecret": "abcdef0123456",
	"RedirectURL": "/callback",
	"AdminLogin": "some_github_username"
}
```
