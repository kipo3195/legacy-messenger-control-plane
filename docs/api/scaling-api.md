# 스케일링 판단 API

## 1. 개요

스케일링 판단 API는 현재 서비스의 연결 부하와 Service Registry에 정의된 스케일링 정책을 비교하여 다음 정보를 계산합니다.

* Scale-out 필요 여부
* Scale-in 가능 여부
* 현재 Task 수 유지 여부
* 권장 `desiredCount`

이 API는 판단 결과만 반환하며 ECS Service의 `desiredCount`를 직접 변경하지 않습니다.

실제 Task 수를 변경하려면 판단 결과를 확인한 뒤 제어 API를 별도로 호출해야 합니다.

```text
관측 데이터 조회
      │
      ▼
스케일링 정책 비교
      │
      ▼
Action 및 권장 Task 수 반환
      │
      ▼
운영자 확인
      │
      ▼
Scale 제어 API 호출
```

---

## 2. 관측·판단·제어 API의 차이

| 구분                  | 역할                             | ECS 상태 변경 |
| ------------------- | ------------------------------ | --------: |
| Connection Pressure | 현재 연결 부하와 Task당 연결 수 조회        |         X |
| Scaling Evaluate    | 정책을 적용하여 Scale-in/out 필요 여부 계산 |         X |
| Scale               | ECS Service의 `desiredCount` 변경 |         O |

Connection Pressure API는 현재 부하 상태를 관측하고, Scaling Evaluate API는 관측값에 서비스별 스케일링 정책을 적용합니다.

```text
GET /connection-pressure
현재 부하 관측
        │
        ▼
POST /scaling-evaluate
정책 기반 판단
        │
        ▼
POST /scale
실제 Task 수 변경
```

---

## 3. 스케일링 판단 실행

현재 연결 상태와 Service Registry의 스케일링 정책을 기준으로 권장 Action과 Task 수를 계산합니다.

### 3.1 요청

```http
POST /api/v1/services/{serviceName}/scaling-evaluate
```

### Path Parameters

| 이름            | 타입     | 필수 | 설명                            |
| ------------- | ------ | -: | ----------------------------- |
| `serviceName` | string |  예 | Service Registry에 등록된 논리 서비스명 |

### Query Parameters

없음

### Request Body

없음

스케일링 판단에 필요한 연결 수, Task 수 및 정책값은 Control Plane이 내부에서 조회합니다.

### 3.2 요청 예시

```bash
curl -X POST \
  "${BASE_URL}/api/v1/services/ws/scaling-evaluate"
```

---

## 4. 성공 응답

* HTTP Status: `200 OK`

```json
{
  "serviceName": "ws",
  "ecsServiceName": "ucware-ws-service",
  "clusterName": "ucware-cluster",
  "action": "SCALE_IN",
  "reason": "connectionPerTask is below scale-in threshold",
  "current": {
    "activeConnectionCount": 2,
    "runningTaskCount": 2,
    "desiredCount": 2,
    "connectionPerTask": 1
  },
  "policy": {
    "targetConnectionsPerTask": 1500,
    "scaleOutThreshold": 1200,
    "scaleInThreshold": 450,
    "minCount": 0,
    "maxCount": 3
  },
  "recommendation": {
    "currentDesiredCount": 2,
    "recommendedDesiredCount": 1
  }
}
```

---

## 5. 응답 필드

### 5.1 기본 필드

| 필드               | 타입     | 설명                            |
| ---------------- | ------ | ----------------------------- |
| `serviceName`    | string | Service Registry에 등록된 논리 서비스명 |
| `ecsServiceName` | string | AWS ECS에 등록된 실제 Service 이름    |
| `clusterName`    | string | ECS Cluster 이름                |
| `action`         | string | 정책 적용 결과로 계산된 권장 Action       |
| `reason`         | string | 해당 Action으로 판단한 이유            |
| `current`        | object | 판단에 사용된 현재 운영 상태              |
| `policy`         | object | 판단에 적용된 스케일링 정책               |
| `recommendation` | object | 현재 Task 수와 권장 Task 수          |

