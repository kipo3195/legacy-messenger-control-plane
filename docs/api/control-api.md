# 제어 API

## 1. 개요

제어 API는 Control Plane의 Service Registry에 등록된 AWS ECS Service의 실행 상태를 변경합니다.

현재 다음 두 가지 제어 기능을 제공합니다.

* ECS Service의 `desiredCount` 변경
* ECS Service의 강제 재배포

제어 API는 AWS ECS에 상태 변경 요청을 전달한 직후 응답합니다. 따라서 성공 응답은 요청이 정상적으로 접수되었다는 의미이며, 신규 Task 실행, 기존 Task 종료 또는 Deployment 완료까지 보장하지 않습니다.

제어 요청 이후에는 관측 API를 이용하여 실제 반영 상태를 확인해야 합니다.

```text
제어 API 호출
      │
      ▼
AWS ECS 변경 요청
      │
      ▼
서비스 상태 조회
      │
      ├─ Deployment 상태 확인
      ├─ Task 상태 확인
      └─ Target Health 확인
```

---

## 2. API 목록

| Method | Endpoint                                  | 설명                             |
| ------ | ----------------------------------------- | ------------------------------ |
| `POST` | `/api/v1/services/{serviceName}/scale`    | ECS Service의 `desiredCount` 변경 |
| `POST` | `/api/v1/services/{serviceName}/redeploy` | ECS `forceNewDeployment` 실행    |

---

## 3. 사용 시 주의사항

제어 API는 실제 AWS ECS Service 상태를 변경합니다.

현재 POC에서는 REST API 호출자에 대한 애플리케이션 수준의 인증과 권한 검증을 구현하지 않았습니다.

운영 환경에서는 다음 구성이 필요합니다.

* 내부 관리망 또는 허용된 네트워크에서만 접근
* 운영자 또는 관리자 인증
* 관측 권한과 제어 권한 분리
* 제어 요청에 대한 감사 로그
* 중복 요청 및 오작동 방지
* 변경 전후 상태 기록

API 경로의 `{serviceName}`은 실제 ECS Service 이름이 아니라 Service Registry에 등록된 논리 서비스명을 사용합니다.

```text
ws → xxxxxx-ws-service
ds → xxxxxx-ds-service
ns → xxxxxx-ns-service
```

---

# 4. 서비스 Task 수 변경

ECS Service의 `desiredCount`를 변경하여 서비스의 목표 Task 수를 조정합니다.

현재 Task 수보다 큰 값을 요청하면 Scale-out, 작은 값을 요청하면 Scale-in이 진행됩니다.

## 4.1 요청

```http
POST /api/v1/services/{serviceName}/scale
Content-Type: application/json
```

### Path Parameters

| 이름            | 타입     | 필수 | 설명                            |
| ------------- | ------ | -: | ----------------------------- |
| `serviceName` | string |  예 | Service Registry에 등록된 논리 서비스명 |

### Query Parameters

없음

### Request Body

```json
{
  "desiredCount": 1
}
```

| 필드             | 타입      |  필수 | 설명                           |
| -------------- | ------- | --: | ---------------------------- |
| `desiredCount` | integer |   예 | 변경하려는 ECS Service의 목표 Task 수 |
| `reason`       | string  | 아니요 | Scale 요청 사유                  |

`reason`은 선택 입력값이며 현재 POC에서는 별도의 영속 감사 로그 저장을 보장하지 않습니다.

## 4.2 요청 예시

```bash
curl -X POST "${BASE_URL}/api/v1/services/ws/scale" \
  -H "Content-Type: application/json" \
  -d '{
    "desiredCount": 1
  }'
```

요청 사유를 포함하려면 다음과 같이 호출할 수 있습니다.

```bash
curl -X POST "${BASE_URL}/api/v1/services/ws/scale" \
  -H "Content-Type: application/json" \
  -d '{
    "desiredCount": 3,
    "reason": "morning traffic increase"
  }'
```

## 4.3 성공 응답

* HTTP Status: `200 OK`

