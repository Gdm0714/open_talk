# OpenTalk

## 프로젝트 개요

- **프론트엔드**: Flutter (Dart)
- **백엔드**: Go (Gin + GORM)
- **아키텍처**: Feature-First Clean Architecture
- **목표**: 오픈소스 프라이버시 중심 메신저 앱

## 프로젝트 구조

```
open_talk/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── core/
│   │   ├── constants/
│   │   ├── theme/
│   │   ├── utils/
│   │   └── network/
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   ├── chat/
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   └── home/
│   │       ├── data/
│   │       ├── domain/
│   │       └── presentation/
│   └── shared/
│       ├── widgets/
│       └── services/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   └── router/router.go
│   ├── pkg/
│   │   ├── response/response.go
│   │   └── validator/validator.go
│   └── migrations/
├── test/
│   ├── unit/
│   ├── widget/
│   └── helpers/
└── integration_test/
```

## 코딩 컨벤션

### Dart / Flutter
- 네이밍: `snake_case` (파일), `PascalCase` (클래스), `camelCase` (변수/함수)
- 상태관리: Riverpod
- 라우팅: go_router
- 모델: freezed + json_serializable
- 비동기: async/await, FutureBuilder/StreamBuilder
- HTTP: dio

### Go 백엔드
- 네이밍: `PascalCase` (exported), `camelCase` (unexported)
- HTTP: Gin
- ORM: GORM
- 인증: JWT (golang-jwt/jwt/v5)
- 구조: handler → service → repository 패턴
- 에러: 표준 Response 포맷 사용

## 아키텍처 규칙

### Feature 구조
각 feature는 다음 레이어를 포함:
- `data/` - 모델, 데이터소스, 리포지토리 구현
- `domain/` - 엔티티, 유스케이스 (비즈니스 로직)
- `presentation/` - 화면, 위젯, 상태관리

### 의존성 방향
```
presentation → domain ← data
```
presentation과 data는 domain에만 의존. 서로 직접 의존 금지.

## 테스트 전략

### 실행 명령
```bash
# Flutter 단위 + 위젯 테스트
flutter test --coverage

# 통합 테스트
flutter test integration_test/

# 커버리지 리포트
genhtml coverage/lcov.info -o coverage/html

# Go 백엔드 테스트
cd backend && go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 커버리지 목표
- 단위 테스트 (models, repositories, usecases): 80%+
- 위젯 테스트 (주요 화면): 70%+
- 통합 테스트: 핵심 사용자 플로우 커버

### 테스트 규칙
- 새 기능 추가 시 반드시 테스트 동반
- 테스트 삭제 금지
- mock은 mocktail 사용

## 빌드 & 실행

```bash
# Flutter
flutter pub get
flutter run

# Go 백엔드
cd backend && go run cmd/server/main.go

# Go 빌드
cd backend && go build -o bin/server cmd/server/main.go
```

## 환경 변수

`.env` 파일 필수 (`.gitignore`에 포함):
```
# Backend
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=open_talk
DB_USER=postgres
DB_PASSWORD=
JWT_SECRET=your-secret-key

# Flutter
API_BASE_URL=http://localhost:8080
```

## Git 규칙

- 커밋 메시지: conventional commits (feat:, fix:, refactor:, test:, docs:)
- PR 전 `flutter test` 통과 필수
- `.env`, 시크릿 파일 커밋 금지
