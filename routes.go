package main

var router = NewRouter()

func setupRoutes() {
	router.GET("/", func(req Request) Response {
		return Response{
			StatusCode: 200,
			Body:       []byte("Hello, World!"),
		}
	})

	router.GET("/test", func(req Request) Response {
		return Response{
			StatusCode: 200,
			Body:       []byte("test"),
		}
	})
}
