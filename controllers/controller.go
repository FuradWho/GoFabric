package controllers


import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gofabric/services"
)


func StartIris() {
	app := iris.New()
	app.Use(Cors)

	testApi := app.Party("/")
	{
		testApi.Get("/test", func(context context.Context) {
			context.JSON("connection success")
		})

		testApi.Get("/LifeCycleChaincodeTest",services.LifeCycleChaincodeTest)
	}

	// users API operate
	usersApi := app.Party("/user")
	{
		usersApi.Post("/CreateUser",services.CreateUser)
	}

	channelApi := app.Party("/channel")
	{
		channelApi.Post("/CreateChannel",services.CreateChannel)
		channelApi.Post("/JoinChannel",services.JoinChannel)
	}

	ccApi := app.Party("/cc")
	{
		ccApi.Post("/CreateCC",services.CreateCC)
	}

	 app.Listen(":9099")

}



// Cors Resolve the CORS
func Cors(ctx iris.Context) {

	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}

