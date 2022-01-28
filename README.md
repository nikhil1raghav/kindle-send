<p align = "center">
<img src="assets/kindle-send-small.png" width="40%">
</p>

<p align = "center">
<strong>Send blogs, documents, collection of webpages to your kindle</strong>
</p>
<h3 align="center">
<a href="#contribute">Contribute</a>
<span> · </span>
<a href="#documentation">Documentation</a>
<span> · </span>
<a href="#todo">Todo</a>
</h3>

---

## Documentation

`kindle-send` is a command line utility to send files, webpages after converting them to mobi format to your kindle via e-mail. 
Webpages are optimized for viewing on kindle, thanks to [percollate](https://github.com/danburzo/percollate).



<p align="center">
  <img width="100%" src="assets/kindle-send-window.svg">
</p>




### How to use it?

1. Clone this repo
```sh
git clone "https://github.com/nikhil1raghav/kindle-send.git"
```

2. Install percollate

```sh
npm install -g percollate
```

3. Install Calibre

I don't know if `ebook-convert` can be installed as a stand-alone  executable. If you have a kindle then there is a very high chance that you have calibre installed.

```sh
sudo pacman -S calibre
```

4. Install other dependencies

There are only two extra dependencies `selectolax` for crawling and `argparser` for parsing arguments.

```sh
pip install argparser selectolax
```

5. Configure the Email

There is a configuration file `config.py` which stores your email credentials and other global options. Bare minimum you need to do is fill in the credentials of the mail that you're going to use for sending the documents.



### Following modes of operation are supported

__1. Send a file__


```sh
kindle-send --file <path-to-file>
```

__2. Send a webpage__
```sh
kindle-send  --link <link-to-a-webpage>
```

__3. Multiple webpages combined in a single volume__

Create a text file with new line separated links of webpages and then pass it as url file to `--link-file` option

```sh
kindle-send --link-file <path-to-url-file>
```

### Additional options

Specify the title for the document using `--title` option.
A different receiver using `--receiver` option

When sending a collection of pages if no title is provided, volume takes the title of the first page.






---

## Contribute

Currently it is a wrapper on [percollate](https://github.com/danburzo/percollate) and [Calibre's](https://github.com/danburzo/percollate) __ebook-convert__. 
Feel free to create an issue and then working on some feature, so that we don't overwrite each other.


---


## Todo


- [ ] Weekly RSS feed dump, when combined with `cron`
- [ ] Capability to create mobi without using `ebook-convert`
- [ ] `--convert` option to specify subject of email as `convert` so that documents are converted by amazon to supported formats before sending to kindle.


