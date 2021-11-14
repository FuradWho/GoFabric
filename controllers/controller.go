package controllers


import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"

	"gofabric/services"
)


func StartIris() {
	App := iris.New()
	App.Use(Cors)
	// users API operate
	usersApi := App.Party("/user")
	{
		usersApi.Use(iris.Compression)
		usersApi.Post("/CreateUser",services.CreateUser)
	}

	App.Listen(":9098")
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

