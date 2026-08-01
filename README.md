# 🖥️ rdpman

O **rdpman** é um gerenciador de conexões RDP (*Remote Desktop Protocol*) compacto, moderno e seguro para Windows. Desenvolvido para oferecer acesso rápido e organizado aos seus servidores através de uma interface elegante no estilo *widget*.

---

## ⚡ Destaques e Funcionalidades

- **Interface Estilo Widget:** Janela compacta, sem bordas (*frameless*) e em *Dark Mode*.
- **Gerenciamento de Conexões:** Cadastro, edição, exclusão e filtro de busca em tempo real por nome ou host.
- **Credenciais Globais:** Cadastre "Usuários Padrão" para reutilizar credenciais em diferentes conexões de forma rápida.
- **Segurança Nativa (DPAPI):** Criptografia diretamente integrada à DPAPI do Windows para proteger hosts, usuários e senhas gravados no banco SQLite local.
- **Instância Única:** Bloqueio nativo de sistema via Mutex (`kernel32.dll`) para evitar múltiplos processos abertos simultaneamente.
- **Leve e CGO-Free:** Construído em Go com SQLite 100% puro (`modernc.org/sqlite`), facilitando o build sem dependências de compiladores C (GCC/MinGW).

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem Backend:** [Go 1.22+](https://go.dev/)
- **Framework Desktop:** [Wails v2](https://wails.io/)
- **Frontend:** HTML5, CSS3, JavaScript (Vanilla)
- **Banco de Dados:** SQLite (`modernc.org/sqlite`)
- **Segurança:** Windows DPAPI (`syscall` / `kernel32.dll`)

---

## 🚀 Como Executar o Projeto

### Pré-requisitos

1. **Windows 10/11**
2. **Go 1.22** ou superior ([Download](https://go.dev/dl/))
3. **Wails CLI** instalado.

---

## 🔐 Nota de Segurança (DPAPI)
As credenciais no banco de dados são criptografadas com a Data Protection API (DPAPI) do Windows, que vincula a chave de segurança à conta do usuário logado na máquina.

⚠️ Atenção: Mover o arquivo de banco de dados para outro computador ou formatar a máquina tornará os dados criptografados inacessíveis, garantindo a proteção contra cópias não autorizadas do banco.
