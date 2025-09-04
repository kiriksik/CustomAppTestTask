# Numeric Multiplier Generator Service

**.env не в .gitignore для удобной проверки.**

- Переменная `PORT` в .env отвечает за порт, на котором будет запущен сервер (64333)

- docker-compose запускает сервер вместе с тестом, при наличии его в папке

## Алгоритм генерации

1. **Вычисление текущего среднего ожидаемого RES:**

```go
if s.count == 0 {
    expectedRES = 0
} else {
    expectedRES = s.sumExpectedGenerated / s.count
}
```

2. **Вычисление вероятности получения нового числа:** 
```go
if expectedRES > s.rtp {
    p_zero = 1 // если среднее превышает целевой RTP, зануляем каждый следующий мультипликатор
} else {
    meanGenerator := s.maxGenerate / 2
    p_zero = 1 - (((s.count + 1) * s.rtp - s.sumExpectedGenerated) / meanGenerator)
}
```

3. **Генерация случайного числа от 0 до maxGenerate:** 

`maxGenerate` = 500, т.к. при больших значениях возможен излишний дребезг и падает точность. При большм количестве итераций следует уменьшать maxGenerate (вплоть до 100)
```go
multiplier = rand.Float64() * s.maxGenerate
if rand.Float64() < p_zero {
    multiplier = 0
}
```

4. **Обновление статистики:** 
```go
s.count++
s.sumExpectedGenerated += multiplier * (multiplier / s.maxClientGenerate) / 2
```

