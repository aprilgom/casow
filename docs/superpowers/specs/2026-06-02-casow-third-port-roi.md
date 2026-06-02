# casow 3차 포팅 ROI 산정

## 기준

- 원본: `/Users/aprilgom/casow/kasuari`
- 대상: `/Users/aprilgom/casow/casow`
- 원칙: 테스트를 먼저 포팅하고, 그 후 코드 포팅
- 현재 상태: 2차 ROI A/B/C/D/E가 `main`에 merge됨. Upstream 외부 integration 테스트(`removal`, `quadrilateral`, `traits`)와 README/lib horizontal boxes 시나리오는 Go 테스트로 반영됨.

## 현재 포팅 상태

- public value API: `Variable`, `Term`, `Expression`, `Strength`, `Constraint`, relation/error 값 타입 포팅됨.
- solver API: `AddConstraint`, `AddConstraints`, `RemoveConstraint`, `HasConstraint`, `AddEditVariable`, `RemoveEditVariable`, `HasEditVariable`, `SuggestValue`, `FetchChanges`, `GetValue`, `Reset` 포팅됨.
- solver behavior: required/soft constraints, edit variables, failed add rollback, change tracking, removal/pivot stress, quadrilateral scenario가 테스트됨.
- Go 차이: Rust operator overloading은 Go helper 메서드와 `NewConstraint(any)`로 대체됨.

## ROI 순 다음 작업

### A. README API 예제 현대화

ROI: 높음

현재 README 예제는 `Var(...)`, `Const(...)`, `ExpressionFromVariable(...)` 중심이라, 최근 추가된 `NewConstraint(any)` API의 장점을 보여주지 못한다.

테스트 먼저 포팅:

- 기존 `usage_test.go`가 README/lib 시나리오를 이미 고정하므로 새 동작 테스트는 불필요
- 문서 예제 빌드 가능성을 높이려면 추후 `ExampleSolver`로 승격 가능

예상 구현:

- README 예제를 `NewConstraint(windowWidth, GreaterOrEqual, 0, Required)` 같은 간결한 형태로 정리
- ratio constraint 예제도 문서에 포함

병렬성:

- 문서 작업이라 코드 포팅과 병렬 가능

### B. Go doc example 추가

ROI: 높음

Go 패키지는 Rust crate-level docs를 그대로 옮길 수 없다. 대신 `ExampleSolver`나 `doc.go`로 사용 흐름을 Go 도구 체인에서 검증 가능하게 만드는 것이 ROI가 높다.

테스트 먼저 포팅:

- `ExampleSolver`는 `go test`가 컴파일을 검증한다
- 출력은 solver 변경 순서가 비결정적일 수 있으므로 `Output:` 검증 없이 컴파일 예제로 시작

예상 구현:

- `doc.go`에 짧은 패키지 설명 추가
- `solver_example_test.go`에 horizontal boxes 축약 예제 추가

병렬성:

- README 작업과 같은 문서/API 표면이라 한 작업으로 묶는 편이 안전

### C. solver error rollback 추가 스트레스

ROI: 중간

failed `AddConstraint` rollback은 테스트됐지만, artificial variable 경로 실패와 기존 edit variable 조합의 스트레스는 더 보강할 수 있다.

테스트 먼저 포팅:

- 이미 값이 있는 시스템에 unsatisfiable required inequality/equality를 추가해 실패시 기존 `HasConstraint`, `GetValue`, `FetchChanges`, edit workflow가 유지되는지 확인
- `AddConstraints`가 upstream처럼 부분 성공을 유지하는지 더 복합적으로 확인

예상 구현:

- 구현 변경 없음이 기대됨
- 실패 시 snapshot/restore 범위 수정

병렬성:

- `solver_test.go` 중심이라 다른 solver 작업과 동시에 진행하지 않는 편이 안전

### D. public API parity matrix 문서화

ROI: 중간

Rust 문법과 Go API는 1:1 문법 포팅이 불가능하다. 남은 포팅 판단을 빠르게 하려면 upstream item별 Go 대응표가 필요하다.

테스트 먼저 포팅:

- 구현 테스트는 이미 대부분 존재한다

예상 구현:

- spec 문서에 `Variable`, `Term`, `Expression`, `Constraint`, `Solver`별 대응표 작성
- Rust 전용 항목(`no_std`, feature flags, operator traits)은 Go 포팅 제외로 명시

병렬성:

- read-only/docs 작업이라 병렬 가능

## 추천 병렬 배치

첫 배치:

- A. README API 예제 현대화
- B. Go doc example 추가
- D. public API parity matrix 문서화

이 셋은 구현 리스크가 낮고, C의 solver stress 작업과 독립적으로 진행 가능하다.

두 번째 배치:

- C. solver error rollback 추가 스트레스

`solver_test.go` 중심이라 단독으로 진행하는 편이 충돌 가능성이 낮다.