클라이언트는 판단 결과를 처리할 때 `reason` 문자열보다 `action` 값을 기준으로 분기하는 것이 좋습니다.

`reason`은 운영자가 판단 근거를 이해하기 위한 설명이며, 문구는 구현 변경에 따라 달라질 수 있습니다.

### 5.2 `current` 필드

| 필드                      | 타입      | 설명                        |
| ----------------------- | ------- | ------------------------- |
| `activeConnectionCount` | number  | CloudWatch에서 조회한 활성 연결 지표 |
| `runningTaskCount`      | integer | 현재 RUNNING 상태인 ECS Task 수 |
| `desiredCount`          | integer | ECS Service가 유지하려는 Task 수 |
| `connectionPerTask`     | number  | 실행 중인 Task 하나당 계산된 연결 수   |

Task당 연결 수는 다음 값들을 기준으로 계산합니다.

```text
connectionPerTask
    = activeConnectionCount / runningTaskCount
```

예시 응답에서는 다음과 같이 계산됩니다.

```text
activeConnectionCount = 2
runningTaskCount      = 2

connectionPerTask = 2 / 2 = 1
```

`runningTaskCount`가 0인 경우에는 0으로 나눌 수 없으므로 별도의 예외 또는 경계 조건 처리가 필요합니다.

### 5.3 `policy` 필드

| 필드                         | 타입      | 설명                            |
| -------------------------- | ------- | ----------------------------- |
| `targetConnectionsPerTask` | integer | Task 하나가 목표로 관리할 연결 수         |
| `scaleOutThreshold`        | number  | Scale-out 검토를 시작하는 Task당 연결 수 |
| `scaleInThreshold`         | number  | Scale-in 검토를 시작하는 Task당 연결 수  |
| `minCount`                 | integer | 서비스에 허용된 최소 Task 수            |
| `maxCount`                 | integer | 서비스에 허용된 최대 Task 수            |

예시 응답의 정책값은 다음과 같습니다.

```text
targetConnectionsPerTask = 1500
scaleOutThreshold        = 1200
scaleInThreshold         = 450
minCount                 = 0
maxCount                 = 3
```

Threshold는 Task당 목표 연결 수를 기준으로 계산합니다.

```text
scaleOutThreshold
    = targetConnectionsPerTask × 0.8

scaleInThreshold
    = targetConnectionsPerTask × 0.3
```

예시 정책에서는 다음과 같습니다.

```text
1500 × 0.8 = 1200
1500 × 0.3 = 450
```

Threshold 비율은 현재 구현 정책을 기준으로 하며, 향후 설정값으로 분리될 수 있습니다.

### 5.4 `recommendation` 필드

| 필드                        | 타입      | 설명                             |
| ------------------------- | ------- | ------------------------------ |
| `currentDesiredCount`     | integer | 현재 ECS Service의 `desiredCount` |
| `recommendedDesiredCount` | integer | 연결 부하와 정책을 기준으로 권장하는 Task 수    |

권장 Task 수는 활성 연결 수를 Task당 목표 연결 수로 나누어 계산합니다.

```text
recommendedDesiredCount
    = ceil(
        activeConnectionCount
        / targetConnectionsPerTask
      )
```

계산된 값은 Service Registry에 정의된 최소·최대 Task 수 범위 안으로 제한합니다.

```text
minCount
    ≤ recommendedDesiredCount
    ≤ maxCount
```

예시 응답에서는 다음과 같습니다.

```text
activeConnectionCount      = 2
targetConnectionsPerTask   = 1500

ceil(2 / 1500) = 1
```

따라서 현재 `desiredCount=2`보다 작은 `1`이 권장됩니다.

---

## 6. Action 상태값

