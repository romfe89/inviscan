package scans

import "fmt"

func RunFullScan(domain string) error {
	fmt.Println("Iniciando scan para:", domain)

	subdomains, err := EnumerateSubdomains(domain)
	if err != nil {
		return err
	}

	active, err := ProbeActiveSites(subdomains)
	if err != nil {
		return err
	}

	juicy := FilterJuicyTargets(active)

	err = CaptureScreenshots(active)
	if err != nil {
		return err
	}

	err = CompareWithPrevious(domain, subdomains)
	if err != nil {
		return err
	}

	fmt.Printf("Scan concluído: %d subdomínios | %d ativos | %d juicy\n",
		len(subdomains), len(active), len(juicy),
	)
	return nil
}
