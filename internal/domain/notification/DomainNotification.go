package notification

const (
	NotFound      = "Notific : {{.Entity}} não encontrado"
	InvalidData   = "Notific : Dados do {{.Entity}} inválidos: {{if .Data}}{{.Data}}{{end}}"
	InvalidMethod = "Notific : Método não permitido"
)

const (
	ScanErrorRepository    = "Notific : Erro ao realizar o scan de dados para as entidades de domain: {{if .Data}}{{.Data}}{{end}}"
	FindErrorRepository    = "Notific : Erro ao realizar a consulta: {{if .Data}}{{.Data}}{{end}}"
	FindAllErrorRepository = "Notific : Erro ao realizar as consultas: {{if .Data}}{{.Data}}{{end}}"
)

const (
	ErrorDbFatal         = "Erro ao conectar ao banco de dados: %v"
	ErrorRepositoryFatal = "Erro ao iniciar os repositorios: %v"
)

const (
	ErrorConfigDb       = "Notific : Configuração do banco de dados não está definida"
	ErrorOpenConnection = "Notific : Erro ao abrir conexão:{{if .Data}}{{.Data}}{{end}}"
	ErrorTestConnection = "Notific : Erro ao testar conexão: {{if .Data}}{{.Data}}{{end}}"
)

const (
	LogForErrorUpdateUsers   = "Erro ao tentar atualizar os usuários: %v"
	LogForPartialUpdateUsers = "Status de %d usuários antigos atualizado para 2"
	LogStartRotineAction     = "Agendador: Próxima atualização agendada para %v (em %v)"
	LogRotineNoStarted       = "Agendador: Desligando antes da primeira execução programada"
	LogNextRotine            = "Agendador: Próxima atualização agendada para %v"
	LogRotineOff             = "Agendador: Desligando"
	LogRotineStoped          = "Agendador: Parado"
	LogMiddleware            = "[%s] %s %s - Status: %d - Duration: %s"
)
