# Auth
자주 사용하는 인증서비스를 API 로 제공

# Goal
* 이메일 확인 후 사용자 생성
* 비밀번호 인증
* 비밀번호 변경
* 비밀번호 초기화
* OTP 생성
* OTP 인증
* OTP 초기화
* 이메일로 verification code 발송
* 이메일로 발송된 verification code 확인

# Frameworks(plan)
* http - gin
* orm - gorm
* log - logrus
* cli - cobra

# Running Tests

```shell
$ go test -v ./...
```
