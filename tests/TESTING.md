# Guia de Testes para PocGo

Este guia explica como executar testes e adicionar novos testes ao projeto PocGo.

## Estrutura de Testes

O projeto usa uma abordagem estruturada para testes com diretórios separados para diferentes tipos de testes:

```
tests/
├── README.md           # Visão geral da estratégia de testes
├── TESTING.md          # Este guia
├── fixtures/           # Dados de teste
├── helpers/            # Funções auxiliares para testes
├── integration/        # Testes de integração
│   ├── services/       # Testes de integração de serviços
│   └── testutils/      # Utilitários para testes de integração
├── mocks/              # Implementações de mock
└── unit/               # Testes unitários
    └── services/       # Testes unitários de serviços
```

## Executando Testes

### Executando Todos os Testes

Para executar todos os testes do projeto:

```bash
go test ./tests/...
```

### Executando Apenas Testes Unitários

Para executar apenas testes unitários:

```bash
go test ./tests/unit/...
```

### Executando Apenas Testes de Integração

Para executar apenas testes de integração:

```bash
go test ./tests/integration/...
```

### Pulando Testes de Integração

Testes de integração requerem uma conexão com o banco de dados. Se você quiser pulá-los:

```bash
SKIP_DB_TESTS=true go test ./tests/...
```

## Escrevendo Novos Testes

### Testes Unitários

Testes unitários devem testar componentes isoladamente, usando mocks para dependências. Siga estes passos:

1. Crie um novo arquivo de teste no diretório apropriado em `tests/unit/`
2. Use os mocks em `tests/mocks/` para isolar o componente
3. Use as funções auxiliares em `tests/helpers/` para asserções
4. Siga o padrão de testes orientados por tabela para múltiplos casos de teste

Exemplo:

```go
func TestAlgumaCoisa(t *testing.T) {
    // Casos de teste
    tests := []struct {
        nome          string
        entrada       string
        configuracaoMock func(*mocks.AlgumMock)
        saidaEsperada string
        erroEsperado  error
    }{
        // Casos de teste aqui
    }

    for _, tt := range tests {
        t.Run(tt.nome, func(t *testing.T) {
            // Arrange (Preparar)
            mock := mocks.NewAlgumMock()
            tt.configuracaoMock(mock)
            servico := NewServico(mock)

            // Act (Agir)
            resultado, err := servico.FazerAlgumaCoisa(tt.entrada)

            // Assert (Verificar)
            if tt.erroEsperado != nil {
                helpers.AssertError(t, err, "Deve retornar um erro")
            } else {
                helpers.AssertNoError(t, err, "Não deve retornar um erro")
                helpers.AssertEqual(t, tt.saidaEsperada, resultado, "O resultado deve corresponder")
            }
        })
    }
}
```

### Testes de Integração

Testes de integração verificam se os componentes funcionam corretamente juntos. Siga estes passos:

1. Crie um novo arquivo de teste no diretório apropriado em `tests/integration/`
2. Use a função `testutils.WithTestDB` para configurar e limpar o banco de dados de teste
3. Use a função `testutils.SkipIfNoDatabase` para pular o teste se o banco de dados não estiver disponível
4. Use as funções auxiliares em `tests/helpers/` para asserções
5. **Use usuários existentes** para testes que envolvem usuários (não crie novos usuários)

#### Usando Usuários Existentes em Testes

Para testes que envolvem usuários, use o método `GetTestUser` para obter usuários existentes do registro de usuários de teste:

```go
func TestAlgumaCoisa_Integracao(t *testing.T) {
    // Pular se o banco de dados não estiver disponível
    testutils.SkipIfNoDatabase(t)

    // Executar teste com banco de dados
    testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
        // Arrange (Preparar)
        // Obter um usuário de teste existente
        testUser := db.GetTestUser(t, "standard") // Ou "admin", "inactive"
        
        // Criar componentes reais
        repo := repository.NewRepository(db.DB)
        servico := service.NewService(repo)

        // Act (Agir)
        resultado, err := servico.FazerAlgumaCoisa(testUser.ID)

        // Assert (Verificar)
        helpers.AssertNoError(t, err, "Não deve retornar um erro")
        // Mais asserções
    })
}
```

