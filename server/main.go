package main

import (
	"context"
	"fmt"
	"server/handlers"
	"server/judger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	code := `
		#include <iostream>
		using namespace std;
		int main(){
			int n;
			cin>>n;
			cout<<n+1;
		}
	`
	inputs := make([]string, 0)
	inputs = append(inputs, "1")
	outputs := make([]string, 0)
	outputs = append(outputs, "2")
	res := judger.Judge_samples(code, inputs, outputs)
	fmt.Println(res)

	h := handlers.GetHandler(context.Background())
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowCredentials = true
	middle := cors.New(config)
	r.Use(middle)
	h.Init_user(r)
	h.Init_grader(r)
	h.Init_problems(r)
	h.Init_submission(r)
	r.Run()
}
