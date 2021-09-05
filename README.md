# Fastcrawler
A Fast HTTP crawler written in Go

# Install 
```
git clone https://github.com/vmnguyen/fastcrawler

cd fastcrawler

go build

```


# Usage
```
Usage of ./fastcrawl:
  -c int
        Concurency number (default 50)
  -t string
        Target to scan (default "https://example.com")
```

# Example

```./fastcrawl -t https://example.com -c 50```

