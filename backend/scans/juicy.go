package scans

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func FilterJuicyTargets(sites []string, outputDir string) []string {
	terms := []string{
		"dev", "dev1", "dev2", "dev3", "development", "test", "testing", "qa",
		"staging", "hml", "sandbox", "demo", "preview", "beta", "alpha", "preprod", "uat",
		"jenkins", "git", "gitlab", "bitbucket", "ci", "cicd", "pipeline", "artifactory", "nexus", "registry",
		"docker", "harbor", "login", "signin", "auth", "authentication", "sso", "saml", "oauth",
		"register", "signup", "password", "reset", "forgot", "token", "vpn", "remote", "access", "gateway",
		"firewall", "admin", "adminpanel", "manage", "dashboard", "console", "cms", "intranet", "internal",
		"private", "secure", "portal", "support", "help", "helpdesk", "it", "ticket", "jira", "confluence",
		"servicenow", "db", "database", "mysql", "postgres", "mongo", "sql", "redis",
		"api", "backend", "tools", "monitoring", "status", "uptime", "metrics", "grafana", "prometheus",
		"logs", "log", "kibana", "elastic", "public", "static", "files", "uploads", "content", "assets",
		"media", "old", "backup", "bak", "temp", "tmp", "archive",
	}

	var juicy []string
	for _, site := range sites {
		for _, term := range terms {
			if strings.Contains(site, term) {
				juicy = append(juicy, site)
				break
			}
		}
	}

	// Salvar juicytargets.txt
	juicyFile := filepath.Join(outputDir, "juicytargets.txt")
	if err := os.WriteFile(juicyFile, []byte(strings.Join(juicy, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao salvar juicytargets.txt: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Juicy targets encontrados: %d", len(juicy)))
	return juicy
}
