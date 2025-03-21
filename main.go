package main
//entrypoint to our programme.

func main() {
	app := App{}
	app.Initialize(DbUser, DbPassword, DbName)
	app.Run("localhost:10000")
}