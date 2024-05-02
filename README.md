# Punch

[![Go Documentation](https://godocs.io/github.com/daxartio/punch?status.svg)](https://godocs.io/github.com/daxartio/punch)

The punch package provides a convenient way to create worker processes that continuously execute code in a loop. This package is particularly useful for scenarios where you need to perform repetitive tasks concurrently.

Key Features:

- Worker Creation: Users can easily create worker instances that encapsulate the code to be executed in a loop.
- Middleware Support: The package offers middleware functionality, allowing users to extend the behavior of workers with additional functionality. Middleware functions can intercept and modify input, output, or behavior of worker processes, enabling tasks such as logging, error handling, rate limiting, and authentication.
- Flexible Configuration: Configuration options are provided to customize the behavior of worker instances according to specific use cases. Users can adjust parameters such as loop interval, maximum execution duration, and error handling strategies to meet their application's needs.
- Error Handling: The package includes robust error handling mechanisms to handle unexpected errors encountered during the execution of worker processes. Users can define error handling logic to gracefully handle failures and maintain application stability.
- Context Support: Context-based cancellation and timeout mechanisms are integrated into worker processes, ensuring graceful termination and resource cleanup in response to external signals or timeouts.

Overall, the punch package offers a powerful and flexible solution for creating and managing worker processes in Go applications. By combining the simplicity of loop-based execution with the extensibility of middleware architecture, it provides a versatile tool for building concurrent, robust, and scalable systems.

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

func main() {
	p := punch.New()

	p.SetHandler(func(ctx punch.Context) error {
		fmt.Println("tick")

		return nil
	})

	p.Use(middleware.IntervalWithConfig[punch.Context](middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))

	if err := p.Run(); err != nil {
		panic(err.Error())
	}
}

```
