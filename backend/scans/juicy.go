package scans

import (
	"strings"
)

var juicyKeywords = []string{
	"dev", "dev1", "dev2", "dev3", "development",
	"test", "testing", "qa", "staging", "hml", "sandbox",
	"demo", "preview", "beta", "alpha", "preprod", "uat",
	"jenkins", "git", "gitlab", "bitbucket", "ci", "cicd",
	"pipeline", "artifactory", "nexus", "registry", "docker", "harbor",
	"login", "signin", "auth", "authentication", "sso", "saml", "oauth",
	"register", "signup", "password", "reset", "forgot", "token",
	"vpn", "remote", "access", "gateway", "firewall",
	"admin", "adminpanel", "manage", "dashboard", "console", "cms",
	"intranet", "internal", "private", "secure", "portal", "support",
	"help", "helpdesk", "it", "ticket", "jira", "confluence", "servicenow",
	"db", "database", "mysql", "postgres", "mongo", "sql", "redis",
	"api", "backend", "tools", "monitoring", "status", "uptime",
	"metrics", "grafana", "prometheus", "logs", "log", "kibana", "elastic",
	"public", "static", "files", "uploads", "content", "assets", "media",
	"old", "backup", "bak", "temp", "tmp", "archive",
}

func FilterJuicyTargets(sites []string) []string {
	var juicy []string

	for _, site := range sites {
		for _, keyword := range juicyKeywords {
			if strings.Contains(site, keyword) {
				juicy = append(juicy, site)
				break
			}
		}
	}
	return juicy
}
