# Inviscan - Projeto em Andamento

Sistema de varredura e reconhecimento de superfície de ataque com frontend em React + TypeScript e backend em Go.

## ✨ Funcionalidades

- Descoberta de subdomínios (`subfinder`, `assetfinder`, `crt.sh`)
- Verificação de serviços ativos com `httprobe`
- Detecção de alvos sensíveis (Juicy Targets)
- Captura de screenshots com `gowitness`
- Histórico de varreduras
- Interface web interativa

---

## ⚙️ Requisitos

### 🔧 Programas obrigatórios

| Ferramenta          | Função                                         |
| ------------------- | ---------------------------------------------- |
| **Go**              | Backend e execução de ferramentas auxiliares   |
| **Bun**             | Gerenciador para frontend com React/Vite       |
| `curl`, `jq`, `sed` | Utilitários shell para integração com `crt.sh` |

### 📦 Ferramentas Go (instalar com um único comando)

```bash
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest ;go install github.com/tomnomnom/assetfinder@latest ;go install github.com/tomnomnom/httprobe@latest ;go install github.com/sensepost/gowitness@latest
```

E adicione ao PATH:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
```

🚀 Como rodar
🔁 Backend

```bash
cd backend
go run main.go
```

💻 Frontend

```bash
cd frontend
bun install
bun run dev
```

📂 Resultados

Os resultados de cada varredura ficam salvos em backend/resultados/<domínio>\_<timestamp>/, contendo:

    subdomains.txt

    active_sites.txt

    juicytargets.txt

    gowitness.sqlite3 (capturas)

Você pode visualizar as capturas com:

```bash
gowitness server
```
