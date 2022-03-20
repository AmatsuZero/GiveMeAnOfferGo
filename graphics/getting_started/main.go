package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
	"log"
)

func main() {
	if !glfw.Init() {
		log.Fatalln("初始化失败")
	}
	defer glfw.Terminate()

	/* 创建一个Window 和 OpenGL上下文 */
	window := &glfw.Window{}
	window.SetSize(960, 640)
	window.SetTitle("Hello World")

	/* 激活上面创建的OpenGL上下文 */
	window.MakeContextCurrent()

	/* 进入游戏引擎主循环 */
	for !window.ShouldClose() {
		/* Render here */
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(49.0/255.0, 77.0/255.0, 121.0/255, 1.0)

		/* Swap front and back buffers */
		window.SwapBuffers()
		/* 处理鼠标 键盘事件 */
		glfw.PollEvents()
	}
}
