# Auth

[![Go Report Card](https://goreportcard.com/badge/github.com/loganstone/auth)](https://goreportcard.com/report/github.com/loganstone/auth)

인증 서비스 API 만들기

# Goal
- [x] 이메일 확인 후 가입
- [x] 비밀번호 인증
- [x] 비밀번호 변경
- [ ] 비밀번호 초기화
- [x] OTP 생성
- [x] OTP 인증
- [x] OTP 초기화
- [ ] 이메일로 인증 코드 발송
- [ ] 이메일로 발송된 인증 코드 확인
- [ ] 관리자 기능 추가

# Frameworks(plan)
* http - gin
* orm - gorm
* log - logrus
* cli - cobra
* uuid - google/uuid
* jwt - jwt-go
* otp - xlzd/gotp
* test - stretchr/testify

# Running Tests

```shell
$ go test -v ./...
```
