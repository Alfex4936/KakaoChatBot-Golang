# Go 언어 공부

## Gin-Gonic을 이용한 Go언어 공부 및 카카오 챗봇 예제

# 엔드 포인트

## /v1/notices/:num

num만큼 MySQL DB에서 공지를 불러옴

```console
$ [GIN-debug] GET    /v1/notices/:num          --> kakao/controllers.GetAllNotices (4 handlers)
```

```console
$ [GIN-debug] Listening and serving HTTP on :8000

$ [GIN] 2021/02/11 - 22:46:32 | 200 |     29.9501ms |             ::1 | GET      "/v1/notices/5"
```

## /v1/last

MySQL DB에서 가장 마지막 공지를 불러옴

```console
$ [GIN-debug] POST   /v1/last/                 --> kakao/controllers.GetLastNotice (4 handlers)
```

## /v1/today

오늘 공지를 불러옴

```console
$ [GIN-debug] POST   /v1/today/                --> kakao/controllers.GetTodayNotices (4 handlers)
```

## /v1/ask

카테고리 선택 유도

```console
$ [GIN-debug] POST   /v1/ask/                  --> kakao/controllers.AskCategory (4 handlers)
```

## /v1/ask/category

카테고리 선택에 따른 공지 5개 불러옴

```console
$ [GIN-debug] POST   /v1/ask/category          --> kakao/controllers.ShowCategory (4 handlers)
```