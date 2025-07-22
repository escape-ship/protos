# Escape Ship - Protocol Buffers Definitions

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![Buf Version](https://img.shields.io/badge/buf-v2-blue.svg)](https://buf.build)
[![Latest Release](https://img.shields.io/github/v/tag/escape-ship/protos)](https://github.com/escape-ship/protos/releases)

이 저장소는 Escape Ship 전자상거래 플랫폼의 공유 Protocol Buffers 정의를 포함합니다. 모든 마이크로서비스에서 사용되는 gRPC 서비스 정의를 제공합니다.

This repository contains shared Protocol Buffers definitions for the Escape Ship e-commerce platform. It provides gRPC service definitions used across all microservices.

## 📋 목차 (Table of Contents)

- [서비스 개요 (Services Overview)](#서비스-개요-services-overview)
- [아키텍처 (Architecture)](#아키텍처-architecture)
- [프로젝트 구조 (Project Structure)](#프로젝트-구조-project-structure)
- [설치 및 사용법 (Installation & Usage)](#설치-및-사용법-installation--usage)
- [개발 가이드 (Development Guide)](#개발-가이드-development-guide)
- [빌드 명령어 (Build Commands)](#빌드-명령어-build-commands)
- [서비스 통합 (Service Integration)](#서비스-통합-service-integration)
- [버전 관리 (Version Management)](#버전-관리-version-management)
- [기여하기 (Contributing)](#기여하기-contributing)

## 🚀 서비스 개요 (Services Overview)

### AccountService - 계정 관리
- **Kakao OAuth 통합**: 카카오 로그인 URL 생성 및 콜백 처리
- **사용자 인증**: 로그인 및 회원가입 기능
- **엔드포인트**:
  - `GET /oauth/kakao/login` - 카카오 로그인 URL 조회
  - `POST /oauth/kakao/callback` - 카카오 OAuth 콜백
  - `POST /login` - 사용자 로그인
  - `POST /register` - 사용자 회원가입

### OrderService - 주문 관리
- **주문 생성**: 새로운 주문 등록
- **주문 조회**: 전체 주문 목록 조회
- **엔드포인트**:
  - `POST /v1/order/insert` - 주문 생성
  - `GET /v1/order` - 주문 목록 조회

### PaymentService - 결제 관리
- **Kakao Pay 통합**: 카카오페이 결제 처리
- **결제 플로우**: Ready → Approve → Cancel
- **엔드포인트**:
  - `POST /payment/kakao/ready` - 결제 준비
  - `POST /payment/kakao/approve` - 결제 승인
  - `POST /payment/kakao/cancel` - 결제 취소
- **OpenAPI 문서** 자동 생성 지원

### ProductService - 상품 관리
- **상품 카탈로그**: 상품 및 카테고리 관리
- **상품 옵션**: 상품별 옵션 및 옵션값 관리
- **엔드포인트**:
  - `GET /products` - 전체 상품 조회
  - `GET /products/{id}` - 특정 상품 조회
  - `POST /products` - 상품 등록
  - `POST /product/{id}/options` - 상품 옵션 조회

## 🏗️ 아키텍처 (Architecture)

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Gateway       │    │   Frontend      │    │   Mobile App    │
│   Service       │    │   Application   │    │   Application   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
          ┌────────────────────────────────────────────────────┐
          │                gRPC Gateway                        │
          └────────────────────┬───────────────────────────────┘
                               │
          ┌────────────────────┴───────────────────────────────┐
          │                                                    │
          │              shared-protos                         │
          │        (Protocol Buffer Definitions)               │
          │                                                    │
          └┬─────────┬─────────┬─────────┬─────────┬──────────┘
           │         │         │         │         │
     ┌─────▼───┐ ┌───▼───┐ ┌───▼────┐ ┌──▼────┐ ┌──▼─────┐
     │Account  │ │ Order │ │Product │ │Payment│ │ Other  │
     │Service  │ │Service│ │Service │ │Service│ │Services│
     └─────────┘ └───────┘ └────────┘ └───────┘ └────────┘
```

## 📁 프로젝트 구조 (Project Structure)

```
protos/
├── account.proto           # 계정 및 인증 서비스 정의
├── order.proto            # 주문 관리 서비스 정의
├── payment.proto          # 결제 서비스 정의 (Kakao Pay)
├── product.proto          # 상품 카탈로그 서비스 정의
├── gen/                   # 생성된 Go 코드 디렉토리
│   ├── *.pb.go           # Protocol Buffer 생성 파일
│   ├── *_grpc.pb.go      # gRPC 생성 파일
│   └── *.pb.gw.go        # gRPC-Gateway 생성 파일
├── buf.yaml              # Buf 설정 파일
├── buf.gen.yaml          # Buf 코드 생성 설정
├── go.mod                # Go 모듈 정의
├── tools.mod             # 빌드 도구 의존성
├── Makefile              # 빌드 자동화 스크립트
└── README.md             # 프로젝트 문서
```

## 📦 설치 및 사용법 (Installation & Usage)

### Go 모듈로 사용하기

```bash
# Go 프로젝트에 추가
go get github.com/escape-ship/protos@latest

# go.mod에서 사용
require github.com/escape-ship/protos v1.0.2
```

### 코드에서 임포트

```go
package main

import (
    pb "github.com/escape-ship/protos/gen"
    "google.golang.org/grpc"
)

func main() {
    // gRPC 클라이언트 생성 예제
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // 서비스 클라이언트 생성
    accountClient := pb.NewAccountServiceClient(conn)
    orderClient := pb.NewOrderServiceClient(conn)
    paymentClient := pb.NewPaymentServiceClient(conn)
    productClient := pb.NewProductServiceClient(conn)
    
    // 사용 예제
    resp, err := accountClient.GetKakaoLoginURL(ctx, &pb.GetKakaoLoginURLRequest{})
}
```

## 🛠️ 개발 가이드 (Development Guide)

### 필수 도구 설치

```bash
# Buf CLI 설치
go install github.com/bufbuild/buf/cmd/buf@latest

# Protocol Buffer 컴파일러 (선택사항, buf가 처리함)
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install protobuf-compiler
```

### 개발 워크플로우

1. **프로토 파일 수정**
   ```bash
   # .proto 파일을 수정하세요
   vim account.proto
   ```

2. **빌드 및 생성**
   ```bash
   # 전체 빌드 (lint + generate)
   make build
   
   # 또는 개별 단계
   make proto_gen
   ```

3. **릴리스**
   ```bash
   # 패치 버전 자동 증가 및 릴리스
   make quick-release
   
   # 또는 수동 버전 지정
   make release
   ```

### Protocol Buffer 컨벤션

- **패키지명**: `go.escape.ship.proto.v1`
- **서비스명**: `Service` 접미사 사용 (예: `AccountService`)
- **필드명**: `snake_case` 사용 (예: `user_id`, `access_token`)
- **메시지명**: `PascalCase` 사용

## 🔧 빌드 명령어 (Build Commands)

### 기본 명령어
```bash
make init              # 초기 설정 (도구 다운로드 + 빌드)
make build             # 전체 빌드 (tidy + update + generate)
make proto_gen         # 프로토 파일만 생성
make clean             # 생성된 파일 정리
```

### 도구 관리
```bash
make tool_download     # 빌드 도구 다운로드
make tool_update       # 빌드 도구 업데이트
```

### 모듈 관리
```bash
make mod-tidy          # Go 모듈 정리
make mod-vendor        # 의존성 벤더링
```

### 버전 관리
```bash
make tag              # 수동 태그 생성
make tag-patch        # 패치 버전 자동 증가
make publish          # Git에 변경사항 푸시
make update-consumers # 소비자 서비스 업데이트
```

### 릴리스 워크플로우
```bash
make release          # 전체 릴리스 프로세스 (수동 버전)
make quick-release    # 빠른 패치 릴리스 (자동 버전)
```

### 품질 검사
```bash
make linter-golangci  # Go 코드 린팅
buf lint              # Protocol Buffer 린팅
```

## 🔗 서비스 통합 (Service Integration)

### 마이크로서비스에서 사용하기

각 마이크로서비스(`accountsrv`, `ordersrv`, `paymentsrv`, `productsrv`, `gatewaysrv`)는 이 공유 프로토 정의를 사용합니다:

```go
// 서비스 구현 예제
type accountServer struct {
    pb.UnimplementedAccountServiceServer
}

func (s *accountServer) GetKakaoLoginURL(
    ctx context.Context, 
    req *pb.GetKakaoLoginURLRequest,
) (*pb.GetKakaoLoginURLResponse, error) {
    // 구현 로직
    return &pb.GetKakaoLoginURLResponse{
        LoginUrl: "https://kauth.kakao.com/oauth/authorize?...",
    }, nil
}
```

### HTTP/JSON API 자동 생성

gRPC-Gateway를 통해 HTTP/JSON API가 자동으로 생성됩니다:

```bash
# gRPC 호출
grpcurl -plaintext localhost:50051 \
  go.escape.ship.proto.v1.AccountService/GetKakaoLoginURL

# HTTP 호출 (동일한 결과)
curl -X GET http://localhost:8080/oauth/kakao/login
```

## 📋 버전 관리 (Version Management)

### 시맨틱 버저닝

이 프로젝트는 [시맨틱 버저닝](https://semver.org/lang/ko/)을 따릅니다:

- **MAJOR** (v2.0.0): 호환되지 않는 API 변경
- **MINOR** (v1.1.0): 하위 호환되는 기능 추가
- **PATCH** (v1.0.1): 하위 호환되는 버그 수정

### 릴리스 노트

각 릴리스의 변경사항은 [GitHub Releases](https://github.com/escape-ship/protos/releases)에서 확인할 수 있습니다.

### 최신 버전 사용하기

```bash
# 소비자 서비스에서 최신 버전으로 업데이트
cd ../paymentsrv
go get github.com/escape-ship/protos@latest
go mod tidy
```

## 🧪 테스트

### Protocol Buffer 유효성 검사

```bash
# Buf를 통한 린팅
buf lint

# 스키마 호환성 검사
buf breaking --against '.git#branch=main'
```

### 생성된 코드 빌드 테스트

```bash
# Go 코드 컴파일 테스트
go build ./gen

# 전체 모듈 검증
go mod verify
```

## 🤝 기여하기 (Contributing)

### 개발 프로세스

1. **이슈 생성**: 새로운 기능이나 버그 수정을 위한 이슈 생성
2. **브랜치 생성**: `feature/새기능` 또는 `fix/버그수정` 브랜치 생성
3. **개발**: Protocol Buffer 정의 수정
4. **테스트**: `make build`로 빌드 테스트
5. **풀 리퀘스트**: 메인 브랜치로 PR 생성

### 코딩 가이드라인

- Protocol Buffer 필드는 `snake_case` 사용
- 서비스명은 `Service` 접미사 필수
- HTTP 엔드포인트는 RESTful 패턴 준수
- 모든 RPC에 적절한 HTTP 매핑 추가

### 리뷰 프로세스

- 모든 변경사항은 코드 리뷰 필수
- CI/CD 파이프라인 통과 필수
- 하위 호환성 유지 필수

## 📞 지원 및 문의 (Support)

- **이슈 트래커**: [GitHub Issues](https://github.com/escape-ship/protos/issues)
- **문서**: [README.md](README.md), [CLAUDE.md](CLAUDE.md)
- **Wiki**: [GitHub Wiki](https://github.com/escape-ship/protos/wiki)

## 📄 라이선스 (License)

이 프로젝트는 [MIT License](LICENSE) 하에 배포됩니다.

---

> **참고**: 이 문서는 지속적으로 업데이트됩니다. 최신 정보는 [GitHub 저장소](https://github.com/escape-ship/protos)에서 확인하세요.