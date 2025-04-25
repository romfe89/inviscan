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
		"dev", "dev1", "dev2", "dev3", "development",
		"test", "testing", "qa",
		"staging", "stage",
		"hml", "homolog",
		"sandbox",
		"demo",
		"preview", "beta", "alpha",
		"preprod", "uat",
		"jenkins",
		"git", "gitlab", "bitbucket", "gitea", "gogs",
		"ci", "cicd", ".ci.", "-ci-", "pipeline",
		"artifactory", "nexus", "registry", "harbor",
		"docker",
		"tools", "tooling", "util",
		"jira", "confluence", "youtrack",
		"swagger", "openapi", "postman",
		"login", "signin", "logon",
		"auth", "authentication", "authenticate",
		"sso", "saml", "oauth", "oidc", "kerberos",
		"register", "signup",
		"password", "passwd", "pwd", "reset", "forgot",
		"token", "jwt", "session",
		"vpn", "remote", "access", "gateway", "rdp",
		"firewall", "proxy",
		"admin", "administrator", "adm",
		"adminpanel", "panel", "manage", "manager",
		"dashboard", "dash", "board",
		"console", "control", "webconsole",
		"cms", "wordpress", "wp-admin", "drupal", "joomla",
		"intranet", "internal", "corp", "private",
		"secure",
		"portal",
		"support", "help", "helpdesk", "servicedesk",
		"it", "itsm",
		"ticket", "osticket", "zammad", "zendesk",
		"servicenow",
		"db", "database",
		"sql", "mysql", "postgres", "pgsql", "mssql",
		"mongo", "mongodb", "redis", "nosql",
		"api", "api-docs", "graphql", "rest",
		"backend", "internal-api",
		"monitoring", "monitor", "metrics", "stats", "status", "health", "uptime",
		"grafana", "prometheus", "zabbix", "nagios", "icinga",
		"logs", "log", "logging", "kibana", "elastic", "elasticsearch", "graylog", "splunk",
		"public",
		"static", "files", "uploads", "downloads", "content", "assets", "media", "storage",
		"backup", "bak", "bkp", "dump",
		"old", "temp", "tmp", "archive", "export", "import",
		"config", "conf", "settings", "env", ".env",
		"debug", "trace",
		"phpmyadmin", "pgadmin",
		"internal",
	}

	var juicy []string
	seenJuicy := make(map[string]bool)

	for _, site := range sites {
		siteLower := strings.ToLower(site)
		for _, term := range terms {
			if strings.Contains(siteLower, term) {
				if !seenJuicy[site] {
					juicy = append(juicy, site)
					seenJuicy[site] = true
				}
				break
			}
		}
	}

	juicyFile := filepath.Join(outputDir, "juicytargets.txt")
	if err := os.WriteFile(juicyFile, []byte(strings.Join(juicy, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao salvar juicytargets.txt: %v", err))
	}

	if len(juicy) > 0 {
		utils.LogSuccess(fmt.Sprintf("Juicy targets encontrados: %d (salvos em %s)", len(juicy), juicyFile))
	} else {
		utils.LogInfo("Nenhum Juicy Target encontrado com os termos atuais.")
	}

	return juicy
}
