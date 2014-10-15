Skuë
====

![logo](https://raw.githubusercontent.com/greivinlopez/skue/master/assets/logo.png)

Skuë is a [Go](http://golang.org/) package intended to facilitate the creation of pragmatic [REST](http://en.wikipedia.org/wiki/Representational_state_transfer) APIs for Go servers.

## What "skuë" means?

Skuë means **"mouse"** in [Bribrí](http://en.wikipedia.org/wiki/Bribri_language) which is the language of an indigenous group of people of [Costa Rica](https://www.youtube.com/watch?v=pNTirQ9eoLo), my Country.

## What is it?

Skuë is just some helper interfaces and functions working together to link separate software pieces in order to create an API REST-like server.

It does not force you to use a particular web server implementation or even a particular routing solution.  And it allows you to control the different moving parts of your API separately so you can replace them in the future if you want to.

## How it works?

Let's look at the following diagram describing the architecture of the API server that you'll create with Skuë and explain each part separately

<p align="center">
  <img src="https://raw.githubusercontent.com/greivinlopez/skue/master/assets/archdiagram.png"/>
</p>

### The web server

For Skuë it's no important what web framework or http router do you use, as long as you follow REST style you will be OK.  Let's see a basic example using [martini](https://github.com/go-martini/martini):

~~~ go
func main() {
	m := martini.Classic()
	
	m.Get("/resources/:id", getResourceHandler)
	/* This will respond with a 405 Method Not Allowed
	   status code for an HTTP request with a method
	   different than GET */
	m.Any("/resources/:id", skue.NotAllowed)
	
	http.ListenAndServe(":3020", m)
}
~~~

So far so good. Nothing different of what you would expect of any other REST server.

### The Skuë layer

Here is the place you start using the helper functions.  The most important functions you will be using are the persistance utils:

~~~
skue.Create
skue.Read
skue.Update
skue.Delete
skue.List
~~~ 

Through those functions you create the API by providing valid implementations of the interfaces defined by Skuë: 

~~~
Consumer
Producer
DatabasePersistor
MemoryCacher
~~~

The interface implementations are passed as parameters to the persistance functions.  How this layer looks like in your server code? Let's continue with the basic example:

~~~ go
// GET a resource by id
func getResourceHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	id := params["id"]
	resource := models.NewResource(id)
	skue.Read(view, resource, cache, w, r)
}
~~~

In the above code `resource` represents an implementation of the `skue.DatabasePersistor`, `view` represents an implementation of `skue.ViewLayer` and `cache` represents an implementation of the `skue.MemoryCacher` interface.

### The view layer

The view layer represents the implementation of two interfaces: `skue.Consumer` and `skue.Producer`. 
`skue.Producer` is intended to be an encoder that writes a value to http writers for a particular [MIME](http://en.wikipedia.org/wiki/Internet_media_type) type.
`skue.Consumer` is intended to be a decoder of HTTP requests that uses a particular [MIME](http://en.wikipedia.org/wiki/Internet_media_type) type to decode the intended object into a value.

~~~ go
type ViewLayer struct {
	Producer Producer
	Consumer Consumer
}

type Producer interface {
	MimeType() string
	Out(w http.ResponseWriter, statusCode int, value interface{})
}

type Consumer interface {
	MimeType() string
	In(r *http.Request, value interface{}) error
}
~~~

Skuë provides two already implemented view layers: `XmlView` and `JSONView`.

To continue with the basic example, let's consume and produce JSON format in our API:

~~~ go
import (
	"github.com/greivinlopez/skue/views"
)

var (
	view   skue.ViewLayer
)

func init() {
	view = *views.NewJSONView()
}
~~~

### The memory layer

`skue.MemoryCacher` represents an abstraction of any memory caching system used to speed up data driven systems by caching data in RAM instead of HD.  That technique is called [caching](http://en.wikipedia.org/wiki/Cache_(computing)).

There are two important details about this layer:

* It's responsability of the `skue.DatabasePersistor` to interact properly with the `skue.MemoryCacher` to ensure it is actually used.
* This layer is optional.

Skuë provides an implementation of the `skue.MemoryCacher` for [Redis](http://redis.io/).

## Credits

### Icons

* Tablet designed by <a href="http://www.thenounproject.com/dreamer810">Qing Li</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Imac designed by <a href="http://www.thenounproject.com/sofiamoya">Sofía Moya</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Import designed by <a href="http://www.thenounproject.com/howlettstudios">Christopher T. Howlett</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Export designed by <a href="http://www.thenounproject.com/howlettstudios">Christopher T. Howlett</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Eye designed by <a href="http://www.thenounproject.com/sergidelgado">Sergi Delgado</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* RAM designed by <a href="http://www.thenounproject.com/brynbodayle">Bryn Bodayle</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Database designed by <a href="http://www.thenounproject.com/anton.outkine">Anton Outkine</a> from the <a href="http://www.thenounproject.com">Noun Project</a>
* Server designed by <a href="http://www.thenounproject.com/idiotbox">Norbert Kucsera</a> from the <a href="http://www.thenounproject.com">Noun Project</a>


[![License](http://img.shields.io/:license-mit-blue.svg)](http://opensource.org/licenses/MIT)

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)
