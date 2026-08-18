package main

var router = NewRouter()

func setupRoutes() {
	router.GET("/", func(req Request) Response {
		return Response{
			StatusCode: 200,
			Body:       []byte("Hello, World!"),
		}
	})

	router.GET("/test/", func(req Request) Response {
		return Response{
			StatusCode: 200,
			Body:       []byte("test"),
		}
	})

	router.GET("/test/:user", func(req Request) Response {
		return Response{
			StatusCode: 200,
			Body:       []byte("hello user"),
		}
	})
}
