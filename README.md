레거시 Java 메신저를 ECS로 올린 뒤,
운영자가 AWS 콘솔에 직접 들어가지 않아도
서비스 상태 조회, Target Health 확인, Task 위치 확인,
desiredCount 조정, 재배포, 연결 기반 스케일 판단을 할 수 있는
운영 제어 계층을 Go로 구현