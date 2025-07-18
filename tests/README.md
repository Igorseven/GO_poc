# Estratégia de Testes para PocGo

Este documento descreve a abordagem de testes para o projeto PocGo, com foco em testes unitários e de integração com alta coesão e baixo acoplamento.

## Princípios de Teste

1. **Alta Coesão**: Cada teste se concentra numa unidade específica de funcionalidade
2. **Baixo Acoplamento**: Os testes devem minimizar dependências de sistemas externos
3. **Legibilidade**: Os testes devem ser fáceis de entender e manter
4. **Confiabilidade**: Os testes devem produzir resultados consistentes
5. **Cobertura**: Os testes devem cobrir caminhos críticos e casos de borda

## Tipos de Teste

### Testes Unitários

Os testes unitários focam em testar componentes individuais isoladamente. As dependências são simuladas (mocked) ou substituídas (stubbed).

- Localização: `tests/unit/`
- Convenção de nomenclatura: `*_test.go`
- Executar com: `go test ./tests/unit/...`

### Testes de Integração

Os testes de integração verificam se diferentes componentes funcionam corretamente juntos.

- Localização: `tests/integration/`
- Convenção de nomenclatura: `*_test.go`
- Executar com: `go test ./tests/integration/...`

## Utilitários de Teste

Utilitários comuns de teste estão localizados em:
- `tests/mocks/`: Implementações simuladas de interfaces
- `tests/fixtures/`: Dados e fixtures de teste
- `tests/helpers/`: Funções auxiliares para testes
- `tests/integration/testutils/`: Utilitários para testes de integração

### TestBootstrap para Testes de Integração

O `TestBootstrap` é um componente que simula o `Bootstrap` da aplicação para testes de integração. Ele permite testar componentes reais do sistema sem iniciar o servidor ou agendadores.

#### Benefícios do TestBootstrap

1. **Inicialização simplificada**: Configura automaticamente banco de dados, repositórios e serviços
2. **Componentes reais**: Usa implementações reais em vez de mocks para testes mais abrangentes
3. **Limpeza de recursos**: Gerencia automaticamente a liberação de recursos após os testes
4. **API consistente**: Fornece uma interface similar ao Bootstrap real da aplicação

#### Como usar o TestBootstrap

O TestBootstrap está disponível em `tests/integration/testutils/TestBootstrap.go`. Para usá-lo:

```go
// Importar o pacote
import "PocGo/tests/integration/testutils"

// Verificar se o banco de dados está disponível
testutils.SkipIfNoDatabase(t)

// Usar o TestBootstrap para o teste
testutils.WithTestApplication(t, func(t *testing.T, app *testutils.TestApplication) {
    // Arrange - Preparar dados de teste
    testUser, err := app.GetTestUser(t, "standard")
    if err != nil {
        t.Fatalf("Falha ao obter usuário de teste: %v", err)
    }

    // Act - Usar serviços diretamente do TestApplication
    user, err := app.Services.User.GetById(testUser.ID)

    // Assert - Verificar resultados
    helpers.AssertNoError(t, err, "Não deve retornar um erro")
    helpers.AssertNotNil(t, user, "Usuário não deve ser nulo")
})
```

#### Exemplo completo

Veja exemplos de uso do TestBootstrap em `tests/integration/services/UserService_testbootstrap_test.go`.

## Usuários de Teste para Integração

### Abordagem com Usuários Existentes

Para testes de integração que envolvem usuários, utilizamos usuários existentes na base de dados ao invés de criar novos. Isso é necessário porque o gerenciamento de usuários é feito pelo Identity no projeto C#.

#### Como Utilizar Usuários Existentes nos Testes

1. Os usuários de teste estão definidos no `TestUserRegistry` em `tests/integration/testutils/DatabaseTestUtils.go`
2. Para obter um usuário de teste, use o método `GetTestUser`:

```go
// Obter um usuário de teste padrão
testUser := db.GetTestUser(t, "standard")

// Usar o ID do usuário em testes
user, err := userService.GetById(testUser.ID)
```

3. Usuários disponíveis:
   - `standard`: Usuário padrão para a maioria dos testes
   - `admin`: Usuário com permissões de administrador
   - `inactive`: Usuário inativo para testar casos específicos

4. Ao modificar um usuário em testes, sempre restaure os valores originais após o teste:

```go
// Armazenar valores originais
originalName := testUser.Name
originalEmail := testUser.Email

// Modificar o usuário
updateUser := &entity.User{
    ID:     testUser.ID,
    Name:   "Nome Modificado",
    Email:  "email.modificado@exemplo.com",
}
userService.Update(updateUser)

// Restaurar valores originais após o teste
restoreUser := &entity.User{
    ID:     testUser.ID,
    Name:   originalName,
    Email:  originalEmail,
}
userService.Update(restoreUser)
```

### Configuração de Novos Usuários de Teste

Para adicionar novos usuários de teste:

1. Adicione o usuário ao banco de dados através do projeto C# com Identity
2. Adicione o ID e informações do usuário ao `TestUserRegistry` em `DatabaseTestUtils.go`
3. Utilize o usuário nos testes através do método `GetTestUser`

## Melhores Práticas

1. Use testes orientados por tabela para testar múltiplos cenários
2. Use nomes de teste significativos que descrevam o que está a ser testado
3. Siga o padrão AAA (Arrange, Act, Assert - Preparar, Agir, Verificar)
4. Simule dependências externas
5. Use subtestes para melhor organização
6. Mantenha os testes independentes e idempotentes
7. Utilize usuários existentes para testes de integração ao invés de criar novos
8. Restaure o estado original dos dados após modificações em testes
