# URL Shortener Service

## Summary

Nowadays, people can get information easily in many ways. Either browse websites directly or via social network sharing.  Basically all those informations are shared in form of URL links or screenshots/multimedia attachments. 

### Business Requirement

As LinkAja customer, i want to receive information from the company and view it on my desktop and mobile browser. But i don't wanna type that long URL address specially from media that unable me to do copy-paste (eg: prints, instagram, etc).

As LinkAja marketing officer, i want to share information about our promotions, company news from our website to customers easily without causing them problem when accessing it. Therefore, instead of sharing this exact URL:  
* https://www.linkaja.id/promo/linkaja-cashback-20-tiket-masuk-ciputra-waterpark-pake-linkaja?utm_source=sms&utm_medium=lba&utm_campaign=ciputrawaterparksurabaya_cb20_juli

i want to share this following URL instead:
* https://lnkj.me/wtps20

The the short URL version, customers will not be hassled to type the URL on their browser.
I want to do this by calling an API and later on using a web as the interface for creating and maintaining the shortened URL

### Acceptance Criteria

- **as customer**
-- upon accessing the shortened URL, i will be redirected to the exact URL and my browser can display the content
- **as marketing officer**
-- i will call an API with exact long URL as input and the system will show me the shortened version.


### Technical Approach

1. There will be a mapping between long URL to short URL
2. For starting, we will use Redis as our database and the URL map will be saved as key-value pair.

#### ***Customer***
1. Type the shortened URL in my browser, eg: https://lnkj.me/wtps20 the method can be HTTP GET or HTTP POST
2. Web server will receive the request and pass it to URL Shortener Service to extract the first segment (*wtps20*)
3. URL Shortener Service (lets call it URL-SS) will lookup to Redis whether key ***wtps20*** exist in database
4. if EXIST, URL-SS will get the long URL and construct HTTP 301 or HTTP 302 request using the records as destination.
5.  if NOT EXIST, URL-SS will return HTTP 404 to tell the browser that URL does not exist 

**Sample of Redis records**

| KEY            |VALUE                          |
|----------------|-------------------------------|
|wtps20|`https://www.linkaja.id/promo/linkaja-cashback-20-tiket-masuk-ciputra-waterpark-pake-linkaja?utm_source=sms&utm_medium=lba&utm_campaign=ciputrawaterparksurabaya_cb20_juli`|
|zahra |`https://www.instagram.com/zahra.daniar/` |
|wheelock |`http://www.bu.edu/academics/programs-for-masters-and-bachelors-degree-students-who-transitioned-to-boston-university-from-wheelock-college`|

#### ***Marketing Officer***
1. Using [Postman](https://www.getpostman.com/) to access URL-SS API for generating short url 
2. The URL will be http://zahra.mac:8080/short
3. Method will be HTTP POST only
4. Request payload will be in JSON format as follow:
`{ 
 "longUrl":"TheExactLongURL", 
 shortUrl:"(optional)" 
 }`

5. Response will also be in JSON format as follow:

    `{"status":"statusCode", "shortUrl":"theShortUrl", "timestamp":unixTimeStamp}`

**shortURL field is optional** 
 -- if it is empty then URL-SS should create 6 random characters, case sensitive as the shortened URL 
 -- if it is defined, it will be used as the shortened URL

___Example with empty shortURL___:
___(Request)___:

`{ 
 "longUrl":"https://github.com/zahradaniar/First-Go-Project/blob/master/hello.go", 
 shortUrl:"" 
 }`

___(Response)___:

`{ 
 "status":"OK", 
 shortUrl:"wt5y6A",
 timestamp:1564065576 
 }`

the URL you have to type will be http://zahra.mac/wt5y6A

___Example with given shortURL___:

___(Request)___

`{ 
 "longUrl":"https://github.com/zahradaniar/First-Go-Project/blob/master/hello.go", 
 shortUrl:"MyGo" 
 }`

___(Response)___

`{ 
 "status":"OK", 
 shortUrl:"MyGo",
 timestamp:1564065576 
 }`

the URL you have to type will be http://zahra.mac/MyGo

### Complementary information
1. URL for customer will be http://zahra.mac/TheShortUrl using HTTP GET/POST
2. URL for officer will be http://zahra.mac:8080/short using HTTP POST only
3. [How to create your own zahra.mac address](https://www.macworld.com/article/1137189/servebyname.html) you dont need the apache section. You have to build your own webserver using Go.
4. Some references:
- https://golang.org/doc/articles/wiki/
- https://hackernoon.com/how-to-create-a-web-server-in-go-a064277287c9
- https://gowebexamples.com/http-server/
- https://yourbasic.org/golang/http-server-example/
- https://www.restapiexample.com/golang-tutorial/simple-example-uses-redis-golang/
- https://socketloop.com/tutorials/golang-how-to-return-http-status-code

![Sample Success Request, shortURL found in Redis](https://gitlab.com/nurikidy/zahra-daniar-internship/blob/master/zahra_sample_01.jpg)

![Sample Failed Request, shortURL not found in Redis](https://gitlab.com/nurikidy/zahra-daniar-internship/blob/master/zahra_sample_02.jpg)
