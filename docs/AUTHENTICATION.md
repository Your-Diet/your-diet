# Autenticação e Autorização

Este documento descreve como funciona o sistema de autenticação e autorização da aplicação.

## Visão Geral

A autenticação é feita usando JWT (JSON Web Tokens). Após o login bem-sucedido, um token JWT é retornado e deve ser enviado no cabeçalho `Authorization` das requisições subsequentes.

## Fluxo de Autenticação

1. **Login**: Envie uma requisição POST para `/v1/users/login` com email e senha.
2. **Token**: O servidor retorna um token JWT no formato `Bearer <token>`.
3. **Requisições Autenticadas**: Inclua o token no cabeçalho `Authorization: Bearer <token>`.

## Middleware de Autenticação

O middleware de autenticação verifica a validade do token JWT e adiciona as informações do usuário ao contexto da requisição.

### Como Usar o Middleware

```go
import "github.com/victorgiudicissi/your-diet/internal/middleware"

// Criar o middleware com a chave secreta
authMiddleware := middleware.AuthMiddleware([]byte("sua-chave-secreta"))

// Aplicar a rotas
router := gin.Default()
protected := router.Group("/api")
protected.Use(authMiddleware)
{
    // Rotas protegidas aqui
}
```

### Acessando Dados do Usuário nos Handlers

```go
func MeuHandler(c *gin.Context) {
    // Obter as claims do contexto
    claims, exists := c.Get(string(middleware.TokenContextKey))
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
        return
    }

    // Fazer type assertion para o tipo Claims
    userClaims, ok := claims.(*middleware.Claims)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar autenticação"})
        return
    }

    // Usar os dados do usuário
    userID := userClaims.UserID
    email := userClaims.Email
    permissions := userClaims.Permissions
    
    // ... resto do handler
}
```

## Autorização Baseada em Permissões

O sistema suporta autorização baseada em permissões. Você pode proteger rotas específicas exigindo permissões específicas:

```go
import "github.com/victorgiudicissi/your-diet/internal/constants"

router := gin.Default()
protected := router.Group("/api")
protected.Use(authMiddleware)
{
    // Rota que requer permissão específica
    protected.POST("/diets", middleware.HasPermission(constants.PermissionCreateDiet), dietHandler.Create)
}
```

## Tipos de Usuário e Permissões

- **Usuário Padrão (DEFAULT)**:
  - `list_diet`: Visualizar dietas

- **Nutricionista (NUTRITIONIST)**:
  - `list_diet`: Visualizar dietas
  - `create_diet`: Criar novas dietas
  - `update_diet`: Atualizar dietas existentes
  - `upload_file`: Fazer upload de arquivos

## Segurança

- O token JWT tem uma validade de 10 minutos por padrão
- A chave secreta deve ser armazenada em uma variável de ambiente em produção
- Todas as rotas protegidas requerem um token válido
- As senhas são armazenadas usando bcrypt com salt

## Variáveis de Ambiente

- `JWT_SECRET`: Chave secreta para assinar os tokens JWT
- `JWT_EXPIRATION`: Tempo de expiração do token (padrão: 10m)

## Endpoints de Dieta

### Listar Dietas do Usuário

Retorna todas as dietas do usuário autenticado.

**Endpoint:** `GET /v1/diets`

**Headers:**
- `Authorization: Bearer <seu-token-jwt>`

**Resposta de Sucesso (200 OK):**
```json
[
  {
    "id": "60d5f1b3b58d8b001f8e4e1a",
    "user_id": "60d5f1b3b58d8b001f8e4e19",
    "user_email": "usuario@exemplo.com",
    "name": "Dieta de Exemplo",
    "duration_in_days": 30,
    "status": "active",
    "created_at": "2023-06-06T12:00:00Z",
    "updated_at": "2023-06-06T12:00:00Z"
  }
]
```

**Possíveis Erros:**
- `401 Unauthorized`: Token inválido ou ausente
- `500 Internal Server Error`: Erro ao processar a requisição

## Exemplo de Uso com cURL

```bash
# Fazer login
curl -X POST http://localhost:8080/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario@exemplo.com", "password": "senha123"}'

# Listar dietas do usuário
curl http://localhost:8080/v1/diets \
  -H "Authorization: Bearer <seu-token-jwt>"
```
