
# DJ2DC

This projects permits importing the CSV file generated from a Deejay command into one's Discogs collection.

### Build

```
git clone git@github.com:etiennejournet/dj2dc.git
cd dj2dc && go install
```

### Usage
Find your binary in `$HOME/go/bin/dj2dc` or `$GOPATH/bin/dj2dc`

```
Usage of dj2dc:
  -file string
    	path of your csv file (default "csv")
  -token string
    	token for discogs api
```
