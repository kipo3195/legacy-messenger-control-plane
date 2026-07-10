## 1. 프로젝트 개요

`legacy-messenger-control-plane`은 기존 Java 기반 레거시 메신저를 AWS ECS 환경에서 운영하는 데 필요한 **관측·제어·스케일링 판단 기능을 Go 기반 REST API로 구현한 Control Plane POC**입니다.

이 프로젝트의 핵심 목표는 단순히 AWS ECS 운영 명령을 API로 제공하는 것이 아니라, **업무용 메신저의 WebSocket 연결 부하를 지속적으로 관측하고 실행 중인 Task 수를 조정할 수 있는 스케일링 제어 구조를 검증하는 것**입니다.

업무용 메신저는 일반적인 HTTP 서비스와 달리 사용자가 로그인한 동안 WebSocket 연결을 장시간 유지합니다. 특히 출근 시간대에는 로그인과 WebSocket 연결이 짧은 시간 안에 집중될 수 있으므로, 연결 증가를 빠르게 감지하고 Task당 세션 부하를 기준으로 Scale-out 필요 여부를 판단할 수 있어야 합니다.

선행 프로젝트인 [Legacy Messenger ECS Ops POC](https://github.com/kipo3195/legacy-messenger-ecs-ops-poc)에서는 서버별로 직접 실행하던 Java 메신저 서비스를 컨테이너 이미지로 전환하고, AWS ECS EC2 환경에 배포하는 구조를 구성하여 레거시 메신저를 ECS Service 단위로 배포하고 제어할 수 있는 기반을 마련했습니다.

다만 서비스 상태 조회, Task 수 변경, Target Health 확인, 재배포 등의 운영 작업이 AWS 콘솔, AWS CLI, 개별 Shell Script에 분산되어 있었습니다. 또한 CPU와 메모리 사용률만으로는 각 WebSocket Task가 실제로 얼마나 많은 사용자 세션을 관리하고 있는지 파악하기 어려웠습니다.

`legacy-messenger-control-plane`은 이러한 운영 기능을 하나의 API 계층으로 통합하고, **연결 부하를 기반으로 ECS Service의 확장 필요 여부를 판단하기 위해 구현했습니다.**

Control Plane은 관리 대상 ECS 서비스의 설정과 운영 정책을 Service Registry로 관리하고, AWS ECS, Elastic Load Balancing, CloudWatch API를 조합하여 다음 역할을 수행합니다.

| 영역 | 역할                                                       |
| -- | -------------------------------------------------------- |
| 관측 | ECS Service 상태, 실행 중인 Task, Target Health 및 연결 부하 조회     |
| 제어 | ECS Service의 `desiredCount` 변경 및 강제 재배포                  |
| 판단 | 연결 수와 실행 중인 Task 수를 바탕으로 Scale-out, Scale-in 또는 유지 여부 계산 |

현재 POC에서는 ALB의 `ActiveConnectionCount`와 실행 중인 Task 수를 이용하여 Task당 연결 부하를 계산합니다. 다만 CloudWatch 지표는 일정 주기로 집계되는 인프라 지표이므로, 각 WebSocket Task가 현재 관리하는 로그인 세션 수를 즉시 나타내는 애플리케이션 카운터로 사용하기에는 한계가 있습니다.

따라서 현재 구현은 CloudWatch 기반 연결 지표를 이용해 스케일링 판단 구조를 먼저 검증하는 단계이며, 이후에는 각 WebSocket Task가 실제 세션 수를 Control Plane에 직접 보고하도록 확장할 예정입니다.

이를 통해 최종적으로는 애플리케이션이 직접 보고한 세션 부하와 최소·최대 Task 수, Cooldown 등의 운영 정책을 결합하여 ECS Service의 `desiredCount`를 자동으로 조정하는 구조를 목표로 합니다.
