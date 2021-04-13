# Auth

[![Go Report Card](https://goreportcard.com/badge/github.com/loganstone/auth)](https://goreportcard.com/report/github.com/loganstone/auth)

인증 서비스 API 만들기


# Goal
- [ ] GORM 을 SQLboiler + sql-migrate 로 변경하기
  * SQLboiler 가 SQL 과 직관적인 맵핑
  * GORM 일부 기능이 제대로 작동하지 않는 문제
- [x] 이메일 확인 후 가입
- [x] 비밀번호 인증
- [x] 비밀번호 변경
- [ ] 비밀번호 초기화
- [x] OTP 생성
- [x] OTP 인증
- [x] OTP 초기화
- [ ] 이메일로 인증 코드 발송
- [ ] 이메일로 발송된 인증 코드 확인
- [x] 관리자 기능 추가
- [ ] 기능이 처리 되었음을 알리는 이벤트 전달 (kafka 지원)
  - https://github.com/confluentinc/confluent-kafka-go 사용 예정


# Prerequisites

* 로컬 메일서버(postfix)가 필요합니다. 테스트 실행 시에는 필요하지 않습니다.
* MariaDB 서버가 필요합니다.
* 필수 환경 변수 설정이 필요합니다.

```shell
$ export AUTH_DB_NAME=<your dbname, required>
$ export AUTH_DB_ID=<your db id, required>
$ export AUTH_DB_PW=<your db password, required>
```

# Running Tests

```shell
$ go test -v -count=1 ./...  # no cached
```
