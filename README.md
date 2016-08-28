Downloads zip files from a given URL.  Zip files must be linked in anchor tags in the html returned by the URL.  Files are then unzipped, and child files are read and loaded idempotently into a NEWS_XML list in redis.

##Usage:
    news-downloader <baseUrl> <tempDir> <redisHostNameAndPort>
where <baseUrl> is the URL of the index page which indexes the zips
<tempDir> is a local directory, writeable to the current process, where files will be stored temporarily
<redisHostNameAndPort> is hostname:port to access redis, such as localhost:6379