### 6.1 `SCALE_OUT`

Task당 연결 수가 Scale-out 기준 이상이고 현재 Task 수가 최대값보다 작은 경우 반환합니다.

```text
connectionPerTask ≥ scaleOutThreshold
currentDesiredCount < maxCount
```

예시:

```text
connectionPerTask = 1300
scaleOutThreshold = 1200

action = SCALE_OUT
```

실제 Task 수를 증가시키려면 Scale 제어 API를 호출해야 합니다.

```http
POST /api/v1/services/ws/scale
```

### 6.2 `SCALE_IN`

Task당 연결 수가 Scale-in 기준 이하이고 현재 Task 수를 줄일 수 있는 경우 반환합니다.

```text
connectionPerTask ≤ scaleInThreshold
currentDesiredCount > minCount
```

업로드한 응답은 다음 조건에 해당합니다.

```text
connectionPerTask = 1
scaleInThreshold  = 450
currentDesiredCount = 2
recommendedDesiredCount = 1
```

따라서 다음 결과가 반환됩니다.

```json
{
  "action": "SCALE_IN",
  "reason": "connectionPerTask is below scale-in threshold",
  "recommendation": {
    "currentDesiredCount": 2,
    "recommendedDesiredCount": 1
  }
}
```

### 6.3 `KEEP`

현재 연결 부하가 Scale-in과 Scale-out 기준 사이에 있거나 정책 경계로 인해 Task 수를 변경할 수 없는 경우 반환합니다.

예를 들어 연결 부하가 높더라도 현재 Task 수가 이미 `maxCount`라면 현재 상태를 유지합니다.

```text
connectionPerTask ≥ scaleOutThreshold
currentDesiredCount = maxCount

action = KEEP
```

연결 부하가 낮더라도 현재 Task 수가 이미 `minCount`인 경우에도 유지합니다.

```text
connectionPerTask ≤ scaleInThreshold
currentDesiredCount = minCount

action = KEEP
```

### 6.4 `NOT_SCALABLE`

Service Registry에서 해당 서비스가 스케일링 대상이 아닌 것으로 설정된 경우 반환합니다.

예시 설정:

```yaml
ds:
  scalable: false
```

이 경우 연결 부하와 관계없이 자동 스케일링 판단 대상에서 제외합니다.

```text
action = NOT_SCALABLE
```

---

## 7. 판단 순서

스케일링 판단은 다음 순서로 수행됩니다.

```text
Service Registry 조회
        │
        ▼
스케일링 가능 여부 확인
        │
        ▼
ECS Service 상태 조회
        │
        ├─ runningTaskCount
        └─ desiredCount
        │
        ▼
CloudWatch 연결 지표 조회
        │
        ▼
connectionPerTask 계산
        │
        ▼
Threshold 비교
        │
        ▼
권장 desiredCount 계산
        │
        ▼
Action 반환
```

개념적인 판단 규칙은 다음과 같습니다.

```text
scalable = false
    → NOT_SCALABLE

connectionPerTask ≥ scaleOutThreshold
and desiredCount < maxCount
    → SCALE_OUT

connectionPerTask ≤ scaleInThreshold
and desiredCount > minCount
    → SCALE_IN

그 외
    → KEEP
```

---

## 8. 판단 결과와 실제 제어

Scaling Evaluate API의 성공 응답은 ECS Service 상태 변경을 의미하지 않습니다.

예를 들어 다음 결과가 반환되어도:

```json
{
  "action": "SCALE_IN",
  "recommendation": {
    "currentDesiredCount": 2,
    "recommendedDesiredCount": 1
  }
}
```

ECS Service의 `desiredCount`는 여전히 2입니다.

실제로 Task 수를 1개로 변경하려면 다음 API를 별도로 호출합니다.

```bash
curl -X POST \
  "${BASE_URL}/api/v1/services/ws/scale" \
  -H "Content-Type: application/json" \
  -d '{
    "desiredCount": 1,
    "reason": "scaling evaluation result"
  }'
```

