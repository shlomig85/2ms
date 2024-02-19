package cmd

import (
	"sync"

	"github.com/checkmarx/2ms/secrets"
)

func processItems(engine *secrets.Engine) {
	defer channels.WaitGroup.Done()

	wgItems := &sync.WaitGroup{}
	for item := range channels.Items {
		report.TotalItemsScanned++
		wgItems.Add(1)
		go engine.Detect(item, secretsChan, wgItems, ignoreVar)
	}
	wgItems.Wait()
	close(secretsChan)
}

func processSecrets() {
	defer channels.WaitGroup.Done()

	for secret := range secretsChan {
		report.TotalSecretsFound++
		if validateVar {
			validationChan <- secret
		}
		report.Results[secret.ID] = append(report.Results[secret.ID], secret)
	}
	close(validationChan)
}

func processValidation() {
	defer channels.WaitGroup.Done()

	wgValidation := &sync.WaitGroup{}
	for secret := range validationChan {
		wgValidation.Add(1)
		go secret.Validate(wgValidation)
	}
	wgValidation.Wait()
}
