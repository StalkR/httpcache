/*
Package httpcache implements net/http RoundTripper with caching.
Persistent with a cache using files. Volatile with a cache in memory.

The following rules applies for caching:
 - only responses of GET requests are cached (no POST, etc.)
 - only body response is cached (no headers)
 - only response status 200 (OK) and response status 301 (Moved Permanently)
*/
package httpcache
