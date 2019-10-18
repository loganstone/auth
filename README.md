# Auth

[![Go Report Card](https://goreportcard.com/badge/github.com/loganstone/auth)](https://goreportcard.com/report/github.com/loganstone/auth)

자주 사용하는 인증서비스를 API 로 제공

# Goal
- [x] 이메일 확인 후 가입
- [x] 비밀번호 인증
- [ ] 비밀번호 변경
- [ ] 비밀번호 초기화
- [ ] OTP 생성
- [ ] OTP 인증
- [ ] OTP 초기화
- [ ] 이메일로 인증 코드 발송
- [ ] 이메일로 발송된 인증 코드 확인

# Frameworks(plan)
* http - gin
* orm - gorm
* log - logrus
* cli - cobra
* uuid - google/uuid
* jwt - jwt-go
* otp - xlzd/gotp

# Running Tests

```shell
$ go test -v ./...
```