현재 구조는 판단과 제어를 분리하여 운영자가 결과를 검토한 뒤 실제 변경을 수행할 수 있도록 합니다.

```text
Scaling Evaluate
판단 결과 제공
        │
        ▼
운영자 검토
        │
        ▼
Scale API
실제 ECS 변경
```

---

## 9. 경계 조건

### 9.1 최소 Task 수

권장 Task 수가 `minCount`보다 작으면 `minCount`로 제한합니다.

```text
calculatedCount = 0
minCount = 1

recommendedDesiredCount = 1
```

예시 응답에서는 `minCount=0`이므로 활성 연결이 전혀 없다면 0이 권장될 수 있습니다.

다만 실제 서비스를 완전히 중지하는 Scale-to-zero 정책을 허용할지는 서비스 특성과 운영 정책에 따라 별도로 결정해야 합니다.

### 9.2 최대 Task 수

권장 Task 수가 `maxCount`보다 크면 `maxCount`로 제한합니다.

```text
calculatedCount = 5
maxCount = 3

recommendedDesiredCount = 3
```

Task당 연결 수가 높더라도 이미 `maxCount`에 도달했다면 추가 Scale-out은 권장할 수 없습니다.

### 9.3 실행 중인 Task가 없는 경우

`runningTaskCount=0`이면 Task당 연결 수를 일반적인 나눗셈으로 계산할 수 없습니다.

이 경우 구현에서는 다음 사항을 별도로 처리해야 합니다.

* `desiredCount=0`인지 확인
* 활성 연결 지표가 존재하는지 확인
* 최소 Task 수가 0인지 확인
* Scale-out 또는 현재 상태 유지 여부 결정

### 9.4 Desired Count와 Running Count 차이

Deployment 또는 Scale 작업이 진행 중이면 다음 값이 다를 수 있습니다.

```text
desiredCount = 3
runningTaskCount = 2
```

이 상태에서 Task당 연결 수는 실제 실행 중인 Task 수를 기준으로 계산되지만, 권장 Task 수 비교는 현재 `desiredCount`를 기준으로 수행됩니다.

진행 중인 Deployment가 있는 상태에서 반복적으로 판단 API를 호출하면 불필요한 중복 제어 판단이 발생할 수 있으므로 운영 시 주의해야 합니다.

---

## 10. 오류 응답

판단 과정에서 오류가 발생하면 다음 형식으로 응답합니다.

```json
{
  "message": "error message"
}
```

* HTTP Status: `500 Internal Server Error`

발생 가능한 조건은 다음과 같습니다.

* `serviceName`이 Service Registry에 등록되지 않음
* ECS Service 이름이 설정되지 않음
* ECS Service 상태 조회 실패
* 실행 중인 Task 수 조회 실패
* CloudWatch 연결 지표 조회 실패
* CloudWatch Datapoint가 존재하지 않음
* 스케일링 정책값이 올바르지 않음
* AWS 인증 또는 IAM 권한 오류

CloudWatch 조회 범위 안에 Datapoint가 없다면 다음과 유사한 오류가 발생할 수 있습니다.

```json
{
  "message": "failed to get active connection count: cloudwatch metric datapoint not found"
}
```

현재 POC에서는 오류 유형이 세분화되지 않아 대부분 `500 Internal Server Error`로 반환될 수 있습니다.


---

## 11. 운영 확인 순서

스케일링 판단 이후 실제 제어가 필요한 경우 다음 순서로 확인합니다.

```text
POST /scaling-evaluate
        │
        ▼
Action 및 recommendedDesiredCount 확인
        │
        ▼
POST /scale
        │
        ▼
GET /status
        │
        ├─ desiredCount
        ├─ runningCount
        └─ pendingCount
        │
        ▼
GET /task
        │
        ▼
GET /target-health
```