```json
{
  "serviceName": "ws",
  "ecsServiceName": "xxxxxx-ws-service",
  "clusterName": "xxxxxx-cluster",
  "previousDesiredCount": 3,
  "desiredCount": 1,
  "runningCount": 0,
  "pendingCount": 0,
  "status": "SCALING_REQUESTED",
  "message": "scale request accepted"
}
```

## 4.4 응답 필드

| 필드                     | 타입      | 설명                                  |
| ---------------------- | ------- | ----------------------------------- |
| `serviceName`          | string  | Service Registry에 등록된 논리 서비스명       |
| `ecsServiceName`       | string  | AWS ECS에 등록된 실제 Service 이름          |
| `clusterName`          | string  | ECS Cluster 이름                      |
| `previousDesiredCount` | integer | 변경 요청 전 ECS Service의 `desiredCount` |
| `desiredCount`         | integer | 요청한 새로운 `desiredCount`              |
| `runningCount`         | integer | 응답 생성 시점에 조회된 RUNNING Task 수        |
| `pendingCount`         | integer | 응답 생성 시점에 조회된 PENDING Task 수        |
| `status`               | string  | Scale 요청 처리 상태                      |
| `message`              | string  | 처리 결과에 대한 설명                        |

## 4.5 상태값

| 값                   | 설명                                   |
| ------------------- | ------------------------------------ |
| `SCALING_REQUESTED` | ECS에 `desiredCount` 변경 요청이 정상적으로 전달됨 |
| `NOOP`              | 현재 `desiredCount`와 요청값이 같아 변경하지 않음   |

`SCALING_REQUESTED`는 변경 요청이 접수되었다는 의미입니다.

다음 항목이 완료되었다는 의미는 아닙니다.

* 신규 Task 생성
* 신규 Task의 RUNNING 전환
* 기존 Task 종료
* Load Balancer Target 등록
* Target Health Check 통과
* Scale-in 대상 연결 Draining 완료

## 4.6 응답 시점 해석

다음과 같은 응답이 반환될 수 있습니다.

```json
{
  "previousDesiredCount": 3,
  "desiredCount": 1,
  "runningCount": 0,
  "pendingCount": 0,
  "status": "SCALING_REQUESTED"
}
```

`runningCount`와 `pendingCount`는 요청 처리 과정에서 조회된 특정 시점의 값입니다.

따라서 응답값만으로 최종 Task 수가 1개로 조정되었다고 판단하면 안 됩니다. 실제 반영 결과는 서비스 상태 조회 API를 통해 확인해야 합니다.

```http
GET /api/v1/services/ws/status
```

필요한 경우 Task와 Target Health도 함께 확인합니다.

```http
GET /api/v1/services/ws/task
GET /api/v1/services/ws/target-health
```

## 4.7 Scale-out 흐름

```text
desiredCount 증가 요청
        │
        ▼
ECS UpdateService 호출
        │
        ▼
신규 Task PENDING
        │
        ▼
신규 Task RUNNING
        │
        ▼
Target Group 등록
        │
        ▼
Target Health 정상
```

## 4.8 Scale-in 흐름

```text
desiredCount 감소 요청
        │
        ▼
ECS UpdateService 호출
        │
        ▼
종료 대상 Task 선정
        │
        ▼
Target Draining
        │
        ▼
Task 종료
        │
        ▼
Service Steady State
```

## 4.9 요청 검증 조건

Scale 요청은 Service Registry에 정의된 서비스별 정책을 기준으로 검증합니다.

주요 검증 조건은 다음과 같습니다.

* `serviceName`이 Service Registry에 등록되어 있어야 함
* 서비스에 실제 `ecsServiceName`이 설정되어 있어야 함
* 서비스가 스케일링 가능한 대상으로 설정되어 있어야 함
* `desiredCount`가 0 이상이어야 함
* `desiredCount`가 서비스별 최소값 이상이어야 함
* `desiredCount`가 서비스별 최대값 이하여야 함

예를 들어 Service Registry에 다음과 같이 정의되어 있다면:

```yaml
ws:
  scalable: true
  minCount: 1
  maxCount: 3
```

허용되는 `desiredCount` 범위는 다음과 같습니다.

```text
1 ≤ desiredCount ≤ 3
```

## 4.10 오류 응답

### 요청 Body 형식 오류

* HTTP Status: `400 Bad Request`

