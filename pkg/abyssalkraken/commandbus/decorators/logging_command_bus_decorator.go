package decorators

import (
	"github.com/abyssal-kraken/abyssalkraken/pkg/abyssalkraken/commandbus"
	"log"
	"time"
)

type LoggingCommandBusDecorator struct{}

func (d LoggingCommandBusDecorator) Decorate(command commandbus.Command[any], execution func() (any, error)) func() (any, error) {
	return func() (any, error) {
		log.Printf("Starting commandbus: %T", command)

		startTime := time.Now()
		result, err := execution()
		duration := time.Since(startTime)

		if err != nil {
			log.Printf("Error executing commandbus: %T, duration: %v, error: %v", command, duration, err)
		} else {
			log.Printf("Finished commandbus: %T in %v", command, duration)
		}

		return result, err
	}
}