#### Modificando Usuários em Testes

Se você precisar modificar um usuário durante um teste, sempre restaure os valores originais após o teste:

```go
func TestAtualizarUsuario_Integracao(t *testing.T) {
    testutils.SkipIfNoDatabase(t)

    testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
        // Arrange (Preparar)
        testUser := db.GetTestUser(t, "standard")
        
        // Guardar valores originais
        originalName := testUser.Name
        originalEmail := testUser.Email
        
        // Criar componentes reais
        userRepo := repository.NewUserRepository(db.DB)
        userService := service.NewUserService(userRepo)
        
        // Dados de atualização
        updateUser := &entity.User{
            ID:     testUser.ID,
            Name:   "Nome Temporário para Teste",
            Email:  "email.temporario@teste.com",
            Status: testUser.Status,
        }

        // Act (Agir)
        err := userService.Update(updateUser)

        // Assert (Verificar)
        helpers.AssertNoError(t, err, "Não deve retornar um erro")
        
        // Verificar se a atualização foi bem-sucedida
        updatedUser, _ := userService.GetById(testUser.ID)
        helpers.AssertEqual(t, updateUser.Name, updatedUser.Name, "Nome deve ser atualizado")
        
        // Restaurar valores originais após o teste
        restoreUser := &entity.User{
            ID:     testUser.ID,
            Name:   originalName,
            Email:  originalEmail,
            Status: testUser.Status,
        }
        userService.Update(restoreUser)
    })
}
```

#### Exemplo Básico de Teste de Integração

```go
func TestAlgumaCoisa_Integracao(t *testing.T) {
    // Pular se o banco de dados não estiver disponível
    testutils.SkipIfNoDatabase(t)

    // Executar teste com banco de dados
    testutils.WithTestDB(t, func(t *testing.T, db *testutils.TestDB) {
        // Arrange (Preparar)
        // Configurar dados de teste

        // Criar componentes reais
        repo := repository.NewRepository(db.DB)
        servico := service.NewService(repo)

        // Act (Agir)
        resultado, err := servico.FazerAlgumaCoisa()

        // Assert (Verificar)
        helpers.AssertNoError(t, err, "Não deve retornar um erro")
        // Mais asserções
    })
}
```

## Melhores Práticas

1. **Use Testes Orientados por Tabela**: Para testar múltiplos cenários com a mesma lógica de teste
2. **Siga o Padrão AAA**: Arrange (Preparar), Act (Agir), Assert (Verificar)
3. **Teste Casos de Borda**: Inclua testes para condições de erro e casos de borda
4. **Mantenha Testes Independentes**: Cada teste deve poder ser executado independentemente
5. **Limpe Após os Testes**: Sempre limpe os dados de teste, especialmente em testes de integração
6. **Nomes Descritivos**: Use nomes de teste descritivos que expliquem o que está sendo testado
7. **Use Funções Auxiliares**: Use as funções auxiliares para asserções consistentes
8. **Mock de Dependências**: Em testes unitários, faça mock de todas as dependências externas

## Adicionando Novos Mocks

Se você precisar fazer mock de uma nova interface:

1. Crie um novo arquivo em `tests/mocks/` nomeado após a interface
2. Implemente todos os métodos da interface
3. Adicione campos de rastreamento para registrar chamadas de método
4. Adicione campos de função que podem ser configurados para personalizar o comportamento
5. Adicione uma função construtora

Exemplo:

```go
type AlgumRepositorioMock struct {
    EncontrarFunc func(id string) (*entity.AlgumaCoisa, error)

    // Para rastrear chamadas
    EncontrarChamadas []string
}

func NewAlgumRepositorioMock() *AlgumRepositorioMock {
    return &AlgumRepositorioMock{
        EncontrarChamadas: []string{},
    }
}

func (m *AlgumRepositorioMock) Encontrar(id string) (*entity.AlgumaCoisa, error) {
    m.EncontrarChamadas = append(m.EncontrarChamadas, id)
    if m.EncontrarFunc != nil {
        return m.EncontrarFunc(id)
    }
    return nil, errors.New("EncontrarFunc não implementada")
}
```