```json
{
  "message": "invalid request body"
}
```

발생 가능한 조건:

* JSON 문법 오류
* `desiredCount` 누락
* `desiredCount` 타입 오류

### Scale 처리 오류

* HTTP Status: `500 Internal Server Error`

```json
{
  "message": "error message"
}
```

발생 가능한 조건:

* 등록되지 않은 `serviceName`
* 스케일링 불가 서비스
* 최소·최대 Task 수 정책 위반
* ECS Service 조회 실패
* ECS `UpdateService` 호출 실패
* AWS 인증 또는 IAM 권한 문제

현재 POC에서는 정책 위반과 AWS API 오류가 모두 `500 Internal Server Error`로 반환될 수 있습니다.

---

# 5. 서비스 재배포

ECS Service에 `forceNewDeployment`를 요청하여 동일한 Task Definition을 기준으로 Task를 순차적으로 교체합니다.

다음과 같은 상황에서 사용할 수 있습니다.

* 동일 Image Tag의 Image가 갱신된 경우
* 외부 설정 또는 Secret 변경 사항을 Task에 반영해야 하는 경우
* 실행 중인 Task를 순차적으로 재기동해야 하는 경우
* Task 상태 이상으로 전체 교체가 필요한 경우

## 5.1 요청

```http
POST /api/v1/services/{serviceName}/redeploy
Content-Type: application/json
```

### Path Parameters

| 이름            | 타입     | 필수 | 설명                            |
| ------------- | ------ | -: | ----------------------------- |
| `serviceName` | string |  예 | Service Registry에 등록된 논리 서비스명 |

### Query Parameters

없음

### Request Body

```json
{
  "reason": "test"
}
```

| 필드       | 타입     |  필수 | 설명        |
| -------- | ------ | --: | --------- |
| `reason` | string | 아니요 | 재배포 요청 사유 |

Request Body 없이 호출할 수 있도록 구현된 경우 빈 JSON Object를 전달할 수 있습니다.

```json
{}
```

## 5.2 요청 예시

```bash
curl -X POST "${BASE_URL}/api/v1/services/ws/redeploy" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "configuration changed"
  }'
```

## 5.3 성공 응답

* HTTP Status: `200 OK`

```json
{
  "serviceName": "ws",
  "ecsServiceName": "xxxxxx-ws-service",
  "clusterName": "xxxxxx-cluster",
  "action": "REDEPLOY",
  "status": "STARTED",
  "reason": "test",
  "desiredCount": 2,
  "runningCount": 1,
  "pendingCount": 0,
  "deployments": [
    {
      "id": "ecs-svc/0198899074006879968",
      "status": "PRIMARY",
      "rolloutState": "IN_PROGRESS",
      "rolloutStateReason": "ECS deployment ecs-svc/0198899074006879968 in progress.",
      "taskDefinition": "arn:aws:ecs:ap-northeast-2:123456789012:task-definition/xxxxxx-ws-task:7",
      "desiredCount": 0,
      "runningCount": 0,
      "pendingCount": 0,
      "createdAt": "2026-07-09T14:10:52.059Z",
      "updatedAt": "2026-07-09T14:10:52.059Z"
    },
    {
      "id": "ecs-svc/1398569661056762211",
      "status": "ACTIVE",
      "rolloutState": "COMPLETED",
      "rolloutStateReason": "ECS deployment ecs-svc/1398569661056762211 completed.",
      "taskDefinition": "arn:aws:ecs:ap-northeast-2:123456789012:task-definition/xxxxxx-ws-task:7",
      "desiredCount": 1,
      "runningCount": 1,
      "pendingCount": 0,
      "createdAt": "2026-06-24T08:20:25.089Z",
      "updatedAt": "2026-07-09T14:10:44.083Z"
    }
  ],
  "message": "force new deployment requested"
}
```

## 5.4 응답 필드

### 기본 필드

