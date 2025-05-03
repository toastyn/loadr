### Currently in development, API is inconsistent

loadr is a library which extends the functionality of the standard html/template functionality. The library drew inspiration from [this article](https://philipptanlak.com/web-frontends-in-go/) and essentially solidifies the idea into a re-usable pattern across projects.

Major concepts are:
1. loadr allows HTML pages to be live reloaded or cached upfront 
2. loadr has a fail-fast approach, where non-existent html files or fragments error out on application startup instead of during application runtime
3. HTML files and fragments behave like native go code, allowing the developer to choose how the parsing, data injection and caching is done


# install

```
go get github.com/nesbyte/loadr
```

# examples
See [_examples](_examples) for more complete and involved examples
