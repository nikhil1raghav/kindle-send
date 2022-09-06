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

`kindle-send` is a command line utility to send files and webpages to your kindle via e-mail. 

Webpages are optimized for viewing on kindle


<p align = "center">
<figure>
<img width="90%" src="assets/toepub.png">
<figcaption>Credits - Netflix tech blog<a href="https://netflixtechblog.com/fixing-performance-regressions-before-they-happen-eab2602b86fe">Fixing Performance Regressions Before they Happen</a></figcaption>
</figure>
</p>


---


An epub is created from a url, then mailed to the kindle. Amazon converts that epub into azw3 for viewing on kindle.

So you can use kindle-send, even if you're using a different ereader like Kobo and Remarkable if it supports pushing ebooks via email.



---

### Installation

To run kindle-send you just need the compiled binary, no other dependency is required.

As this was not the case with the older [python version](https://github.com/nikhil1raghav/kindle-send/tree/python) which required percollate, calibre etc.






Download the binary for your operating system and architecture from [release page](https://github.com/nikhil1raghav/kindle-send/releases) and add it to your [PATH](https://en.wikipedia.org/wiki/PATH_(variable)).
If there is no binary compatible for your system. Please create an issue.


For the first time when you run `kindle-send`, you need to answer some questions to create a configuration file, which has options like sender, receiver, password and path to store the generated files.


If you're using gmail to send mails to kindle, please consider creating an [app password](https://support.google.com/mail/answer/185833?hl=en-GB) and then using it.


---



### Following modes of operation are supported

__1. Send a file__

Using `kindle-send` to mail an already existing file.

```sh
kindle-send --file <path-to-file>
```


<p align="center">
  <img width="100%" src="assets/file-send.svg">
</p>


__2. Send a webpage__

Quickly send a webpage to kindle


```sh
kindle-send  --url <link-to-a-webpage>
```

<p align="center">
  <img width="100%" src="assets/kindle-send-window.svg">
</p>


__3. Multiple webpages combined in a single volume__


Create a text file with new line separated links of webpages and then pass it as link file to `--linkfile` option


```sh
kindle-send --linkfile <path-to-url-file>
```

<p align="center">
  <img width="100%" src="assets/linkfile.svg">
</p>






### Additional options

Default timeout for mail is 2 minutes, if you get timeout error while sending bigger files. Please increase the timeout using `--mail-timeout <number of seconds>` option



Specify the title for the document using `--title` option.

Specify a different configuration file using `--config` option. Configuration is stored in home directory as `KindleConfig.json`. You can directly edit it if you want.

When sending a collection of pages if no title is provided, volume takes the title of the first page.


---

## Contribute

Feel free to create an issue and then working on some feature, so that we don't overwrite each other.


## Todo

- [ ] Weekly RSS feed dump, when combined with `cron`
- [ ] Better CSS & formatting for epub
- [ ] Compressing images before embedding to reduce final file size
- [x] Remove dependency on percollate and calibre
- [x] Make installation easier