| 필드               | 타입      | 설명                              |
| ---------------- | ------- | ------------------------------- |
| `serviceName`    | string  | Service Registry에 등록된 논리 서비스명   |
| `ecsServiceName` | string  | AWS ECS에 등록된 실제 Service 이름      |
| `clusterName`    | string  | ECS Cluster 이름                  |
| `action`         | string  | 수행한 제어 작업                       |
| `status`         | string  | 재배포 요청 처리 상태                    |
| `reason`         | string  | 요청 시 전달한 재배포 사유                 |
| `desiredCount`   | integer | 서비스가 유지하려는 Task 수               |
| `runningCount`   | integer | 응답 시점에 실행 중인 Task 수             |
| `pendingCount`   | integer | 응답 시점에 시작 대기 중인 Task 수          |
| `deployments`    | array   | 재배포 요청 직후 조회한 ECS Deployment 목록 |
| `message`        | string  | 요청 처리 결과 설명                     |

### `deployments` 필드

| 필드                   | 타입      | 설명                                      |
| -------------------- | ------- | --------------------------------------- |
| `id`                 | string  | ECS Deployment 식별자                      |
| `status`             | string  | Deployment 상태 또는 역할                     |
| `rolloutState`       | string  | Deployment 진행 상태                        |
| `rolloutStateReason` | string  | Deployment 상태에 대한 AWS 설명                |
| `taskDefinition`     | string  | 해당 Deployment가 사용하는 Task Definition ARN |
| `desiredCount`       | integer | 해당 Deployment가 목표로 하는 Task 수            |
| `runningCount`       | integer | 해당 Deployment에서 실행 중인 Task 수            |
| `pendingCount`       | integer | 해당 Deployment에서 시작 대기 중인 Task 수         |
| `createdAt`          | string  | Deployment 생성 시각                        |
| `updatedAt`          | string  | Deployment 최근 갱신 시각                     |

## 5.5 상태값

### Action

| 값          | 설명                          |
| ---------- | --------------------------- |
| `REDEPLOY` | ECS `forceNewDeployment` 요청 |

### 요청 상태

| 값         | 설명                     |
| --------- | ---------------------- |
| `STARTED` | 재배포 요청이 ECS에 정상적으로 전달됨 |

### Deployment 상태

| 값         | 설명                                    |
| --------- | ------------------------------------- |
| `PRIMARY` | 신규 Task를 생성하며 목표 상태로 전환 중인 Deployment |
| `ACTIVE`  | 기존 Task가 남아 있는 이전 Deployment          |

### Rollout 상태

| 값             | 설명                    |
| ------------- | --------------------- |
| `IN_PROGRESS` | Task 교체가 진행 중         |
| `COMPLETED`   | Deployment가 정상적으로 완료됨 |
| `FAILED`      | Deployment 완료에 실패함    |

## 5.6 응답 해석

예시 응답에는 두 개의 Deployment가 존재합니다.

```text
PRIMARY
- 신규 Deployment
- rolloutState: IN_PROGRESS

ACTIVE
- 기존 Deployment
- rolloutState: COMPLETED
```

재배포 직후에는 신규 `PRIMARY` Deployment의 `desiredCount`, `runningCount`, `pendingCount`가 모두 0으로 조회될 수 있습니다.

이는 ECS가 새로운 Deployment를 생성했지만 Task 배치가 아직 반영되지 않은 초기 시점일 수 있습니다.

응답의 다음 값:

```json
{
  "action": "REDEPLOY",
  "status": "STARTED",
  "message": "force new deployment requested"
}
```

은 재배포 완료가 아니라 `forceNewDeployment` 요청 시작을 의미합니다.

## 5.7 재배포 진행 흐름

```text
Redeploy 요청
      │
      ▼
forceNewDeployment 호출
      │
      ▼
신규 PRIMARY Deployment 생성
      │
      ▼
신규 Task 생성
      │
      ▼
Target Group 등록 및 Health Check
      │
      ▼
기존 Task Draining
      │
      ▼
기존 Deployment 종료
      │
      ▼
신규 Deployment COMPLETED
```

## 5.8 후속 상태 확인

재배포 이후에는 서비스 상태 조회 API를 사용합니다.

```http
GET /api/v1/services/ws/status
```

다음 항목을 확인합니다.

* `runningCount`가 `desiredCount`와 일치하는지
* `pendingCount`가 0인지
* `PRIMARY` Deployment가 `COMPLETED`인지
* 기존 `ACTIVE` Deployment가 제거되었는지
* ECS Event에 반복적인 Task 시작 실패가 없는지

