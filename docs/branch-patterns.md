
---

### 📌 **Общий формат**

```
<тип>/<краткое-описание>
```

или если проект связан с таск-трекером:

```
<тип>/<номер-задачи>-<краткое-описание>
```

(номер задачи можно опустить, если его нет)

---

## 🔹 **Основные типы веток**

| Тип        | Когда использовать                         | Примеры                                           |
| ---------- | ------------------------------------------ | ------------------------------------------------- |
| `feature`  | Добавление новой функциональности          | `feature/rabbitmq-producer`, `feature/user-auth`  |
| `fix`      | Исправление бага                           | `fix/message-queue-error`, `fix/login-bug`        |
| `refactor` | Улучшение кода без изменения логики        | `refactor/cleanup-rabbitmq-handler`               |
| `chore`    | Обслуживание проекта (конфиги, обновления) | `chore/update-dependencies`, `chore/setup-linter` |
| `docs`     | Обновление документации                    | `docs/update-readme`, `docs/rabbitmq-guide`       |
| `test`     | Добавление или исправление тестов          | `test/rabbitmq-unit-tests`, `test/fix-auth-tests` |
| `ci`       | Настройка CI/CD                            | `ci/github-actions`, `ci/fix-build`               |
| `hotfix`   | Срочное исправление в продакшене           | `hotfix/fix-payment-crash`                        |
| `release`  | Подготовка к релизу                        | `release/1.2.0`, `release/candidate`              |

---

### 🔹 **Примеры использования**

#### ✅ **Для новой фичи:**

```bash
feature/rabbitmq-integration
feature/user-profile-page
feature/add-dark-mode
```

#### ✅ **Для исправления багов:**

```bash
fix/rabbitmq-connection-timeout
fix/missing-user-avatar
fix/payment-gateway-error
```

#### ✅ **Для рефакторинга:**

```bash
refactor/split-auth-module
refactor/optimize-db-queries
refactor/move-queue-handlers
```

#### ✅ **Для документации и тестов:**

```bash
docs/api-documentation
docs/update-readme
test/add-rabbitmq-tests
test/fix-integration-tests
```

#### ✅ **Для CI/CD и конфигураций:**

```bash
ci/github-actions
ci/fix-dockerfile
chore/update-eslint-config
chore/add-prettier
```

---

### 📌 **Дополнительные рекомендации**

- 🔹 **Используй дефисы вместо подчеркиваний** → `feature/add-user-login`, а не `feature_add_user_login`
- 🔹 **Старайся не делать названия слишком длинными**  
    ❌ `feature/i-want-to-add-some-fancy-rabbitmq-queue-processor`  
    ✅ `feature/rabbitmq-queue-processor`
- 🔹 **Если ветка тестовая или временная, используй префикс `wip/`**
    
    ```bash
    wip/experiment-with-rabbitmq
    wip/new-cache-strategy
    ```
    

Этот паттерн делает `git branch` удобным и читаемым. 🚀
