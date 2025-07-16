# inshorts
# insructions 
    inshortfers/config/.env file --> update it according to credentials


    run .mysql file to create database, table and indexes

    ensure database and redis is running properly

    inshortfers/redis --> credentials

    run --> go run main.go



#  endpoints 
    http://localhost:8080/v1/news   POST --> create news articles
    http://localhost:8080/v1/news/category?name=General  GET --> get news by category
    http://localhost:8080/v1/news/score?val=0.8  GET  articles which have score > 0.8
    http://localhost:8080/v1/news/source?val=News18   GET articles by source
    http://localhost:8080/v1/news/nearby?lat=17.900636&lon=77.465262&radius=40   GET articles with in 40 km radius
    http://localhost:8080/v1/news/search?query=salman khan   GET articles by search query


    http://localhost:8080/v1/interaction   POST create some interactions
    http://localhost:8080/v1/interaction/trending?lat=19.683407&lon=73.067455&limit=3   GET trending articles in lat and lon of radius meter 10000000  -->     10000km