Task의 실행 상태는 다음 API로 확인합니다.

```http
GET /api/v1/services/ws/task
```

Load Balancer의 트래픽 전달 가능 상태는 다음 API로 확인합니다.

```http
GET /api/v1/services/ws/target-health
```

## 5.9 재배포 제한 조건

다음 상황에서는 재배포 요청이 거부될 수 있습니다.

* `serviceName`이 Service Registry에 등록되지 않음
* 실제 ECS Service 이름이 설정되지 않음
* ECS Service의 `desiredCount`가 0
* 기존 Deployment가 이미 진행 중
* ECS Service 상태 조회 실패
* ECS `UpdateService` 호출 실패
* AWS 인증 또는 IAM 권한 부족

`desiredCount=0`인 서비스는 실행 중이거나 교체할 Task가 없으므로 재배포 대상이 아닙니다.

이미 Deployment가 진행 중인 경우 중복 재배포를 요청하면 Deployment 상태 판단이 복잡해질 수 있으므로 제한하는 것이 적절합니다.

## 5.10 오류 응답

### 요청 Body 형식 오류

* HTTP Status: `400 Bad Request`

```json
{
  "message": "invalid request body"
}
```

### 재배포 처리 오류

* HTTP Status: `500 Internal Server Error`

```json
{
  "message": "error message"
}
```

발생 가능한 조건:

* 등록되지 않은 서비스
* 실행 중인 Task가 없는 서비스
* 기존 Deployment 진행 중
* ECS Service 조회 실패
* ECS 재배포 요청 실패
* AWS 인증 또는 IAM 권한 오류

---

# 6. Scale과 Redeploy의 차이

| 구분                 | Scale                              | Redeploy                        |
| ------------------ | ---------------------------------- | ------------------------------- |
| 목적                 | Task 수 변경                          | 기존 Task 순차 교체                   |
| ECS 호출             | `UpdateService`의 `desiredCount` 변경 | `forceNewDeployment=true`       |
| Task Definition 변경 | 없음                                 | 없음                              |
| Task 수 변경          | 발생 가능                              | 기존 `desiredCount` 유지            |
| 기존 Task 교체         | Scale-in 시 일부 종료                   | 전체 Task 순차 교체                   |
| 주요 상태값             | `SCALING_REQUESTED`                | `STARTED`                       |
| 완료 확인              | Service, Task, Target Health       | Deployment, Task, Target Health |

Scale API는 서비스 용량을 변경하기 위한 API입니다.

Redeploy API는 서비스 용량을 유지하면서 실행 중인 Task를 교체하기 위한 API입니다.

---

# 7. 공통 오류 응답

제어 API 처리 중 오류가 발생하면 다음 형식으로 응답합니다.

```json
{
  "message": "error message"
}
```

|                       상태 코드 | 설명                                      |
| --------------------------: | --------------------------------------- |
|           `400 Bad Request` | 요청 Body 형식 또는 필수 입력값 오류                 |
| `500 Internal Server Error` | Service Registry 정책 검증 또는 AWS API 호출 실패 |

현재 POC에서는 다음 오류가 모두 `500 Internal Server Error`로 반환될 수 있습니다.

* 존재하지 않는 서비스
* 스케일링 불가 서비스
* 최소·최대 Task 수 정책 위반
* 실행 중인 Deployment와 충돌
* AWS ECS API 호출 실패

운영 환경에서는 오류 유형을 구분하여 `404`, `409`, `422`, `502` 등의 상태 코드로 세분화할 수 있습니다.

---

# 8. 운영 확인 순서

## Scale 요청 이후

```text
POST /scale
    │
    ▼
GET /status
    │
    ├─ desiredCount 확인
    ├─ runningCount 확인
    └─ pendingCount 확인
    │
    ▼
GET /task
    │
    ▼
GET /target-health
```

## Redeploy 요청 이후

```text
POST /redeploy
    │
    ▼
GET /status
    │
    ├─ PRIMARY Deployment 확인
    ├─ rolloutState 확인
    └─ ECS Event 확인
    │
    ▼
GET /task
    │
    ▼
GET /target-health
```
