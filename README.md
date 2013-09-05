gourl
=====

A URL shortener framework written in Go.

You can use this with Google's App Engine simply by using the gae
package. The front end is rather minimal, but you are welcome to jazz
it up.

You can include this in your own Go application by simply doing two
things: create a structure that implements the DataStore interface,
and attach the handlers to your applications web server. You can see
an example of this in the gae packages source code.

This product includes GeoLite2 data created by MaxMind, available from
<a href="http://www.maxmind.com">http://www.maxmind.com</a>.
