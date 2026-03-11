# XPLR Skills (вызывай через @skill)

## @add-card-type
Когда просят новый тип карты:
1. domain/card/newtype.go
2. обнови domain/card/type.go
3. добавь в application/card/service.go
4. новая миграция goose

## @implement-antifraud
Реализуй антифрод:
- Считай попытки в transaction_repo
- После 3 неудач — блокируй карту через domain/transaction/errors.go
- Используй wrapper для логирования

## @add-migration
Создавай миграцию только через make migrate-create name=...

## @use-wrapper
Всегда используй gitlab.com/libs-artifex/wrapper вместо errors.New / fmt.Errorf

## @use-envparse
Конфиг парси только через gitlab.com/libs-artifex/envparse