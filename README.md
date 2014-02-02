Nord
====

Download Nord binaries for your platform [here](https://github.com/FZambia/nord/releases/latest) 

Social network share, like, tweet, pin statistics for URL address, written in Go. Server originally developed for bluffy.net (not yet released)

To get information about URL all you need is send HTTP GET request on running Nord:

```bash
curl "http://localhost:3000/?url=http://github.com/&providers=googleplus,facebook"
```

In example above we asked for data about `http://github.com/` url address from Google Plus and Facebook. The response is in JSON format:

```javascript
{
    "facebook": {
        "data": {
            "click_count": 1, 
            "comment_count": 2273, 
            "commentsbox_count": 0, 
            "like_count": 2006, 
            "share_count": 6705, 
            "total_count": 10984
        }, 
        "error": null
    }, 
    "googleplus": {
        "data": {
            "count": 4664
        }, 
        "error": null
    }
}
```


At moment the following social networks/providers are supported:

* Facebook (`?providers=facebook&url=http://mail.ru/`)
* Google Plus (`?providers=googleplus,facebook&url=http://mail.ru/`)
* Twitter (the same for `twitter`)
* Pinterest (`pinterest`)
* Delicious (`delicious`)
* LinkedIn (`linkedin`)
* VK (`vk`)


Features
--------

* dead-simple API
* JSONP support
* optional response caching with Redis
* can be integrated with your existing Go web applications


Install and run
---------------

```bash
git clone https://github.com/FZambia/nord nord/
cd nord/
go build
./nord
```

Or just download Nord single precompiled binary file for your platform [here](https://github.com/FZambia/nord/releases/latest) 


API
---

Available GET parameters:

* `providers` - comma separated list of social networks you want to be included into response
* `url` - the URL address of the page of your interest
* `timeout` - maximum request timeout in milliseconds, default 3000 (3 seconds)
* `callback` - wrap response in JSONP callback


From web browser
----------------

Nord supports JSONP, here is an example of jQuery jsonp request from web browser:

```javascript
$(function() {
    $.ajax({
        type: 'GET',
        url: "http://localhost:3000",
        data: {
            'url': location.href,
            'providers': 'googleplus,facebook,twitter,pinterest'
        },
        dataType: 'jsonp',
        success: function(data) {
           console.log(data);
        },
        error: function(e) {
           console.log(e.message);
        }
    });    
})
```

Configuration
-------------

You can change some defaults using command-line options:

* ``--timeout`` - default request timeout in milliseconds (3000 by default)
* ``--address`` - interface to bind to (e.g. 127.0.0.1, default 0.0.0.0)
* ``--port`` - port to listen on (default 3000)
* ``--prefix`` - url path prefix for main handler (default "", i.e. service will be available on "/" path)

```bash
./nord --address=127.0.0.1 --port=3000 --timeout=5000
```

Using command above we started Nord on http://localhost:3000/ . After this you can open your browser and go to http://localhost:3000/?providers=twitter&url=http://my.com

There are also options for configuring caching. See them below.

Caching
-------

Caching responses supported via Redis. To enable response caching with default settings use command-line option ``--cache``

```bash
./nord --cache
```

Data cached for each provider/url pair.

By default cache timeout is 60 seconds. You can change it using ``--cache-timeout`` option.

Also you can change redis address, port and set auth password using ``--redis-address``, ``--redis-port``, ``--redis-password`` options respectively.

```bash
./nord --cache --cache-timeout=300 --redis-host=127.0.0.1 --redis-port=6379 --redis-password="pass"
```


Integration
-----------

Nord functionality can be integrated into your existing Go web application. Below a simple example of such integration:

```go
package main

import (
    "github.com/FZambia/nord/libnord"
    "net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Hello, do you like Nord?"))
}

func main() {
    http.HandleFunc("/", handler)

    config := libnord.DefaultConfig
    config.Prefix = "/nord/"
    m := libnord.GetHandler(config)
    http.Handle("/nord/", m)

    http.ListenAndServe(":8000", nil)
}
```

Now Nord available on http://localhost:8000/nord/

License
-------

MIT


Contribute
----------

Pull requests are welcome

